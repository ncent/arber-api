package attachment

import (
	"log"

	"github.com/DusanKasan/parsemail"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/dchest/uniuri"
)

const BUCKET_NAME = "redb-inbox-attachments"

func SaveAttachments(attachments []parsemail.Attachment) ([]*string, error) {

	sess := session.Must(session.NewSession())

	uploader := s3manager.NewUploader(sess)
	var attachmentURLs []*string
	randomString := uniuri.New()
	publicRead := "public-read"
	for _, attachment := range attachments {
		result, err := uploader.Upload(&s3manager.UploadInput{
			Body:   attachment.Data,
			Bucket: aws.String(BUCKET_NAME),
			Key:    aws.String(randomString + "-" + attachment.Filename),
			ACL:    &publicRead,
		})

		if err != nil {
			log.Printf("Failed to upload: %v", err)
			return nil, err
		}
		attachmentURLs = append(attachmentURLs, &result.Location)
	}
	return attachmentURLs, nil
}
