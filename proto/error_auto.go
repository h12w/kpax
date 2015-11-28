package proto

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
	0:  "proto(0): no error",
	1:  "proto(1): the requested offset is outside the range of offsets maintained by the server for the given topic/partition",
	2:  "proto(2): this indicates that a message contents does not match its CRC",
	3:  "proto(3): this request is for a topic or partition that does not exist on this broker",
	4:  "proto(4): the message has a negative size",
	5:  "proto(5): this error is thrown if we are in the middle of a leadership election and there is currently no leader for this partition and hence it is unavailable for writes",
	6:  "proto(6): this error is thrown if the client attempts to send messages to a replica that is not the leader for some partition. It indicates that the clients metadata is out of date",
	7:  "proto(7): this error is thrown if the request exceeds the user-specified time limit in the request",
	8:  "proto(8): this is not a client facing error and is used mostly by tools when a broker is not alive",
	9:  "proto(9): if replica is expected on a broker, but is not (this can be safely ignored)",
	10: "proto(10): the server has a configurable maximum message size to avoid unbounded memory allocation. This error is thrown if the client attempt to produce a message larger than this maximum",
	11: "proto(11): internal error code for broker-to-broker communication",
	12: "proto(12): if you specify a string larger than configured maximum for offset metadata",
	13: "proto(13)",
	14: "proto(14): the broker returns this error code for an offset fetch request if it is still loading offsets (after a leader change for that offsets topic partition)",
	15: "proto(15): the broker returns this error code for consumer metadata requests or offset commit requests if the offsets topic has not yet been created",
	16: "proto(16): the broker returns this error code if it receives an offset fetch or commit request for a consumer group that it is not a coordinator for",
}
