package main

import (
	"context"
	"log"

	Resolver "gitlab.com/ncent/arber/api/services/appsync"
	"golang.org/x/oauth2"

	"github.com/aws/aws-lambda-go/lambda"
	googleClient "gitlab.com/ncent/arber/api/services/google/client"
	helpers "gitlab.com/ncent/arber/api/services/google/helper"
)

var resolver = Resolver.New()

func Handler(ctx context.Context, event googleClient.GetContactsPayload) error {
	log.Printf("Found message: %+v", event)

	user, err := resolver.GetUser(event.ID)

	if err != nil {
		log.Printf("Failed to get user: %v", err)
		return err
	}
	log.Printf("Got existing user: %+v", user)

	token := oauth2.Token{
		RefreshToken: *user.Token,
	}

	err = helpers.PopulateContacts(resolver, user, token, ctx)

	if err != nil {
		log.Printf("Failed to populate contacts for user user: %+v, err: %v", user, err.Error())
	}
	return nil
}

func main() {
	lambda.Start(Handler)
}
