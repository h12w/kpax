package proto

import (
	"fmt"
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
