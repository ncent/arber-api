package clients

import (
	"testing"

	goblin "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

const streamName = "defaultTestStream"

func Test(t *testing.T) {
	g := goblin.Goblin(t)
	kinesisRecords := KinesisRecords(10)
	kinesisService := MockKinesisService(MockKinesisClient(kinesisRecords))
	var dataRecords [][]byte
	for _, record := range kinesisRecords {
		dataRecords = append(dataRecords, record.Data)
	}

	//special hook for gomega
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Kinesis Service", func() {
		g.Describe("CreateStream", func() {
			g.It("Should not throw an error", func() {
				err := kinesisService.CreateStream(streamName)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		g.Describe("PublishRecords", func() {
			g.It("Should add records into the stream", func() {

				err := kinesisService.PublishRecords(dataRecords, streamName)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		g.Describe("ConsumeRecords", func() {
			g.It("Should read records from the stream", func() {
				_, err := kinesisService.ConsumeRecords(streamName, nil)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
}
