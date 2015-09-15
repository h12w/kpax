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
				"docker:32791",
				"docker:32792",
				"docker:32793",
			},
			BrokerConfig: broker.Config{
				SendQueueLen: 10,
				RecvQueueLen: 10,
			},
			ClientID: "abc",
		},
		RequiredAcks: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	tm := time.Now()
	if err := producer.Produce("test", nil, []byte("hello "+tm.Format(time.RFC3339))); err != nil {
		t.Fatal(err)
	}
}
