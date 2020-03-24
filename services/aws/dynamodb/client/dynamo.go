package clients

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	appsync "gitlab.com/ncent/arber/api/services/appsync"
)

func init() {
	DynamoClient = NewDynamoService(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
}

var DynamoClient *DynamoService

func (ds DynamoService) UpdateState(tableName string, dynamoRecords []DynamoDBStreamRecord) error {
	for _, dr := range dynamoRecords {
		originalTransaction := appsync.Transaction{}
		dynamodbattribute.UnmarshalMap(dr.OldImage, &originalTransaction)
		newTransaction := appsync.Transaction{}
		dynamodbattribute.UnmarshalMap(dr.NewImage, &newTransaction)
		originalTransaction.InitStateMachine()
		if originalTransaction.Action.Status.String() != newTransaction.Action.Status.String() {
			log.Printf("Attempting to UpdateState for a transaction: %v", originalTransaction.ID)
			log.Printf("From %v to %v", originalTransaction.Action.Status.String(), newTransaction.Action.Status.String())
			err := originalTransaction.Action.FSM.Event(newTransaction.Action.Status.String())
			if err != nil {
				log.Printf("Successfully Transitioned.")
			} else {
				log.Printf("Failed to Transition. Error: ", err.Error())
				log.Printf("Need to retransition transaction %v from %v to %v", originalTransaction.ID, originalTransaction.Action.Status.String(), newTransaction.Action.Status.String())
			}
		}
	}
	return nil
}
