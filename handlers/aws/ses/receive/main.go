package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	Resolver "gitlab.com/ncent/arber/api/services/appsync"
	Arber "gitlab.com/ncent/arber/api/services/arber/mail"
	clients "gitlab.com/ncent/arber/api/services/aws/ses/client"
)

var resolver = Resolver.New()

func handler(ctx context.Context, event events.S3Event) error {
	for _, record := range event.Records {

		error := clients.SESClient.ConsumeEmail(ctx, record, Arber.ProcessInbound)
		if error != nil {
			return fmt.Errorf("Event Error %v", error)
		}
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
