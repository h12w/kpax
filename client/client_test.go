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
	req := proto.Request{
		APIKey:        proto.TopicMetadataRequestType,
		APIVersion:    0,
		CorrelationID: 1,
		ClientID:      "abc",
		RequestMessage: &proto.TopicMetadataRequest{
			"test",
		},
	}
	if err := req.Send(conn); err != nil {
		t.Fatal(t)
	}
	resp := &proto.Response{
		ResponseMessage: &proto.TopicMetadataResponse{},
	}
	if err := resp.Receive(conn); err != nil {
		t.Fatal(resp)
	}
	fmt.Printf("%#v\n", resp.ResponseMessage)
}
