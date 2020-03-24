package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"gitlab.com/ncent/arber/api/services/appsync"
	r "gitlab.com/ncent/arber/api/services/appsync"
	"gitlab.com/ncent/arber/api/services/auth0"
	lambdaClient "gitlab.com/ncent/arber/api/services/aws/lambda/client"
	clients "gitlab.com/ncent/arber/api/services/aws/ses/client"
	google "gitlab.com/ncent/arber/api/services/google/client"
	"golang.org/x/oauth2"
	"google.golang.org/api/people/v1"
)

const CLIENT_APP_URL = ""
const API_URL = ""
const SHORTENER_URL = ""

func CreateOrUpdateUser(resolver r.Resolver, existingUsers []appsync.User, usr *auth0.User, emails []*string, googleUserInfo *people.Person, names []*string, phones []*string, photos []*string, token oauth2.Token) (*appsync.User, error) {
	var user *appsync.User
	var err error

	if len(existingUsers) > 0 {
		log.Printf("Updating user: %v", usr)

		user, err = resolver.UpdateUser(
			appsync.UpdateUserInput{
				ID:           *existingUsers[0].ID,
				Identity:     &usr.ID,
				Emails:       emails,
				Etag:         &googleUserInfo.Etag,
				Names:        names,
				PhoneNumbers: phones,
				Pictures:     photos,
				Token:        &token.RefreshToken,
			},
		)
	} else {
		log.Printf("Creating user: %v", usr)

		user, err = resolver.CreateUser(
			appsync.CreateUserInput{
				Identity:     &usr.ID,
				Emails:       emails,
				Etag:         &googleUserInfo.Etag,
				Names:        names,
				PhoneNumbers: phones,
				Pictures:     photos,
				Token:        &token.RefreshToken,
			},
		)

		swe_err := sendWelcomeEmail(user)
		if swe_err != nil {
			log.Printf("Failed to send welcome email: %v", swe_err.Error())
		}
	}
	return user, err
}

func sendWelcomeEmail(user *appsync.User) error {
	html := fmt.Sprintf(`<p><span style="font-weight: 400;">Welcome %v, thank you for signing up!</span></p>
	<p>&nbsp;</p>
	<p><span style="font-weight: 400;">We&rsquo;ve built Arber to help </span><strong>YOU</strong><span style="font-weight: 400;"> and </span><strong>YOUR NETWORK</strong><span style="font-weight: 400;"> find the best jobs from within one anothers networks. We&rsquo;ve seen that everyone ends up happiest when in-network referrals are hired, so we found a way to help you all capitalized on that!</span></p>
	<br>
	<p><span style="font-weight: 400;">You, and your network, will get paid out by the company directly if anyone in your network gets hired </span><strong>or successfully helps</strong><span style="font-weight: 400;"> find the hire in their networks!</span></p>
	<br>
	<p><strong>[$1,000] You</strong><span style="font-weight: 400;"> -&gt; </span><strong>[$2,000] Your friend</strong><span style="font-weight: 400;"> -&gt; </span><strong>[$4,000] Your friend&rsquo;s friend</strong><span style="font-weight: 400;"> -&gt; </span><strong>[Hired!] Your friend&rsquo;s friend&rsquo;s friend! </strong></p>
	<br>
	<p><span style="font-weight: 400;">We created this video to show you what we&rsquo;re all about:</span></p>
	<p><span style="font-weight: 400;">Click here to watch</span></p>
	<br>
	<p><span style="font-weight: 400;">Best,</span></p>
	<p>&nbsp;</p>
	<p><span style="font-weight: 400;">KK</span></p>
	<p><span style="font-weight: 400;">CEO </span><a href="http://arber.redb.ai"><span style="font-weight: 400;">Arber, an nCent Labs Application</span></a></p>`, strings.Fields(*user.Names[0])[0])

	emailerRequestData := clients.EmailRequest{
		Recipient: *user.Emails[0],
		Sender:    "no-reply@redb.ai",
		Subject:   "You’re IN! Quick video inside :)",
		Html:      html,
	}

	err := clients.SESClient.SendEmail(emailerRequestData)
	return err
}

func SendStartEmail(user appsync.User, challenge appsync.Challenge) error {
	startBody, err := generateStartBody(challenge)
	if err != nil {
		log.Printf("Failed to send Start Email: %v", err)
		return err
	}
	emailerRequestData := clients.EmailRequest{
		Recipient: *user.Emails[0],
		Sender:    "no-reply@redb.ai",
		Subject:   fmt.Sprintf("Start Recruiting Now: %s %s", *challenge.SponsorName, *challenge.Name),
		Html:      *startBody,
	}

	err = clients.SESClient.SendEmail(emailerRequestData)
	return err
}

func generateStartBody(challenge appsync.Challenge) (*string, error) {
	reshareLink, err := ShortenUrl(API_URL + "/reshare?transactionId=&challengeId=" + *challenge.ID)
	if err != nil {
		log.Printf("Failed to generate Start Body: %v", err)
		return nil, err
	}
	startBody := fmt.Sprintf(
		`<p>Thank you for using RedB to help find your dream %s!</p>
		<p><a href="%s">Click here to start your search</a></p>
		<p>To learn more about how RedB works <a href="%s">click here</a></p>`,
		*challenge.Name,
		*reshareLink,
		CLIENT_APP_URL,
	)
	return &startBody, nil
}

func PopulateContacts(resolver r.Resolver, user *appsync.User, token oauth2.Token, ctx context.Context) error {
	googleUserContacts, err := google.GoogleClient.GetContacts(google.GoogleOAuthConfig, &token, ctx)
	if err != nil {
		log.Printf("Failed to get user contacts from google: %v", err)
		return err
	}
	log.Printf("Got user contacts from google")

	var contactsEmails []*string
	for _, gc := range googleUserContacts {
		for _, ea := range gc.EmailAddresses {
			if len(ea.Value) > 0 {
				contactsEmails = append(contactsEmails, &ea.Value)
			}
		}
	}
	if len(contactsEmails) > 0 {

		existingUsersToEmailMap, err := resolver.MapUsersByEmails(
			contactsEmails,
		)

		if err != nil {
			log.Printf("There was a problem in PopulateContacts when getting map of emails: %v", err.Error())
			return err
		}

		for _, gc := range googleUserContacts {
			emails, names, phones, photos := ExtractGooglePersonInformation(resolver, gc)
			for _, email := range emails {
				if _, ok := existingUsersToEmailMap[*email]; ok {
					_, err = resolver.UpdateUser(
						appsync.UpdateUserInput{
							ID:           *existingUsersToEmailMap[*email].ID,
							Emails:       emails,
							Etag:         &gc.Etag,
							Names:        names,
							PhoneNumbers: phones,
							Pictures:     photos,
						},
					)
				} else {
					createdUser, _ := resolver.CreateUser(
						appsync.CreateUserInput{
							Emails:       emails,
							Etag:         &gc.Etag,
							Names:        names,
							PhoneNumbers: phones,
							Pictures:     photos,
						},
					)
					resolver.CreateUserContact(
						appsync.CreateUserContactInput{
							UserContactUserId:    createdUser.ID,
							UserContactContactId: user.ID,
						},
					)
				}
			}
		}
	}

	return nil
}

type Payload struct {
	Body string `json:"body"`
}

type ShortenOutput struct {
	ShortURL string `json:"short_url"`
}

func ShortenUrl(url string) (*string, error) {
	body, err := json.Marshal(map[string]interface{}{
		"url": url,
	})

	jsonPayload, err := json.Marshal(Payload{Body: string(body)})
	if err != nil {
		log.Printf("There was an error marshalling the uel into a ShortenInput json for shorten url")
		log.Printf("No shorten url generated for url: %v", url)
		return nil, err
	}

	lambdaName, _ := os.LookupEnv("SHORTEN_URL_LAMBDA")
	log.Printf("Found lambda to call: %v, with payload: %+v", lambdaName, string(jsonPayload))
	result, err := lambdaClient.LambdaClient.Invoke(lambdaName, jsonPayload)

	if err != nil {
		log.Printf("Cannot get shorten response: %v", err)
		return nil, err
	}

	var resultPayload map[string]interface{}
	err = json.Unmarshal(result.Payload, &resultPayload)
	var shortenOutput ShortenOutput
	err = json.Unmarshal([]byte(resultPayload["body"].(string)), &shortenOutput)

	if err != nil {
		log.Printf("Cannot Unmarshal shorten response: %v", err)
		return nil, err
	}

	shortURL := SHORTENER_URL + "/" + shortenOutput.ShortURL
	return &shortURL, nil
}

func ExtractGooglePersonInformation(resolver r.Resolver, person *people.Person) ([]*string, []*string, []*string, []*string) {
	var emails []*string
	for _, ea := range person.EmailAddresses {
		if len(ea.Value) > 0 {
			emails = append(emails, &ea.Value)
		}
	}

	var names []*string
	for _, n := range person.Names {
		if len(n.DisplayName) > 0 {
			names = append(names, &n.DisplayName)
		}
	}

	var phones []*string
	for _, p := range person.PhoneNumbers {
		if len(p.CanonicalForm) > 0 {
			phones = append(phones, &p.CanonicalForm)
		}
	}

	var photos []*string
	for _, pix := range person.Photos {
		if len(pix.Url) > 0 {
			photos = append(photos, &pix.Url)
		}
	}

	return emails, names, phones, photos
}
