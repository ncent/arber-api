package clients

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/yunspace/serverless-golang/examples/aws-golang-kinesis/config"
)

func init() {
	FirehoseClient = NewFirehoseService(config.NewConfig())
}

var FirehoseClient *FirehoseService

const (
	maxBatchSize = 400
)

func (fs FirehoseService) ArchiveRecords(event events.KinesisEvent, streamName string) ([][]*firehose.Record, error) {
	var batchesOfRecords [][]*firehose.Record
	batchOfRecords := make([]*firehose.Record, 0, maxBatchSize)
	records := event.Records
	numRecords := len(records)

	log.Printf("Batching records: %v", numRecords)

	for i, record := range records {
		batchOfRecords = append(batchOfRecords, &firehose.Record{Data: append(record.Kinesis.Data, '\n')})
		if len(batchOfRecords) >= maxBatchSize || i == (numRecords-1) {
			fs.ArchiveBatch(streamName, batchOfRecords)
			batchesOfRecords = append(batchesOfRecords, batchOfRecords)
			batchOfRecords = make([]*firehose.Record, 0, maxBatchSize)
		}
	}

	log.Printf("Completed archiving %v batches", len(batchesOfRecords))
	log.Printf("Recods archived: %v", batchesOfRecords)

	return batchesOfRecords, nil
}

func (fs FirehoseService) ArchiveBatch(streamName string, records []*firehose.Record) error {
	if len(records) < 1 {
		log.Printf("Cannot archive empty records.")
		return nil
	}

	batchResult, err := fs.client.PutRecordBatch(
		&firehose.PutRecordBatchInput{
			DeliveryStreamName: aws.String(streamName),
			Records:            records,
		},
	)
	if err != nil {
		log.Printf("Failed to archive records: %v", err.Error())
		return err
	}
	log.Printf("Successfully archived records: %v", batchResult)
	return nil
}
