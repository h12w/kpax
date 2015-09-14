package client

import (
	"fmt"
	"net"
	"testing"

	"h12.me/kafka/proto"
)

func Test(t *testing.T) {
	conn, err := net.Dial("tcp", "docker:32791")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	broker := NewBroker(&BrokerConfig{
		Conn:         conn,
		SendChanSize: 10,
		RecvChanSize: 10,
	})
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
	fmt.Printf("%#v\n", resp.ResponseMessage)
}
