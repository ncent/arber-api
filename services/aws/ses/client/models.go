package clients

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type SESService struct {
	client *ses.SES
}

func NewSESService(cfg *aws.Config) *SESService {
	return &SESService{
		client: ses.New(session.Must(session.NewSession(cfg))),
	}
}

type EmailRequest struct {
	Recipient string `json:"recipient"`
	Sender    string `json:"sender"`
	Html      string `json:"html,omitempty"`
	Body      string `json:"body,omitempty"`
	Subject   string `json:"subject"`
	ID        string `json:"id,omitempty"`
}
