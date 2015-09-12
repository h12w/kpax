package proto

import (
	"io"
)

type T interface {
	Marshal(io.Writer) error
	Unmarshal([]byte) error
}
