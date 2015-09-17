package broker

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"h12.me/kafka/proto"
)

const (
	kafkaAddr = "docker:32793"
)

func TestMeta(t *testing.T) {
	broker, err := New(&Config{
		Addr:     kafkaAddr,
		QueueLen: 10,
		Timeout:  time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	req := &proto.Request{
		APIKey:        proto.TopicMetadataRequestType,
		APIVersion:    0,
		CorrelationID: 1,
		ClientID:      "abc",
		RequestMessage: &proto.TopicMetadataRequest{
			"test",
		},
	}
	resp := &proto.TopicMetadataResponse{}
	if err := broker.Do(req, resp); err != nil {
		t.Fatal(err)
	}
	fmt.Println(toJSON(resp))
}

func TestConsumeAll(t *testing.T) {
	broker, err := New(&Config{
		Addr:     kafkaAddr,
		QueueLen: 10,
		Timeout:  time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	req := &proto.Request{
		APIKey:        proto.FetchRequestType,
		APIVersion:    0,
		CorrelationID: 1,
		ClientID:      "abc",
		RequestMessage: &proto.FetchRequest{
			ReplicaID:   -1,
			MaxWaitTime: 0,
			MinBytes:    0,
			FetchOffsetInTopics: []proto.FetchOffsetInTopic{
				{
					TopicName: "test",
					FetchOffsetInPartitions: []proto.FetchOffsetInPartition{
						{
							Partition:   0,
							FetchOffset: 0,
							MaxBytes:    10000,
						},
					},
				},
			},
		},
	}
	resp := proto.FetchResponse{}
	if err := broker.Do(req, &resp); err != nil {
		t.Fatal(err)
	}
	fmt.Println(toJSON(resp))
	for _, t := range resp {
		for _, p := range t.FetchMessageSetInPartitions {
			for _, m := range p.MessageSet {
				fmt.Println(string(m.SizedMessage.CRCMessage.Message.Value))
			}
		}
	}
}

func TestOffsetCommit(t *testing.T) {
	broker, err := New(&Config{
		Addr:     kafkaAddr,
		QueueLen: 10,
		Timeout:  time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	tm := time.Now()
	req := &proto.Request{
		APIKey:        proto.OffsetCommitRequestType,
		APIVersion:    1,
		CorrelationID: 1,
		ClientID:      "abc",
		RequestMessage: &proto.OffsetCommitRequestV1{
			ConsumerGroupID:           "test-1",
			ConsumerGroupGenerationID: 1,
			ConsumerID:                "consumer-1",
			OffsetCommitInTopicV1s: []proto.OffsetCommitInTopicV1{
				{
					TopicName: "test",
					OffsetCommitInPartitionV1s: []proto.OffsetCommitInPartitionV1{
						{
							Partition: 0,
							Offset:    1,
							TimeStamp: tm.Unix(),
							Metadata:  fmt.Sprint(tm.Unix()),
						},
					},
				},
			},
		},
	}
	resp := proto.OffsetCommitResponse{}
	if err := broker.Do(req, &resp); err != nil {
		t.Fatal(t)
	}
	fmt.Println(toJSON(resp))
}

func TestOffsetFetch(t *testing.T) {
	broker, err := New(&Config{
		Addr:     kafkaAddr,
		QueueLen: 10,
		Timeout:  time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	req := &proto.Request{
		APIKey:        proto.OffsetFetchRequestType,
		APIVersion:    1,
		CorrelationID: 1,
		ClientID:      "abc",
		RequestMessage: &proto.OffsetFetchRequestV1{
			ConsumerGroup: "test-1",
			PartitionInTopics: []proto.PartitionInTopic{
				{
					TopicName:  "test",
					Partitions: []int32{0, 1, 2},
				},
			},
		},
	}
	resp := proto.OffsetFetchResponse{}
	if err := broker.Do(req, &resp); err != nil {
		t.Fatal(err)
	}
	fmt.Println(toJSON(resp))
}

func TestConsumerMeta(t *testing.T) {
	broker, err := New(&Config{
		Addr:     kafkaAddr,
		QueueLen: 10,
		Timeout:  time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	creq := proto.ConsumerMetadataRequest("test-1")
	req := &proto.Request{
		APIKey:         proto.ConsumerMetadataRequestType,
		APIVersion:     0,
		CorrelationID:  1,
		ClientID:       "abc",
		RequestMessage: &creq,
	}
	resp := proto.ConsumerMetadataResponse{}
	if err := broker.Do(req, &resp); err != nil {
		t.Fatal(err)
	}
	fmt.Println(toJSON(resp))
}

func TestProduce(t *testing.T) {
	broker, err := New(&Config{
		Addr:     kafkaAddr,
		QueueLen: 10,
		Timeout:  time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	tm := time.Now()
	req := &proto.Request{
		APIKey:        proto.ProduceRequestType,
		APIVersion:    0,
		CorrelationID: 1,
		ClientID:      "abc",
		RequestMessage: &proto.ProduceRequest{
			RequiredAcks: 1,
			Timeout:      0,
			MessageSetInTopics: []proto.MessageSetInTopic{
				{
					TopicName: "test",
					MessageSetInPartitions: []proto.MessageSetInPartition{
						{
							Partition: 0,
							MessageSet: []proto.OffsetMessage{
								{
									SizedMessage: proto.SizedMessage{CRCMessage: proto.CRCMessage{
										Message: proto.Message{
											Key:   nil,
											Value: []byte("hello " + tm.Format(time.RFC3339)),
										},
									}}},
							},
						},
					},
				},
			},
		},
	}
	resp := proto.ProduceResponse{}
	if err := broker.Do(req, &resp); err != nil {
		t.Fatal(err)
	}
	fmt.Println(toJSON(resp))
}

func toJSON(v interface{}) string {
	buf, _ := json.MarshalIndent(v, "", "    ")
	return string(buf)
}
