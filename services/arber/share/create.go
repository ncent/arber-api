package user

import (
	"fmt"
	"log"
	"net/mail"

	"gitlab.com/ncent/arber/api/services/appsync"
	Resolver "gitlab.com/ncent/arber/api/services/appsync"
	UserController "gitlab.com/ncent/arber/api/services/arber/user"
)

func CreateShareActionAndTransactionWithParentTransaction(resolver Resolver.Resolver, parentTransactionID string, challengeId string) (*appsync.Transaction, error) {
	log.Printf("Creating a Share Action with Transaction")

	var transaction *appsync.Transaction
	shareAction, err := resolver.CreateShareAction(
		appsync.CreateShareAction{
			ChallengeID: &challengeId,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create ShareAction: %v", err)
	}
	log.Printf("Created ShareAction: %+v", shareAction)

	createTransaction := appsync.CreateTransaction{
		TransactionActionID: shareAction.ID,
	}
	if parentTransactionID != "" {
		createTransaction = appsync.CreateTransaction{
			ParentTransactionID: &parentTransactionID,
			TransactionActionID: shareAction.ID,
		}
	}
	transaction, err = resolver.CreateTransaction(createTransaction)
	if err != nil {
		return nil, fmt.Errorf("Failed to create Transaction: %v", err)
	}
	return transaction, nil
}

func CreateShareActionContacts(resolver Resolver.Resolver, transactionID string, from *mail.Address, tos []*mail.Address) error {
	log.Printf("Creating a Share Action with Transaction")

	if transactionID != "" {
		fromUser, err := UserController.CreateSparseUser(resolver, from)
		if err != nil {
			return err
		}

		transaction, err := resolver.GetTransaction(transactionID)
		if err != nil {
			return fmt.Errorf("Failed to get Parent Transaction: %v", err)
		}

		_, err = resolver.UpdateShareAction(
			appsync.UpdateShareAction{
				ID:     transaction.Action.ID,
				UserID: fromUser.ID,
			},
		)
		if err != nil {
			return fmt.Errorf("Failed to update ShareAction: %v", err)
		}

		for _, to := range tos {
			toUser, _ := UserController.CreateSparseUser(resolver, to)

			_, err = resolver.CreateShareActionContact(
				appsync.CreateShareActionContact{
					ShareActionContactShareActionID: transaction.Action.ID,
					ShareActionContactContactID:     toUser.ID,
				},
			)

			if err != nil {
				break
			}
		}
		if err != nil {
			return fmt.Errorf("Failed to create Share Action Contact: %v", err)
		}
	}

	return nil
}
