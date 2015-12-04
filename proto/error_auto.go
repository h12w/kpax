package proto

const (
	NoError                             ErrorCode = 0
	ErrUnknown                          ErrorCode = -1
	ErrOffsetOutOfRange                 ErrorCode = 1
	ErrInvalidMessage                   ErrorCode = 2
	ErrUnknownTopicOrPartition          ErrorCode = 3
	ErrInvalidMessageSize               ErrorCode = 4
	ErrLeaderNotAvailable               ErrorCode = 5
	ErrNotLeaderForPartition            ErrorCode = 6
	ErrRequestTimedOut                  ErrorCode = 7
	ErrBrokerNotAvailable               ErrorCode = 8
	ErrReplicaNotAvailable              ErrorCode = 9
	ErrMessageSizeTooLarge              ErrorCode = 10
	ErrStaleControllerEpochCode         ErrorCode = 11
	ErrOffsetMetadataTooLargeCode       ErrorCode = 12
	ErrGroupLoadInProgressCode          ErrorCode = 14
	ErrGroupCoordinatorNotAvailableCode ErrorCode = 15
	ErrNotCoordinatorForGroupCode       ErrorCode = 16
	ErrInvalidTopicCode                 ErrorCode = 17
	ErrRecordListTooLargeCode           ErrorCode = 18
	ErrNotEnoughReplicasCode            ErrorCode = 19
	ErrNotEnoughReplicasAfterAppendCode ErrorCode = 20
	ErrInvalidRequiredAcksCode          ErrorCode = 21
	ErrIllegalGenerationCode            ErrorCode = 22
	ErrInconsistentGroupProtocolCode    ErrorCode = 23
	ErrInvalidGroupIdCode               ErrorCode = 24
	ErrUnknownMemberIdCode              ErrorCode = 25
	ErrInvalidSessionTimeoutCode        ErrorCode = 26
	ErrRebalanceInProgressCode          ErrorCode = 27
	ErrInvalidCommitOffsetSizeCode      ErrorCode = 28
	ErrTopicAuthorizationFailedCode     ErrorCode = 29
	ErrGroupAuthorizationFailedCode     ErrorCode = 30
	ErrClusterAuthorizationFailedCode   ErrorCode = 31
)

var errTexts = []string{
	0:  "proto(0): no error--it worked!",
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
	14: "proto(14): the broker returns this error code for an offset fetch request if it is still loading offsets (after a leader change for that offsets topic partition), or in response to group membership requests (such as heartbeats) when group metadata is being loaded by the coordinator",
	15: "proto(15): the broker returns this error code for group coordinator requests, offset commits, and most group management requests if the offsets topic has not yet been created, or if the group coordinator is not active",
	16: "proto(16): the broker returns this error code if it receives an offset fetch or commit request for a group that it is not a coordinator for",
	17: "proto(17): for a request which attempts to access an invalid topic (e.g. one which has an illegal name), or if an attempt is made to write to an internal topic (such as the consumer offsets topic)",
	18: "proto(18): if a message batch in a produce request exceeds the maximum configured segment size",
	19: "proto(19): returned from a produce request when the number of in-sync replicas is lower than the configured minimum and requiredAcks is -1",
	20: "proto(20): returned from a produce request when the message was written to the log, but with fewer in-sync replicas than required",
	21: "proto(21): returned from a produce request if the requested requiredAcks is invalid (anything other than -1, 1, or 0)",
	22: "proto(22): returned from group membership requests (such as heartbeats) when the generation id provided in the request is not the current generation",
	23: "proto(23): returned in join group when the member provides a protocol type or set of protocols which is not compatible with the current group",
	24: "proto(24): returned in join group when the groupId is empty or null",
	25: "proto(25): returned from group requests (offset commits/fetches, heartbeats, etc) when the memberId is not in the current generation",
	26: "proto(26): return in join group when the requested session timeout is outside of the allowed range on the broker",
	27: "proto(27): returned in heartbeat requests when the coordinator has begun rebalancing the group. This indicates to the client that it should rejoin the group",
	28: "proto(28): this error indicates that an offset commit was rejected because of oversize metadata",
	29: "proto(29): returned by the broker when the client is not authorized to access the requested topic",
	30: "proto(30): returned by the broker when the client is not authorized to access a particular groupId",
	31: "proto(31): returned by the broker when the client is not authorized to use an inter-broker or administrative API",
}
