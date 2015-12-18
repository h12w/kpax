package cluster

import (
	"fmt"
	"testing"
)

func TestClient(t *testing.T) {
	client, err := New(DefaultConfig("docker:32791"))
	if err != nil {
		t.Fatal(err)
	}
	partitions, err := client.Partitions("test")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(partitions)
}
