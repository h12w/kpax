package proto

import (
	"errors"
	"fmt"

	"h12.me/wipro"
)

var (
	ErrSizeMismatch = errors.New("proto: size mismatch in response")
	ErrCRCMismatch  = errors.New("proto: CRC mismatch in response")
	ErrConn         = wipro.ErrConn
)

func (code ErrorCode) Error() string {
	if code == -1 {
		return "proto: an unexpected server error"
	} else if code >= 0 && int(code) < len(errTexts) {
		return errTexts[code]
	}
	return fmt.Sprintf("proto: unknown error: %d", code)
}
