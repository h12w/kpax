package proto

import (
	"fmt"
	"net"
	"testing"
	"time"

	"h12.me/kafka/broker"
	"h12.me/realtest/kafka"
	"h12.me/wipro"
)

var (
	New           = broker.New
	DefaultConfig = broker.DefaultConfig
)

func TestTopicMetadata(t *testing.T) {
	t.Parallel()
	k, err := kafka.New()
	if err != nil {
		t.Fatal(err)
	}
	partitionCount := 2
	topic, err := k.NewRandomTopic(partitionCount)
	if err != nil {
		t.Fatal(err)
	}
	defer k.DeleteTopic(topic)
	respMsg := getTopicMetadata(t, k, topic)
	meta := &respMsg.TopicMetadatas[0]
	if len(meta.PartitionMetadatas) != partitionCount {
		t.Fatalf("partition count: expect %d but got %d", partitionCount, len(meta.PartitionMetadatas))
	}
	for _, pMeta := range meta.PartitionMetadatas {
		if pMeta.ErrorCode != NoError {
			t.Fatal(pMeta.ErrorCode)
		}
	}
}

func TestProduceFetch(t *testing.T) {
	t.Parallel()
	k, err := kafka.New()
	if err != nil {
		t.Fatal(err)
	}
	partitionCount := 2
	partition := int32(1)
	topic, err := k.NewRandomTopic(partitionCount)
	if err != nil {
		t.Fatal(err)
	}
	defer k.DeleteTopic(topic)
	leaderAddr := getLeader(t, k, topic, partition)
	b := New(DefaultConfig().WithAddr(leaderAddr))
	defer b.Close()
	key, value := "test key", "test value"
	produceMessage(t, b, topic, partition, key, value)
	messages := fetchMessage(t, b, topic, partition, 0)
	if len(messages) != 1 {
		t.Fatalf("expect 1 message but got %v", messages)
	}
	if m := messages[0]; m[0] != key || m[1] != value {
		t.Fatalf("expect [%s %s] but got %v", key, value, m)
	}
}

func TestProduceSnappy(t *testing.T) {
	t.Parallel()
	k, err := kafka.New()
	if err != nil {
		t.Fatal(err)
	}
	partitionCount := 2
	partition := int32(1)
	topic := "topic1"
	err = k.NewTopic(topic, partitionCount)
	if err != nil {
		t.Fatal(err)
	}
	defer k.DeleteTopic(topic)
	leaderAddr := getLeader(t, k, topic, partition)

	b := New(DefaultConfig().WithAddr(leaderAddr))
	defer b.Close()
	var w wipro.Writer
	ms := MessageSet{
		{SizedMessage: SizedMessage{CRCMessage: CRCMessage{Message: Message{
			Attributes: 0,
			Key:        nil,
			Value:      []byte("hello"),
		}}}},
	}
	ms.Marshal(&w)
	compressedValue := encodeSnappy(w.B[4:])
	fmt.Println(w.B)
	if err := (&Payload{
		Topic:     topic,
		Partition: partition,
		MessageSet: MessageSet{
			{
				SizedMessage: SizedMessage{CRCMessage: CRCMessage{Message: Message{
					Attributes: 2,
					Key:        nil,
					Value:      compressedValue,
				}}}},
		},
		RequiredAcks: AckLocal,
		AckTimeout:   10 * time.Second,
	}).DoProduce(b); err != nil {
		t.Fatal(err)
	}
}

func TestGroupCoordinator(t *testing.T) {
	t.Parallel()
	k, err := kafka.New()
	if err != nil {
		t.Fatal(err)
	}
	group := kafka.RandomGroup()
	coord := getCoord(t, k, group)
	fmt.Println(group, coord)
}

func getTopicMetadata(t *testing.T, k *kafka.Cluster, topic string) *TopicMetadataResponse {
	b := New(DefaultConfig().WithAddr(k.AnyBroker()))
	defer b.Close()
	respMsg, err := Metadata{topic}.Fetch(b)
	if err != nil {
		t.Fatal(err)
	}
	brokers := k.Brokers()
	for i := range brokers {
		if respMsg.Brokers[i].Addr() != brokers[i] {
			t.Fatalf("broker: expect %s but got %s", brokers[i], respMsg.Brokers[i].Addr())
		}
	}
	if len(respMsg.TopicMetadatas) != 1 {
		t.Fatalf("len(TopicMetadatas): expect 1 but got %d", len(respMsg.TopicMetadatas))
	}
	meta := &respMsg.TopicMetadatas[0]
	if meta.ErrorCode != NoError {
		t.Fatal(meta.ErrorCode)
	}
	if meta.TopicName != topic {
		t.Fatalf("topic: expect %s but got %s", topic, meta.TopicName)
	}
	return respMsg
}

func getLeader(t *testing.T, k *kafka.Cluster, topic string, partitionID int32) string {
	metaResp := getTopicMetadata(t, k, topic)
	meta := &metaResp.TopicMetadatas[0]
	leaderAddr := ""
	for _, partition := range meta.PartitionMetadatas {
		if partition.PartitionID == partitionID {
			for _, broker := range metaResp.Brokers {
				if broker.NodeID == partition.Leader {
					leaderAddr = broker.Addr()
				}
			}
		}
	}
	if leaderAddr == "" {
		t.Fatalf("fail to find leader in topic %s partition %d", topic, partitionID)
	}
	return leaderAddr
}

func getCoord(t *testing.T, k *kafka.Cluster, group string) string {
	reqMsg := GroupCoordinatorRequest(group)
	req := &Request{
		RequestMessage: &reqMsg,
	}
	respMsg := &GroupCoordinatorResponse{}
	resp := &Response{ResponseMessage: respMsg}
	conn, err := net.Dial("tcp", k.AnyBroker())
	if err != nil {
		t.Fatal(err)
	}
	sendReceive(t, conn, req, resp)
	if respMsg.ErrorCode.HasError() {
		t.Fatal(respMsg.ErrorCode)
	}
	return respMsg.Broker.Addr()
}

func produceMessage(t *testing.T, b *broker.B, topic string, partition int32, key, value string) {
	if err := (&Payload{
		Topic:     topic,
		Partition: partition,
		MessageSet: MessageSet{
			{
				SizedMessage: SizedMessage{CRCMessage: CRCMessage{Message: Message{
					Key:   []byte(key),
					Value: []byte(value),
				}}},
			},
		},
		RequiredAcks: AckLocal,
		AckTimeout:   10 * time.Second,
	}).DoProduce(b); err != nil {
		t.Fatal(err)
	}
}

func fetchMessage(t *testing.T, b *broker.B, topic string, partition int32, offset int64) [][2]string {
	req := FetchRequest{
		ReplicaID:   -1,
		MaxWaitTime: int32(time.Second / time.Millisecond),
		MinBytes:    1,
		FetchOffsetInTopics: []FetchOffsetInTopic{
			{
				TopicName: topic,
				FetchOffsetInPartitions: []FetchOffsetInPartition{
					{
						Partition:   partition,
						FetchOffset: offset,
						MaxBytes:    1024 * 1024,
					},
				},
			},
		},
	}
	resp := FetchResponse{}
	if err := (client{clientID, b}).Do(&req, &resp); err != nil {
		t.Fatal(err)
	}
	var result [][2]string
	for _, t := range resp {
		if t.TopicName == topic {
			for _, p := range t.FetchMessageSetInPartitions {
				if p.Partition == partition {
					if p.ErrorCode == NoError {
						for _, msg := range p.MessageSet {
							m := &msg.SizedMessage.CRCMessage.Message
							result = append(result, [2]string{string(m.Key), string(m.Value)})
						}
					}
				}
			}
		}
	}
	return result
}

func sendReceive(t *testing.T, conn net.Conn, req *Request, resp *Response) {
	if err := req.Send(conn); err != nil {
		t.Fatal(t)
	}
	if err := resp.Receive(conn); err != nil {
		t.Fatal(err)
	}
	if resp.CorrelationID != req.CorrelationID {
		t.Fatalf("correlation id: expect %d but got %d", req.CorrelationID, resp.CorrelationID)
	}
}
