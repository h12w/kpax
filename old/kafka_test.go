package kafka

import (
	"fmt"
	"testing"

	"h12.me/kafka/connector"
)

func TestIt(t *testing.T) {
	cfg := connector.NewConfig()
	cfg.BrokerList = []string{
		"192.168.59.103:32791",
	}
	conn, err := connector.New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	m, err := conn.GetTopicMetadata([]string{"test"})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%#v", m.TopicsMetadata)
	/*
			producerConfig := &producer.ProducerConfig{
				BatchSize:       1,
				ClientID:        "siesta",
				MaxRequests:     10,
				SendRoutines:    10,
				ReceiveRoutines: 10,
				ReadTimeout:     5 * time.Second,
				WriteTimeout:    5 * time.Second,
				RequiredAcks:    1,
				AckTimeoutMs:    2000,
				Linger:          1 * time.Second,
			}
			p := producer.New(producerConfig, producer.ByteSerializer, producer.StringSerializer, conn)
			errChan := p.Send(&producer.ProducerRecord{Topic: "test1", Value: "hello kafka"})

			select {
			case err := <-errChan:
				t.Fatal(err)
			case <-time.After(5 * time.Second):
				t.Fatal("Could not get produce response within 5 seconds")
			}

			p.Close(1 * time.Second)
		if err := conn.CommitOffset("test1-group1", "test1", 0, 1); err != nil {
			t.Fatal(err)
		}
		r, err := conn.GetOffset("test1-group1", "test1", 0)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(r)
	*/
}
