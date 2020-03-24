package clients

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
)

func MockKinesisService(mockKinesisClient kinesisClientMock) *KinesisService {
	return &KinesisService{
		config: nil,
		client: mockKinesisClient,
	}
}

func EventRecords(records []events.KinesisRecord) []events.KinesisEventRecord {
	var eventRecords []events.KinesisEventRecord
	for _, r := range records {
		kinRecord := events.KinesisEventRecord{
			Kinesis: r,
		}
		eventRecords = append(eventRecords, kinRecord)
	}
	return eventRecords
}

func KinesisRecords(num int) []events.KinesisRecord {
	var records []events.KinesisRecord

	for i := 0; i < num; i++ {
		records = append(records, events.KinesisRecord{
			Data:           []byte(fmt.Sprintf("Data_%d", i)),
			SequenceNumber: fmt.Sprintf("SeqNum_%d", i),
		})
	}
	return records
}

func Records(num int) []*kinesis.Record {
	var records []*kinesis.Record

	for i := 0; i < num; i++ {
		seqNum := fmt.Sprintf("SeqNum_%d", i)
		records = append(records, &kinesis.Record{
			Data:           []byte(fmt.Sprintf("Data_%d", i)),
			SequenceNumber: &seqNum,
		})
	}
	return records
}

func MockKinesisClient(records []events.KinesisRecord) kinesisClientMock {
	return kinesisClientMock{}
}

type kinesisClientMock struct {
	kinesisiface.KinesisAPI
}

func (kinesisClientMock) CreateStream(*kinesis.CreateStreamInput) (*kinesis.CreateStreamOutput, error) {
	return &kinesis.CreateStreamOutput{}, nil
}

func (kinesisClientMock) WaitUntilStreamExists(*kinesis.DescribeStreamInput) error {
	return nil
}

func (kinesisClientMock) PutRecords(params *kinesis.PutRecordsInput) (*kinesis.PutRecordsOutput, error) {
	return nil, nil
}

func (kinesisClientMock) DescribeStream(*kinesis.DescribeStreamInput) (*kinesis.DescribeStreamOutput, error) {
	description := kinesis.StreamDescription{Shards: []*kinesis.Shard{}}
	return &kinesis.DescribeStreamOutput{StreamDescription: &description}, nil
}
