package proto

import (
	"errors"
	"io"
)

var (
	errNilUnmarshaler = errors.New("nil unmarshaler")
)

type T interface {
	Marshal(io.Writer) error
	Unmarshal(io.Reader) error
}
