package proto

import (
	"fmt"
	"net"
	"testing"
)

func Test(t *testing.T) {
	conn, err := net.Dial("tcp", "docker:32791")
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
