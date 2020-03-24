package user

import (
	"fmt"
	"log"
	"net/mail"

	"gitlab.com/ncent/arber/api/services/appsync"
	Resolver "gitlab.com/ncent/arber/api/services/appsync"
)

func CreateChallenge(resolver Resolver.Resolver, subject string, from *mail.Address, sponsor string, name string, description string, attachmentURL string) (*appsync.Challenge, error) {
	log.Printf("Creating a challenge")
	challenge, err := resolver.CreateChallenge(
		appsync.CreateChallenge{
			Name:          &name,
			Description:   &description,
			SponsorName:   &sponsor,
			AttachmentURL: &attachmentURL,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create challenge: %v", err)
	}
	log.Printf("Created challenge: %+v", challenge)

	return challenge, nil
}

func GetChallenge(resolver Resolver.Resolver, id string) (*appsync.Challenge, error) {
	log.Printf("Get a challenge")
	challenge, err := resolver.GetChallenge(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get challenge: %v", err)
	}
	log.Printf("Get challenge: %+v", challenge)

	return challenge, nil
}
