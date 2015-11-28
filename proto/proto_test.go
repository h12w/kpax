package proto

import (
	"net"
	"testing"

	"h12.me/realtest/kafka"
)

func TestTopicMetadata(t *testing.T) {
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
	topic, err := k.NewRandomTopic(0)
	if err != nil {
		t.Fatal(err)
	}
	defer k.DeleteTopic(topic)
	reqMsg := &TopicMetadataRequest{
		topic,
	}
	req := Request{
		APIKey:         reqMsg.APIKey(),
		APIVersion:     reqMsg.APIVersion(),
		CorrelationID:  1,
		ClientID:       "abc",
		RequestMessage: reqMsg,
	}
	if err := req.Send(conn); err != nil {
		t.Fatal(t)
	}
	resp := &Response{
		ResponseMessage: &TopicMetadataResponse{},
	}
	if err := resp.Receive(conn); err != nil {
		t.Fatal(resp)
	}
}
