package client

import (
	"h12.me/kafka/broker"
	"h12.me/kafka/proto"
)

type Config struct {
	Brokers []string
}

type C struct {
	pm partitionMap
}

func New() *C {
	return &C{
		pm: make(partitionMap),
	}
}

type (
	topicPartition struct {
		Topic     string
		Partition int32
	}
	partitionMap map[topicPartition]*broker.B
)

func (c *C) Produce(msg *proto.Message) error {
	return nil
}
