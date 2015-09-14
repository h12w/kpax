package proto

import (
	"errors"
	"fmt"
)

var (
	ErrUnexpectedEOF = errors.New("proto: unexpected EOF")
	ErrSizeMismatch  = errors.New("proto: size mismatch")
	ErrCRCMismatch   = errors.New("proto: CRC mismatch")
)

type Error struct {
	Msg string
	FilePos
}

func newError(format string, v ...interface{}) *Error {
	return &Error{
		Msg:     fmt.Sprintf(format, v...),
		FilePos: getFilePos(1),
	}
}

func (err *Error) Error() string {
	return fmt.Sprint(err.Msg, err.FilePos)
}
