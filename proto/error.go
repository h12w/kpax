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

//type ErrorCode int

const (
	ErrUnknown ErrorCode = iota - 1
	NoError
	ErrOffsetOutOfRange
	ErrInvalidMessage
	ErrUnknownTopicOrPartition
	ErrInvalidMessageSize
	ErrLeaderNotAvailable
	ErrNotLeaderForPartition
	ErrRequestTimedOut
	ErrBrokerNotAvailable
	ErrReplicaNotAvailable
	ErrMessageSizeTooLarge
	ErrStaleControllerEpochCode
	ErrOffsetMetadataTooLargeCode
	_
	ErrOffsetsLoadInProgressCode
	ErrConsumerCoordinatorNotAvailableCode
	ErrNotCoordinatorForConsumerCode
)

var errTexts = []string{
	0:  "proto: no error",
	1:  "proto: the requested offset is outside the range of offsets maintained by the server for the given topic/partition",
	2:  "proto: this indicates that a message contents does not match its CRC",
	3:  "proto: this request is for a topic or partition that does not exist on this broker",
	4:  "proto: the message has a negative size",
	5:  "proto: this error is thrown if we are in the middle of a leadership election and there is currently no leader for this partition and hence it is unavailable for writes",
	6:  "proto: this error is thrown if the client attempts to send messages to a replica that is not the leader for some partition. It indicates that the clients metadata is out of date",
	7:  "proto: this error is thrown if the request exceeds the user-specified time limit in the request",
	8:  "proto: this is not a client facing error and is used mostly by tools when a broker is not alive",
	9:  "proto: if replica is expected on a broker, but is not (this can be safely ignored)",
	10: "proto: the server has a configurable maximum message size to avoid unbounded memory allocation. This error is thrown if the client attempt to produce a message larger than this maximum",
	11: "proto: internal error code for broker-to-broker communication",
	12: "proto: if you specify a string larger than configured maximum for offset metadata",
	// 13
	14: "proto: the broker returns this error code for an offset fetch request if it is still loading offsets (after a leader change for that offsets topic partition)",
	15: "proto: the broker returns this error code for consumer metadata requests or offset commit requests if the offsets topic has not yet been created",
	16: "proto: the broker returns this error code if it receives an offset fetch or commit request for a consumer group that it is not a coordinator for",
}

func (code ErrorCode) Error() string {
	if code == -1 {
		return "proto: an unexpected server error"
	} else if code >= 0 && int(code) < len(errTexts) {
		return errTexts[code]
	}
	return fmt.Sprintf("proto: unknown error: %d", code)
}
