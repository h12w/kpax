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
