package clients

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesis"
	uuid "github.com/satori/go.uuid"
	"github.com/yunspace/serverless-golang/examples/aws-golang-kinesis/config"
)

func init() {
	KinesisClient = NewKinesisService(config.NewConfig())
}

var KinesisClient *KinesisService

func (ks *KinesisService) CreateStream(streamName string) error {
	_, err := ks.client.CreateStream(&kinesis.CreateStreamInput{
		ShardCount: aws.Int64(1),
		StreamName: &streamName,
	})

	if err != nil {
		log.Printf("Failed to create stream: %v", err.Error())
		return err
	} else {
		log.Printf("Stream created: %s", streamName)
		err := ks.client.WaitUntilStreamExists(&kinesis.DescribeStreamInput{StreamName: &streamName})
		if err != nil {
			log.Println(err.Error())
			return err
		}
	}
	return nil
}

func (ks *KinesisService) PublishRecords(data [][]byte, streamName string) error {
	log.Printf("In function...PublishRecords")
	if len(data) < 1 {
		log.Printf("Not enough data...")
		dataJSON, _ := json.Marshal(data)
		log.Printf("Must pass data to publisher, nothing to publish: %s", dataJSON)
		return nil
	}

	entries := make([]*kinesis.PutRecordsRequestEntry, len(data))

	for i := 0; i < len(entries); i++ {
		entries[i] = &kinesis.PutRecordsRequestEntry{
			Data:         data[i],
			PartitionKey: aws.String(uuid.NewV4().String()),
		}
	}

	params := &kinesis.PutRecordsInput{
		Records:    entries,
		StreamName: &streamName,
	}

	recordsJSON, err := json.Marshal(params.Records)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	log.Printf("Publishing to Kinesis Stream: %s, Data: %s", streamName, recordsJSON)

	putsOutput, err := ks.client.PutRecords(params)

	if err != nil {
		log.Println(err.Error())
		return err
	}
	log.Printf("Published data result: %s", putsOutput)

	return nil
}

func (ks *KinesisService) ConsumeRecords(streamName string, executor KinesisRecordsExecutor) (map[string][][]byte, error) {
	log.Printf("In function...ConsumeRecords")

	stream, err := ks.client.DescribeStream(&kinesis.DescribeStreamInput{StreamName: &streamName})
	if err != nil {
		log.Printf("Failed to get stream description: %v", err.Error())
		return nil, err
	}
	log.Printf("Stream Description: %v\n", stream)

	var returnRecords map[string][][]byte
	for _, shard := range stream.StreamDescription.Shards {
		iteratorOutput, err := ks.client.GetShardIterator(&kinesis.GetShardIteratorInput{
			ShardId:           shard.ShardId,
			ShardIteratorType: aws.String("LATEST"),
			StreamName:        &streamName,
		})

		if err != nil {
			log.Printf("Failed to get iterator: %v", err.Error())
			return nil, err
		}

		records, err := ks.client.GetRecords(&kinesis.GetRecordsInput{
			ShardIterator: iteratorOutput.ShardIterator,
		})

		if err != nil {
			log.Printf("Failed to get records: %v", err.Error())
			return nil, err
		}

		log.Printf("Attaching Records: %v\n", records)
		for _, r := range records.Records {
			var p Payload
			err := json.Unmarshal(r.Data, &p)
			if err != nil {
				log.Printf("json unmarshal error: %v", err)
			} else {
				returnRecords[p.ProxyStreamName] = append(returnRecords[p.ProxyStreamName], r.Data)
			}
		}
	}

	log.Printf("Found all records: %v\n", returnRecords)

	if executor != nil {
		executor(ks, returnRecords)
	}
	return returnRecords, nil
}
