package clients

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/aws/aws-sdk-go/service/firehose/firehoseiface"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
	"github.com/yunspace/serverless-golang/examples/aws-golang-kinesis/config"
)

type KinesisService struct {
	config *config.Config
	client kinesisiface.KinesisAPI
}

func NewKinesisService(cfg *config.Config) *KinesisService {
	return &KinesisService{
		config: cfg,
		client: kinesis.New(session.Must(session.NewSession(aws.NewConfig()))),
	}
}

type FirehoseService struct {
	config *config.Config
	client firehoseiface.FirehoseAPI
}

func NewFirehoseService(cfg *config.Config) *FirehoseService {
	return &FirehoseService{
		config: cfg,
		client: firehose.New(session.Must(session.NewSession(aws.NewConfig()))),
	}
}

type KinesisRecordsExecutor func(*KinesisService, map[string][][]byte)

func (krx *KinesisRecordsExecutor) execute(ks *KinesisService, streamNameRecords map[string][][]byte) error {
	log.Printf("In KinesisRecordsExecutor execute.")
	for streamName, records := range streamNameRecords {
		log.Printf("Publishing %d records to %s proxy stream", len(records), streamName)
		ks.PublishRecords(records, streamName)
	}
	return nil
}

type Payload struct {
	ProxyStreamName string `json:"proxyStreamName"`
	Event           string `json:"event"`
	Timestamp       string `json:"timestamp"`
}
