package producer

import (
	"testing"
	"time"
)

func TestProducer(t *testing.T) {
	producer, err := New(DefaultConfig(
		"docker:32791",
		"docker:32792",
		"docker:32793",
		"docker:32794",
		"docker:32795",
	))
	if err != nil {
		t.Fatal(err)
	}
	if err := producer.Produce("test", nil, []byte("hello "+time.Now().Format(time.RFC3339))); err != nil {
		t.Fatal(err)
	}
}
