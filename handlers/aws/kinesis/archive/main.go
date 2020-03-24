package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/firehose"
	clients "gitlab.com/ncent/arber/api/services/aws/kinesis/client"
)

var streamName, _ = os.LookupEnv("AWS_KINESIS_NAME")

func Handler(ctx context.Context, event events.KinesisEvent) ([][]*firehose.Record, error) {
	return clients.FirehoseClient.ArchiveRecords(event, streamName)
}

func main() {
	lambda.Start(Handler)
}
