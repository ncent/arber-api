package main

import (
	"context"
	"log"

	clients "gitlab.com/ncent/arber/api/services/aws/ses/client"
	googleclient "gitlab.com/ncent/arber/api/services/google/client"

	"encoding/json"

	Resolver "gitlab.com/ncent/arber/api/services/appsync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/oauth2"
)

var resolver = Resolver.New()

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var emailerRequestData clients.EmailRequest
	json.Unmarshal([]byte(event.Body), &emailerRequestData)
	log.Printf("Found message: %+v", emailerRequestData)

	user, err := resolver.GetUser(emailerRequestData.ID)

	token := oauth2.Token{
		RefreshToken: *user.Token,
	}

	log.Printf("Generated token %+v", token)

	if err != nil {
		log.Printf("Failed to get user: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       err.Error(),
		}, err
	} else {
		log.Printf("Got existing user: %+v", user)
	}

	err = googleclient.GoogleClient.SendMail(googleclient.GoogleOAuthConfig, &token, emailerRequestData, ctx)

	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 201,
	}, nil
}

func main() {
	lambda.Start(handler)
}
