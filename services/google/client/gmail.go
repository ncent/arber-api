package clients

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	clients "gitlab.com/ncent/arber/api/services/aws/ses/client"
	"golang.org/x/oauth2"
	gmail "google.golang.org/api/gmail/v1"
)

func (gs *GoogleService) SendMail(config *GoogleConfig, token *oauth2.Token, emailRequestData clients.EmailRequest, ctx context.Context) error {
	return gs.sendMail(config, token, emailRequestData, ctx)
}

func (gs *GoogleService) sendMail(config *GoogleConfig, token *oauth2.Token, emailRequestData clients.EmailRequest, ctx context.Context) error {
	svc, err := gs.getGmailService(config, token, ctx)

	if err != nil {
		log.Printf("There was an issue getting the gmail service: %v", err.Error())
		return err
	}

	message := gs.getGmailMessage(emailRequestData)

	_, err = svc.Users.Messages.Send("me", message).Do()
	if err != nil {
		log.Printf("Unable to send message: %+v, error: %v", message, err.Error())
		return err
	}

	log.Printf("Result fromt sending message: %+v", message)
	return nil
}

func (gs *GoogleService) getGmailService(config *GoogleConfig, token *oauth2.Token, ctx context.Context) (*gmail.Service, error) {
	client := gs.getGmailClient(config, token, ctx)
	srv, err := gmail.New(client)

	if err != nil {
		log.Fatalf("Unable to create service %v", err)
		return nil, err
	}
	return srv, err
}

func (gs *GoogleService) getGmailClient(config *GoogleConfig, token *oauth2.Token, ctx context.Context) *http.Client {
	log.Printf("Refresh token is set to: %v", token.RefreshToken)
	var context context.Context
	if ctx != nil {
		context = ctx
	} else {
		context = gs.ctx
	}
	tkr := &TokenRefresher{
		ctx:          context,
		conf:         config,
		refreshToken: token.RefreshToken,
	}

	rtks := &ReuseTokenSource{
		t:   token,
		new: tkr,
	}
	log.Printf("Token Refresher: %+v", tkr)
	log.Printf("Reuse Token Soruce: %+v", rtks)
	t, err := rtks.Token()
	if err != nil {
		log.Printf("There was an error in getGmailClient retrieving the token: %v", err.Error())
		panic("There was an error in getGmailClient retrieving the token")
	}
	return config.OAuth2Config.Client(gs.ctx, t)
}

func (gs *GoogleService) getGmailMessage(emailRequestData clients.EmailRequest) *gmail.Message {
	header := make(map[string]string)
	header["From"] = emailRequestData.Sender
	header["To"] = emailRequestData.Recipient
	header["Subject"] = emailRequestData.Subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"
	var msg string
	for k, v := range header {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msg += "\r\n" + emailRequestData.Body
	return &gmail.Message{
		Raw: base64.RawURLEncoding.EncodeToString([]byte(msg)),
	}
}
