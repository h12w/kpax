package proto

import (
	"errors"
)

var (
	ErrUnexpectedEOF = errors.New("proto: unexpected EOF")
	ErrSizeMismatch  = errors.New("proto: size mismatch")
	ErrCRCMismatch   = errors.New("proto: CRC mismatch")
	ErrConn          = errors.New("proto: network connection error")
)
