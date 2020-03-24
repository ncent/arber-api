package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	clients "gitlab.com/ncent/arber/api/services/dynamodb/client"
)

func handler(event DynamoDBEvent) error {
	var dynamoRecords []DynamoDBStreamRecord
	for _, record := range event.Records {
		switch record.EventName {
		case "INSERT":
			fallthrough
		case "MODIFY":
			dynamoRecords = append(dynamoRecords, record.Change)
		default:
			log.Printf("No handler for event %v %v, skipping.", record.EventName, record.EventID)
		}
	}
	return clients.DynamoDBClient.UpdateState("actions", dynamoRecords)
}

func main() {
	lambda.Start(handler)
}
