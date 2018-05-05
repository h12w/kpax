package proto

import "h12.io/wipro"

type RequestOrResponse struct {
	Size int32
	wipro.M
}
type Request struct {
	APIKey        int16
	APIVersion    int16
	CorrelationID int32
	ClientID      string
	RequestMessage
}
type Response struct {
	CorrelationID int32
	ResponseMessage
}
type ResponseMessage wipro.M
type MessageSet []OffsetMessage
type OffsetMessage struct {
	Offset int64
	SizedMessage
}
type SizedMessage struct {
	Size int32
	CRCMessage
}
type CRCMessage struct {
	CRC uint32
	Message
}
type Message struct {
	MagicByte  int8
	Attributes int8
	Key        []byte
	Value      []byte
}
type TopicMetadataRequest []string
type TopicMetadataResponse struct {
	Brokers        []Broker
	TopicMetadatas []TopicMetadata
}
type Broker struct {
	NodeID int32
	Host   string
	Port   int32
}
type TopicMetadata struct {
	ErrorCode
	TopicName          string
	PartitionMetadatas []PartitionMetadata
}
type PartitionMetadata struct {
	ErrorCode
	PartitionID int32
	Leader      int32
	Replicas    []int32
	ISR         []int32
}
type ProduceRequest struct {
	RequiredAcks       int16
	Timeout            int32
	MessageSetInTopics []MessageSetInTopic
}
type MessageSetInTopic struct {
	TopicName              string
	MessageSetInPartitions []MessageSetInPartition
}
type MessageSetInPartition struct {
	Partition int32
	MessageSet
}
type ProduceResponse []OffsetInTopic
type OffsetInTopic struct {
	TopicName          string
	OffsetInPartitions []OffsetInPartition
}
type OffsetInPartition struct {
	Partition int32
	ErrorCode
	Offset int64
}
type FetchRequest struct {
	ReplicaID           int32
	MaxWaitTime         int32
	MinBytes            int32
	FetchOffsetInTopics []FetchOffsetInTopic
}
type FetchOffsetInTopic struct {
	TopicName               string
	FetchOffsetInPartitions []FetchOffsetInPartition
}
type FetchOffsetInPartition struct {
	Partition   int32
	FetchOffset int64
	MaxBytes    int32
}
type FetchResponse []FetchMessageSetInTopic
type FetchMessageSetInTopic struct {
	TopicName                   string
	FetchMessageSetInPartitions []FetchMessageSetInPartition
}
type FetchMessageSetInPartition struct {
	Partition int32
	ErrorCode
	HighwaterMarkOffset int64
	MessageSet
}
type OffsetRequest struct {
	ReplicaID    int32
	TimeInTopics []TimeInTopic
}
type TimeInTopic struct {
	TopicName        string
	TimeInPartitions []TimeInPartition
}
type TimeInPartition struct {
	Partition          int32
	Time               int64
	MaxNumberOfOffsets int32
}
type OffsetResponse []OffsetsInTopic
type OffsetsInTopic struct {
	TopicName           string
	OffsetsInPartitions []OffsetsInPartition
}
type OffsetsInPartition struct {
	Partition int32
	ErrorCode
	Offsets []int64
}
type GroupCoordinatorRequest string
type GroupCoordinatorResponse struct {
	ErrorCode
	Broker
}
type OffsetCommitRequestV0 struct {
	ConsumerGroupID        string
	OffsetCommitInTopicV0s []OffsetCommitInTopicV0
}
type OffsetCommitInTopicV0 struct {
	TopicName                  string
	OffsetCommitInPartitionV0s []OffsetCommitInPartitionV0
}
type OffsetCommitInPartitionV0 struct {
	Partition int32
	Offset    int64
	Metadata  string
}
type OffsetCommitRequestV1 struct {
	ConsumerGroupID           string
	ConsumerGroupGenerationID int32
	ConsumerID                string
	OffsetCommitInTopicV1s    []OffsetCommitInTopicV1
}
type OffsetCommitInTopicV1 struct {
	TopicName                  string
	OffsetCommitInPartitionV1s []OffsetCommitInPartitionV1
}
type OffsetCommitInPartitionV1 struct {
	Partition int32
	Offset    int64
	TimeStamp int64
	Metadata  string
}
type OffsetCommitRequestV2 struct {
	ConsumerGroup             string
	ConsumerGroupGenerationID int32
	ConsumerID                string
	RetentionTime             int64
	OffsetCommitInTopicV2s    []OffsetCommitInTopicV2
}
type OffsetCommitInTopicV2 struct {
	TopicName                  string
	OffsetCommitInPartitionV2s []OffsetCommitInPartitionV2
}
type OffsetCommitInPartitionV2 struct {
	Partition int32
	Offset    int64
	Metadata  string
}
type OffsetCommitResponse []ErrorInTopic
type ErrorInTopic struct {
	TopicName         string
	ErrorInPartitions []ErrorInPartition
}
type ErrorInPartition struct {
	Partition int32
	ErrorCode
}
type OffsetFetchRequestV0 struct {
	ConsumerGroup     string
	PartitionInTopics []PartitionInTopic
}
type PartitionInTopic struct {
	TopicName  string
	Partitions []int32
}
type OffsetFetchRequestV1 struct {
	ConsumerGroup     string
	PartitionInTopics []PartitionInTopic
}
type OffsetFetchResponse []OffsetMetadataInTopic
type OffsetMetadataInTopic struct {
	TopicName                  string
	OffsetMetadataInPartitions []OffsetMetadataInPartition
}
type OffsetMetadataInPartition struct {
	Partition int32
	Offset    int64
	Metadata  string
	ErrorCode
}
type JoinGroupRequest struct {
	GroupID        string
	SessionTimeout int32
	MemberID       string
	ProtocolType   string
	GroupProtocols
}
type GroupProtocols []GroupProtocol
type GroupProtocol struct {
	ProtocolName string
	ProtocolMetadata
}
type ProtocolMetadata struct {
	Version int16
	Subscription
	UserData []byte
}
type Subscription []string
type JoinGroupResponse struct {
	ErrorCode
	GenerationID      int32
	GroupProtocolName string
	LeaderID          string
	MemberID          string
	MemberWithMetas
}
type MemberWithMetas []MemberWithMeta
type MemberWithMeta struct {
	MemberID       string
	MemberMetadata []byte
}
type SyncGroupRequest struct {
	GroupID      string
	GenerationID int32
	MemberID     string
	GroupAssignments
}
type GroupAssignments []GroupAssignment
type GroupAssignment struct {
	MemberID string
	MemberAssignment
}
type MemberAssignment struct {
	Version int16
	PartitionAssignments
}
type PartitionAssignments []PartitionAssignment
type PartitionAssignment struct {
	Topic      string
	Partitions []int32
}
type SyncGroupResponse struct {
	ErrorCode
	MemberAssignment
}
type HeartbeatRequest struct {
	GroupID      string
	GenerationID int32
	MemberID     string
}
type HeartbeatResponse ErrorCode
type LeaveGroupRequest struct {
	GroupID  string
	MemberID string
}
type LeaveGroupResponse ErrorCode
type ListGroupsRequest struct {
}
type ListGroupsResponse struct {
	ErrorCode
	Groups
}
type Groups []Group
type Group struct {
	GroupID      string
	ProtocolType string
}
type DescribeGroupsRequest []string
type DescribeGroupsResponse []GroupDescription
type GroupDescription struct {
	ErrorCode
	GroupID      string
	State        string
	ProtocolType string
	Protocol     string
	Members
}
type Members []Member
type Member struct {
	MemberID       string
	ClientID       string
	ClientHost     string
	MemberMetadata []byte
	MemberAssignment
}
type ErrorCode int16
