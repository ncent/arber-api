package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	clients "gitlab.com/ncent/arber/api/services/aws/ses/client"
)

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var emailerRequestData clients.EmailRequest
	log.Printf("body: %+v", req.Body)
	err := json.Unmarshal([]byte(req.Body), &emailerRequestData)
	if err != nil {
		log.Printf("Cannot get body from request: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(err.Error()),
		}, err
	}
	err = clients.SESClient.SendEmail(emailerRequestData)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(err.Error()),
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string("Successfully sent."),
	}, nil
}

func main() {
	lambda.Start(handler)
}
