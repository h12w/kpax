package proto

import (
	"errors"
	"fmt"
)

var (
	ErrUnexpectedEOF = errors.New("proto: unexpected EOF")
	ErrSizeMismatch  = errors.New("proto: size mismatch in response")
	ErrCRCMismatch   = errors.New("proto: CRC mismatch in response")
	ErrConn          = errors.New("proto: network connection error")
)

func (code ErrorCode) Error() string {
	if code == -1 {
		return "proto: an unexpected server error"
	} else if code >= 0 && int(code) < len(errTexts) {
		return errTexts[code]
	}
	return fmt.Sprintf("proto: unknown error: %d", code)
}
