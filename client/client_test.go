package client

import (
	"fmt"
	"io"
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
	req := proto.RequestOrResponse{
		T: &proto.Request{
			APIKey:        proto.TopicMetadataRequestType,
			APIVersion:    0,
			CorrelationID: 1,
			ClientID:      "abc",
			RequestMessage: &proto.TopicMetadataRequest{
				"test",
			},
		},
	}
	w := &proto.Writer{}
	req.Marshal(w)
	if _, err := conn.Write(w.B); err != nil {
		t.Fatal(t)
	}
	resp := proto.RequestOrResponse{
		T: &proto.Response{
			ResponseMessage: &proto.TopicMetadataResponse{},
		},
	}
	r := &proto.Reader{B: make([]byte, 4)}
	if _, err := conn.Read(r.B); err != nil {
		t.Fatal(err)
	}
	size := int(r.ReadInt32())
	fmt.Println("response size: ", size, r.B)
	b := r.B
	r.B = make([]byte, 4+size)
	copy(r.B[:4], b)
	r.Offset = 0
	if _, err := io.ReadAtLeast(conn, r.B[4:], size); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%d: %#v\n", len(r.B), string(r.B))
	resp.Unmarshal(r)
	fmt.Printf("%#v\n", resp.T.(*proto.Response).ResponseMessage)
	if r.Err != nil {
		t.Fatal(r)
	}
}
