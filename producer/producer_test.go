package producer

import (
	"testing"
	"time"
)

func TestProducer(t *testing.T) {
	producer, err := New(DefaultConfig(
		"docker:32776",
		"docker:32777",
		"docker:32778",
		"docker:32779",
		"docker:32780",
	))
	if err != nil {
		t.Fatal(err)
	}
	if err := producer.Produce("test", nil, []byte("hello "+time.Now().Format(time.RFC3339))); err != nil {
		t.Fatal(err)
	}
}
