package proto

const (
	ProduceRequestType          = 0
	FetchRequestType            = 1
	OffsetRequestType           = 2
	TopicMetadataRequestType    = 3
	OffsetCommitRequestType     = 8
	OffsetFetchRequestType      = 9
	ConsumerMetadataRequestType = 10
)

const (
	AckNone  = 0
	AckLocal = 1
	AckAll   = -1
)
