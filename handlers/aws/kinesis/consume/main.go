package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	clients "gitlab.com/ncent/arber/api/services/aws/kinesis/client"
)

var streamName, _ = os.LookupEnv("AWS_KINESIS_NAME")

func Handler(ctx context.Context, event events.KinesisEvent) (map[string][][]byte, error) {
	return clients.KinesisClient.ConsumeRecords(streamName, nil)
}

func main() {
	lambda.Start(Handler)
}
