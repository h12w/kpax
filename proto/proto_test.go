package proto

import (
	"math/rand"
	"net"
	"testing"

	"h12.me/realtest/kafka"
)

func TestTopicMetadata(t *testing.T) {
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
	reqMsg := &TopicMetadataRequest{topic}
	correlationID := rand.Int31()
	req := Request{
		APIKey:         reqMsg.APIKey(),
		APIVersion:     reqMsg.APIVersion(),
		CorrelationID:  correlationID,
		ClientID:       "abc",
		RequestMessage: reqMsg,
	}
	conn, err := net.Dial("tcp", k.AnyBroker())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	if err := req.Send(conn); err != nil {
		t.Fatal(t)
	}
	respMsg := &TopicMetadataResponse{}
	resp := &Response{ResponseMessage: respMsg}
	if err := resp.Receive(conn); err != nil {
		t.Fatal(resp)
	}
	if resp.CorrelationID != correlationID {
		t.Fatalf("correlation id: expect %d but got %d", correlationID, resp.CorrelationID)
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
	k, err := kafka.New()
	if err != nil {
		t.Fatal(err)
	}
	brokers := k.Brokers()
	conn, err := net.Dial("tcp", brokers[0])
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	partitionCount := 2
	topic, err := k.NewRandomTopic(partitionCount)
	if err != nil {
		t.Fatal(err)
	}
	defer k.DeleteTopic(topic)
	reqMsg := &TopicMetadataRequest{topic}
	correlationID := rand.Int31()
	req := Request{
		APIKey:         reqMsg.APIKey(),
		APIVersion:     reqMsg.APIVersion(),
		CorrelationID:  correlationID,
		ClientID:       "abc",
		RequestMessage: reqMsg,
	}
	if err := req.Send(conn); err != nil {
		t.Fatal(t)
	}
	respMsg := &TopicMetadataResponse{}
	resp := &Response{ResponseMessage: respMsg}
	if err := resp.Receive(conn); err != nil {
		t.Fatal(resp)
	}
	if resp.CorrelationID != correlationID {
		t.Fatalf("correlation id: expect %d but got %d", correlationID, resp.CorrelationID)
	}
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
	if len(meta.PartitionMetadatas) != partitionCount {
		t.Fatalf("partition count: expect %d but got %d", partitionCount, len(meta.PartitionMetadatas))
	}
	pMeta := meta.PartitionMetadatas[0]
	t.Log(respMsg.Brokers[int(pMeta.Leader)].Addr())
}
