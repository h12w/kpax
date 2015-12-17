package broker

/*
import (
	"fmt"
	"testing"
	"time"
)

const (
	kafkaAddr = "docker:32793"
)

func newTestBroker() *B {
	cfg := DefaultConfig()
	cfg.Addr = kafkaAddr
	return New(cfg)
}

func TestMeta(t *testing.T) {
	broker := newTestBroker()
	req := &Request{
		CorrelationID: 1,
		ClientID:      "abc",
		RequestMessage: &TopicMetadataRequest{
			"test",
		},
	}
	resp := &TopicMetadataResponse{}
	if err := broker.Do(req, resp); err != nil {
		t.Fatal(err)
	}
	fmt.Println(toJSON(resp))
}

func TestConsumeAll(t *testing.T) {
	broker := newTestBroker()
	req := &Request{
		CorrelationID: 1,
		ClientID:      "abc",
		RequestMessage: &FetchRequest{
			ReplicaID:   -1,
			MaxWaitTime: 0,
			MinBytes:    0,
			FetchOffsetInTopics: []FetchOffsetInTopic{
				{
					TopicName: "test",
					FetchOffsetInPartitions: []FetchOffsetInPartition{
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
	resp := FetchResponse{}
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
	broker := New(&Config{
		Addr:     kafkaAddr,
		QueueLen: 10,
		Timeout:  time.Second,
	})
	tm := time.Now()
	req := &Request{
		CorrelationID: 1,
		ClientID:      "abc",
		RequestMessage: &OffsetCommitRequestV1{
			ConsumerGroupID:           "test-1",
			ConsumerGroupGenerationID: 1,
			ConsumerID:                "consumer-1",
			OffsetCommitInTopicV1s: []OffsetCommitInTopicV1{
				{
					TopicName: "test",
					OffsetCommitInPartitionV1s: []OffsetCommitInPartitionV1{
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
	resp := OffsetCommitResponse{}
	if err := broker.Do(req, &resp); err != nil {
		t.Fatal(t)
	}
	fmt.Println(toJSON(resp))
}

func TestOffsetFetch(t *testing.T) {
	broker := newTestBroker()
	req := &Request{
		CorrelationID: 1,
		ClientID:      "abc",
		RequestMessage: &OffsetFetchRequestV1{
			ConsumerGroup: "test-1",
			PartitionInTopics: []PartitionInTopic{
				{
					TopicName:  "test",
					Partitions: []int32{0, 1, 2},
				},
			},
		},
	}
	resp := OffsetFetchResponse{}
	if err := broker.Do(req, &resp); err != nil {
		t.Fatal(err)
	}
	fmt.Println(toJSON(resp))
}

func TestConsumerMeta(t *testing.T) {
	broker := newTestBroker()
	creq := GroupCoordinatorRequest("test-1")
	req := &Request{
		CorrelationID:  1,
		ClientID:       "abc",
		RequestMessage: &creq,
	}
	resp := GroupCoordinatorResponse{}
	if err := broker.Do(req, &resp); err != nil {
		t.Fatal(err)
	}
	fmt.Println(toJSON(resp))
}

func TestProduce(t *testing.T) {
	broker := newTestBroker()
	tm := time.Now()
	req := &Request{
		CorrelationID: 1,
		ClientID:      "abc",
		RequestMessage: &ProduceRequest{
			RequiredAcks: 1,
			Timeout:      0,
			MessageSetInTopics: []MessageSetInTopic{
				{
					TopicName: "test",
					MessageSetInPartitions: []MessageSetInPartition{
						{
							Partition: 0,
							MessageSet: []OffsetMessage{
								{
									SizedMessage: SizedMessage{CRCMessage: CRCMessage{
										Message: Message{
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
	resp := ProduceResponse{}
	if err := broker.Do(req, &resp); err != nil {
		t.Fatal(err)
	}
	fmt.Println(toJSON(resp))
}

*/
