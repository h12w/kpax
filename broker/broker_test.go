package broker

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"h12.me/kafka/proto"
)

func TestMeta(t *testing.T) {
	broker, err := New(&Config{
		Addr:         "docker:32791",
		SendQueueLen: 10,
		RecvQueueLen: 10,
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
	resp := &proto.Response{
		ResponseMessage: &proto.TopicMetadataResponse{},
	}
	if err := broker.Do(req, resp); err != nil {
		t.Fatal(t)
	}
	fmt.Println(toJSON(resp.ResponseMessage))
}

func TestProduce(t *testing.T) {
	broker, err := New(&Config{
		Addr:         "docker:32793",
		SendQueueLen: 10,
		RecvQueueLen: 10,
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
	if err := broker.Do(req, &proto.Response{ResponseMessage: &resp}); err != nil {
		t.Fatal(err)
	}
	fmt.Println(toJSON(resp))
}

func toJSON(v interface{}) string {
	buf, _ := json.MarshalIndent(v, "", "    ")
	return string(buf)
}
