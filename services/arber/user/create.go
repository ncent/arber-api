package user

import (
	"fmt"
	"log"
	"net/mail"
	"strings"

	"gitlab.com/ncent/arber/api/services/appsync"
	Resolver "gitlab.com/ncent/arber/api/services/appsync"
)

func CreateSparseUser(resolver Resolver.Resolver, from *mail.Address) (*Resolver.User, error) {
	fromAddress := strings.ToLower(from.Address)
	emails := []*string{&fromAddress}
	existingUsers, err := resolver.ListUsersByEmails(
		emails,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to get user: %v", err)
	}

	log.Printf("Got existing user: %+v", existingUsers)
	if len(existingUsers) != 0 {
		return &existingUsers[0], nil
	}

	log.Printf("Creating a sparse user")
	blankUserName := " "
	user, err := resolver.CreateUser(
		appsync.CreateUserInput{
			Emails: emails,
			Names:  []*string{&blankUserName},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create sprase user: %v", err)
	}
	log.Printf("Created sparse user: %+v", user)

	return user, nil
}
