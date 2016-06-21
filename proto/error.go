package proto

import (
	"errors"
	"fmt"
)

var (
	ErrSizeMismatch = errors.New("proto: size mismatch in response")
	ErrCRCMismatch  = errors.New("proto: CRC mismatch in response")
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
	case ErrOffsetOutOfRange, ErrInvalidMessage, ErrInvalidMessageSize, ErrMessageSizeTooLarge,
		ErrStaleControllerEpochCode, ErrOffsetMetadataTooLargeCode, ErrRecordListTooLargeCode,
		ErrInvalidRequiredAcksCode, ErrIllegalGenerationCode, ErrInconsistentGroupProtocolCode,
		ErrTopicAuthorizationFailedCode, ErrGroupAuthorizationFailedCode, ErrClusterAuthorizationFailedCode:
		return false
	}
	return true
}

func IsNotCoordinator(err error) bool {
	return IsNotLeader(err)
}
