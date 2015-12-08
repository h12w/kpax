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

func (code ErrorCode) HasError() bool {
	switch code {
	case NoError, ErrReplicaNotAvailable:
		return false
	}
	return true
}

func IsNotLeader(err error) bool {
	switch err {
	case ErrConn, ErrUnknown, ErrNotLeaderForPartition:
		return true
	}
	return false
}

func IsNotCoordinator(err error) bool {
	switch err {
	case ErrConn, ErrUnknown, ErrNotCoordinatorForGroupCode:
		return true
	}
	return false
}
