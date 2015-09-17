package producer

import (
	"testing"
	"time"

	"h12.me/kafka/broker"
	"h12.me/kafka/client"
)

func TestProducer(t *testing.T) {
	producer, err := New(&Config{
		Client: client.Config{
			Brokers: []string{
				"docker:32771",
			},
			BrokerConfig: broker.Config{
				QueueLen: 10,
				Timeout:  time.Second,
			},
			ClientID: "abc",
		},
		RequiredAcks: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := producer.Produce("test", nil, []byte("hello "+time.Now().Format(time.RFC3339))); err != nil {
		t.Fatal(err)
	}
}
