package proto

import (
	"errors"
	"fmt"
)

var (
	ErrUnexpectedEOF = errors.New("proto: unexpected EOF")
	ErrSizeMismatch  = errors.New("proto: size mismatch")
	ErrCRCMismatch   = errors.New("proto: CRC mismatch")
	ErrConn          = errors.New("proto: network connection error")
)

var errTexts = []string{
	0:  "no error",
	1:  "the requested offset is outside the range of offsets maintained by the server for the given topic/partition",
	2:  "this indicates that a message contents does not match its CRC",
	3:  "this request is for a topic or partition that does not exist on this broker",
	4:  "the message has a negative size",
	5:  "this error is thrown if we are in the middle of a leadership election and there is currently no leader for this partition and hence it is unavailable for writes",
	6:  "this error is thrown if the client attempts to send messages to a replica that is not the leader for some partition. It indicates that the clients metadata is out of date",
	7:  "this error is thrown if the request exceeds the user-specified time limit in the request",
	8:  "this is not a client facing error and is used mostly by tools when a broker is not alive",
	9:  "if replica is expected on a broker, but is not (this can be safely ignored)",
	10: "the server has a configurable maximum message size to avoid unbounded memory allocation. This error is thrown if the client attempt to produce a message larger than this maximum",
	11: "internal error code for broker-to-broker communication",
	12: "if you specify a string larger than configured maximum for offset metadata",
	13: "13",
	14: "the broker returns this error code for an offset fetch request if it is still loading offsets (after a leader change for that offsets topic partition)",
	15: "the broker returns this error code for consumer metadata requests or offset commit requests if the offsets topic has not yet been created",
	16: "the broker returns this error code if it receives an offset fetch or commit request for a consumer group that it is not a coordinator for",
}

func ErrorText(code int16) string {
	if code == -1 {
		return "an unexpected server error"
	} else if code >= 0 && int(code) < len(errTexts) {
		return errTexts[code]
	}
	return fmt.Sprintf("unknown error: %d", code)
}
