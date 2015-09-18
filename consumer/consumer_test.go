package consumer

import (
	"fmt"
	"testing"
	"time"

	"h12.me/kafka/broker"
	"h12.me/kafka/client"
)

func TestGetOffset(t *testing.T) {
	consumer := getConsumer(t)
	offset, err := consumer.Offset("test", 0, "test-consumergroup-b")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("get offset: ", offset)
}

func TestConsumeAll(t *testing.T) {
	consumer := getConsumer(t)
	values, err := consumer.Consume("test", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range values {
		fmt.Println(string(value))
	}
}

func TestCommitOffset(t *testing.T) {
	consumer := getConsumer(t)
	err := consumer.Commit("test", 0, "test-consumergroup-b", 2)
	if err != nil {
		t.Fatal(err)
	}
}

func getConsumer(t *testing.T) *C {
	consumer, err := New(&Config{
		Client: client.Config{
			Brokers: []string{
				"docker:32791",
			},
			BrokerConfig: broker.Config{
				QueueLen: 10,
				Timeout:  time.Second,
			},
			ClientID: "abc",
		},
		MinBytes:        0,
		MaxBytes:        9999,
		OffsetRetention: 7 * 24 * time.Hour,
	})
	if err != nil {
		t.Fatal(err)
	}
	return consumer
}
