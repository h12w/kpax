package client

import (
	"fmt"
	"testing"

	"h12.me/kafka/broker"
)

func TestClient(t *testing.T) {
	client, err := New(&Config{
		Brokers: []string{
			"docker:32791",
			"docker:32792",
			"docker:32793",
		},
		BrokerConfig: broker.Config{
			QueueLen: 10,
		},
		ClientID: "abc",
	})
	if err != nil {
		t.Fatal(err)
	}
	partitions, err := client.Partitions("test")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(partitions)
}
