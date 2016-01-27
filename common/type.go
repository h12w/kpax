package common

import (
	"io"
)

type Doer interface {
	Do(Request, Response) error
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

type Cluster interface {
	Coordinator(group string) (Doer, error)
	CoordinatorIsDown(group string)
	Leader(topic string, partition int32) (Doer, error)
	LeaderIsDown(topic string, partition int32)
}
