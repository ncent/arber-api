package main

import (
	"context"
	"log"
	"os"

	"encoding/json"

	Resolver "gitlab.com/ncent/arber/api/services/appsync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	auth0 "gitlab.com/ncent/arber/api/services/auth0"
	lambdaClient "gitlab.com/ncent/arber/api/services/aws/lambda/client"
	google "gitlab.com/ncent/arber/api/services/google/client"
	helpers "gitlab.com/ncent/arber/api/services/google/helper"
	"golang.org/x/oauth2"
)

var (
	resolver = Resolver.New()
)

func Handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var usr *auth0.User
	log.Printf("Found event body: %v", event.Body)
	json.Unmarshal([]byte(event.Body), &usr)
	log.Printf("Found user: %+v", usr)
	var identity = usr.Identities[0]

	log.Printf("Generating token for identity: %+v", identity)
	token := oauth2.Token{
		AccessToken:  identity.AccessToken,
		RefreshToken: identity.RefreshToken,
		Expiry:       usr.CreatedAt.AddDate(0, 0, identity.ExpiresIn),
	}

	log.Printf("Generated token %+v", token)

	googleUserInfo, err := google.GoogleClient.GetMyPerson(google.GoogleOAuthConfig, &token, ctx)
	if err != nil {
		log.Printf("Failed to get user info from google: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       err.Error(),
		}, err
	}
	log.Printf("Got user info from google")

	emails, names, phones, photos := helpers.ExtractGooglePersonInformation(resolver, googleUserInfo)

	existingUsers, err := resolver.ListUsersByEmails(
		emails,
	)

	if err != nil {
		log.Printf("Failed to get user: %v", err)
	} else {
		log.Printf("Got existing user: %+v", existingUsers)
	}

	user, err := helpers.CreateOrUpdateUser(resolver, existingUsers, usr, emails, googleUserInfo, names, phones, photos, token)
	if err != nil {
		log.Printf("Failed to create or update user: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       err.Error(),
		}, err
	}

	populateUserContactsAsync(*user.ID)

	userJSON, err := json.Marshal(user)
	if err != nil {
		log.Printf("Failed to create user json: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       string(userJSON),
	}, nil
}

func populateUserContactsAsync(userID string) error {
	jsonPayload, err := json.Marshal(google.GetContactsPayload{ID: userID})
	if err != nil {
		log.Printf("There was an error marshalling the userID into a GetContactsPayload json for contant population")
		log.Printf("No contacts generated for user: %v", userID)
	}

	lambdaName, _ := os.LookupEnv("POPULATE_USER_CONTACTS_LAMBDA")
	log.Printf("Found lambda to call: %v, with payload: %+v", lambdaName, jsonPayload)
	return lambdaClient.LambdaClient.InvokeAsync(lambdaName, jsonPayload)
}

func main() {
	lambda.Start(Handler)
}
