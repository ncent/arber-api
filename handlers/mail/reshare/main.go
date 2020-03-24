package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	Resolver "gitlab.com/ncent/arber/api/services/appsync"
	ReshareService "gitlab.com/ncent/arber/api/services/arber/mail/reshare"
)

var (
	resolver = Resolver.New()
)

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	reshareBody, err := ReshareService.GenerateReshareBodyByChallenge(resolver, event.QueryStringParameters["transactionId"], event.QueryStringParameters["challengeId"])

	if err != nil {
		log.Printf("Failed to get challenge: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       err.Error(),
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       *reshareBody,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
