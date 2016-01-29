package model

import (
	"io"
)

type Broker interface {
	Do(Request, Response) error
	Close()
}

type Cluster interface {
	Coordinator(group string) (Broker, error)
	CoordinatorIsDown(group string)
	Leader(topic string, partition int32) (Broker, error)
	LeaderIsDown(topic string, partition int32)
	Partitions(topic string) ([]int32, error)
}

type Request interface {
	Send(io.Writer) error
	ID() int32
	SetID(int32)
}

type Response interface {
	Receive(io.Reader) error
	ID() int32
}
