package proto

import (
	"encoding"
)

type T interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}
