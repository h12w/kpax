package proto

import (
	"fmt"
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
	req := Request{
		APIKey:        TopicMetadataRequestType,
		APIVersion:    0,
		CorrelationID: 1,
		ClientID:      "abc",
		RequestMessage: &TopicMetadataRequest{
			"test",
		},
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
	fmt.Printf("%#v\n", resp.ResponseMessage)
}
