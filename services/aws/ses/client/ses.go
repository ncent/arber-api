package clients

import (
	"context"
	"fmt"
	"log"
	"net/mail"
	"os"
	"strings"

	"github.com/DusanKasan/parsemail"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ses"
	Resolver "gitlab.com/ncent/arber/api/services/appsync"
)

var (
	SESClient *SESService
	s3Client  *s3.S3
	resolver  = Resolver.New()
)

func init() {
	SESClient = NewSESService(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
	sess, err := session.NewSession()
	if err != nil {
		log.Fatalf("AWS NewSession failed: %s", err)
	}
	s3Client = s3.New(sess, aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")))
}

func (sess SESService) ConsumeEmail(
	ctx context.Context,
	record events.S3EventRecord,
	processEmail func(
		resolver Resolver.Resolver,
		sess SESService,
		toAddresses []*mail.Address,
		fromAddress *mail.Address,
		bccAddress []*mail.Address,
		body string,
		subject string,
		attachments []parsemail.Attachment) error,
) error {
	log.Printf("record: %+v", record)

	obj, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
		Key:    aws.String(record.S3.Object.Key),
	})
	if err != nil {
		return fmt.Errorf("S3 GetObject failed: %s", err)
	}

	parsedMail, err := parsemail.Parse(obj.Body)
	if err != nil {
		return fmt.Errorf("ReadMessage failed: %s", err)
	}

	obj.Body.Close()

	log.Printf("Found email Body: %v", parsedMail.TextBody)
	log.Printf("Found email FROM: %v", parsedMail.From)
	log.Printf("Found email To: %+v", parsedMail.To)
	log.Printf("Found email BCC: %v", parsedMail.Bcc)
	log.Printf("Found email Subject: %v", parsedMail.Subject)

	return processEmail(resolver, sess, parsedMail.To, parsedMail.From[0], parsedMail.Bcc, parsedMail.TextBody, parsedMail.Subject, parsedMail.Attachments)
}

func (sess SESService) SendEmail(er EmailRequest) error {
	log.Printf("er: %+v", er)

	err := sess.validateEmails(er.Sender)
	if err != nil {
		return sess.checkMailerError(err)
	}

	input := sess.createEmailInput(er)
	log.Printf("input: %+v", input)

	result, err := sess.client.SendEmail(input)
	if err != nil {
		return sess.checkMailerError(err)
	}

	log.Printf("Email Sent to address: %v", er.Recipient)
	log.Printf("Result: %+v", result)
	return nil
}

func (sess SESService) validateEmails(emails ...string) error {
	identitiesResult, _ := sess.client.ListIdentities(
		&ses.ListIdentitiesInput{
			IdentityType: aws.String("EmailAddress"),
		},
	)
	log.Printf("Identities found: %+v", identitiesResult.Identities)
	identities := make(map[string]string, len(identitiesResult.Identities))
	for _, s := range identitiesResult.Identities {
		identities[*s] = strings.ToLower(*s)
	}
	log.Printf("Identities map: %+v", identities)
	log.Printf("Comparing with emails: %+v", emails)
	for _, email := range emails {
		if _, exists := identities[strings.ToLower(email)]; !exists {
			_, err := sess.client.VerifyEmailAddress(&ses.VerifyEmailAddressInput{EmailAddress: aws.String(email)})
			if err != nil {
				return sess.checkMailerError(err)
			}
		}
	}
	return nil
}

func (sess SESService) createEmailInput(er EmailRequest) *ses.SendEmailInput {
	return &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(er.Recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(er.Html),
				},
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(er.Body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(er.Subject),
			},
		},
		Source: aws.String(er.Sender),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}
}

func (sess SESService) checkMailerError(err error) error {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case ses.ErrCodeMessageRejected:
			log.Printf(ses.ErrCodeMessageRejected, aerr.Error())
		case ses.ErrCodeMailFromDomainNotVerifiedException:
			log.Printf(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
		case ses.ErrCodeConfigurationSetDoesNotExistException:
			log.Printf(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
		default:
			log.Printf(aerr.Error())
		}
	} else {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		log.Printf(err.Error())
	}

	return err
}
