package mail

import (
	"bufio"
	"log"
	"net/mail"
	"strings"

	"github.com/DusanKasan/parsemail"
	Resolver "gitlab.com/ncent/arber/api/services/appsync"
	AttachmentController "gitlab.com/ncent/arber/api/services/arber/attachment"
	ChallengeController "gitlab.com/ncent/arber/api/services/arber/challenge"
	ShareActionController "gitlab.com/ncent/arber/api/services/arber/share"
	UserController "gitlab.com/ncent/arber/api/services/arber/user"
	clients "gitlab.com/ncent/arber/api/services/aws/ses/client"
	helpers "gitlab.com/ncent/arber/api/services/google/helper"
)

func ProcessInbound(
	resolver Resolver.Resolver,
	sess clients.SESService,
	tos []*mail.Address,
	from *mail.Address,
	bcc []*mail.Address,
	body string,
	subject string,
	attachments []parsemail.Attachment) error {
	toAddress := strings.ToLower(tos[0].Address)
	if strings.HasPrefix(toAddress, "start") {
		var firstAttachmentURL string
		if len(attachments) > 0 {
			attachmentURLs, err := AttachmentController.SaveAttachments(attachments)
			if err != nil {
				log.Printf("Failed to save attachments: %+v", err)
			}
			firstAttachmentURL = *attachmentURLs[0]
		}

		user, err := UserController.CreateSparseUser(resolver, from)
		if err != nil {
			return err
		}
		bodyLines, _ := StringToLines(body)
		challenge, err := ChallengeController.CreateChallenge(
			resolver,
			subject,
			from,
			subject,
			bodyLines[0],
			body,
			firstAttachmentURL,
		)
		if err != nil {
			return err
		}

		helpers.SendStartEmail(*user, *challenge)
	} else {
		var bccAddress string
		if len(bcc) > 0 && len(bcc[0].Address) > 0 {
			bccAddress = strings.ToLower(bcc[0].Address)
		}
		if strings.HasPrefix(bccAddress, "share") {
			// bbcAddress will contain txId as such share+txId@redb.ai
			partsArray := strings.Split(strings.Split(bccAddress, "+")[1], "@")
			transactionID := partsArray[0]

			err := ShareActionController.CreateShareActionContacts(resolver, transactionID, from, tos)
			if err != nil {
				return err
			}
		} else {
			log.Printf("Failed to find a proper route for: %v", bccAddress)
			log.Printf("Failed to find a proper route for: %v", toAddress)
		}
	}
	return nil
}

func StringToLines(s string) (lines []string, err error) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	return
}
