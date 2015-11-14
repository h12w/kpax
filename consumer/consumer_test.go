package consumer

import (
	"fmt"
	"testing"
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
	for partition := int32(0); partition < 3; partition++ {
		values, err := consumer.Consume("test", partition, 0)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("Partition", partition)
		for _, value := range values {
			fmt.Println(string(value))
		}
		fmt.Println()
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
	consumer, err := New(DefaultConfig(
		"docker:32776",
		"docker:32777",
		"docker:32778",
		"docker:32779",
		"docker:32780",
	))
	if err != nil {
		t.Fatal(err)
	}
	return consumer
}
