package clients

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/yunspace/serverless-golang/examples/aws-golang-kinesis/config"
)

type LambdaService struct {
	config *config.Config
	client *lambda.Lambda
}

func NewLambdaService(cfg *config.Config) *LambdaService {
	return &LambdaService{
		config: cfg,
		client: lambda.New(session.Must(session.NewSession(aws.NewConfig()))),
	}
}
