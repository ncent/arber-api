package clients

import (
	"bytes"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/yunspace/serverless-golang/examples/aws-golang-kinesis/config"
)

func init() {
	LambdaClient = NewLambdaService(config.NewConfig())
}

var LambdaClient *LambdaService

func (ls LambdaService) InvokeAsync(functionName string, jsonPayload []byte) error {
	out, err := ls.client.InvokeAsync(&lambda.InvokeAsyncInput{
		FunctionName: aws.String(functionName),
		InvokeArgs:   bytes.NewReader(jsonPayload),
	})

	if err != nil {
		log.Printf("There was an error in InvokeAsync: %v", err.Error())
		return err
	}

	log.Printf("InvokeAsync output: %+v", out.String())

	return nil
}

func (ls LambdaService) Invoke(functionName string, jsonPayload []byte) (*lambda.InvokeOutput, error) {
	out, err := ls.client.Invoke(&lambda.InvokeInput{
		FunctionName: aws.String(functionName),
		Payload:      jsonPayload,
	})

	if err != nil {
		log.Printf("There was an error in Invoke: %v", err.Error())
		return nil, err
	}

	log.Printf("Invoke output: %+v", out.String())

	return out, nil
}
