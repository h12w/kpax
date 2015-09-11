package proto

func (t *RequestOrResponse) MarshalBinary() (data []byte, err error) {
	d0 := []byte{byte(t.Size >> 24), byte(t.Size >> 16), byte(t.Size >> 8), byte(t.Size)}
}

func (t *Request) MarshalBinary() (data []byte, err error) {
	d0 := []byte{byte(t.APIKey >> 8), byte(t.APIKey)}
	d1 := []byte{byte(t.APIVersion >> 8), byte(t.APIVersion)}
	d2 := []byte{byte(t.CorrelationID >> 24), byte(t.CorrelationID >> 16), byte(t.CorrelationID >> 8), byte(t.CorrelationID)}
}

// type RequestMessage

func (t *Response) MarshalBinary() (data []byte, err error) {
	d0 := []byte{byte(t.CorrelationID >> 24), byte(t.CorrelationID >> 16), byte(t.CorrelationID >> 8), byte(t.CorrelationID)}
}

// type ResponseMessage

// type MessageSet

func (t *SizedMessage) MarshalBinary() (data []byte, err error) {
	d0 := []byte{byte(t.Offset >> 56), byte(t.Offset >> 48), byte(t.Offset >> 32), byte(t.Offset >> 24), byte(t.Offset >> 16), byte(t.Offset >> 8), byte(t.Offset)}
	d1 := []byte{byte(t.MessageSize >> 24), byte(t.MessageSize >> 16), byte(t.MessageSize >> 8), byte(t.MessageSize)}
}

func (t *Message) MarshalBinary() (data []byte, err error) {
	d0 := []byte{byte(t.CRC >> 24), byte(t.CRC >> 16), byte(t.CRC >> 8), byte(t.CRC)}
	d1 := []byte{byte(t.MagicByte)}
	d2 := []byte{byte(t.Attributes)}
}

// type TopicMetadataRequest

func (t *MetadataResponse) MarshalBinary() (data []byte, err error) {
	// field []
	// field []
}

func (t *Broker) MarshalBinary() (data []byte, err error) {
	d0 := []byte{byte(t.NodeID >> 24), byte(t.NodeID >> 16), byte(t.NodeID >> 8), byte(t.NodeID)}
	d2 := []byte{byte(t.Port >> 24), byte(t.Port >> 16), byte(t.Port >> 8), byte(t.Port)}
}

func (t *TopicMetadata) MarshalBinary() (data []byte, err error) {
	d0 := []byte{byte(t.TopicErrorCode >> 8), byte(t.TopicErrorCode)}
	// field []
}

func (t *PartitionMetadata) MarshalBinary() (data []byte, err error) {
	d0 := []byte{byte(t.PartitionErrorCode >> 8), byte(t.PartitionErrorCode)}
	d1 := []byte{byte(t.PartitionID >> 24), byte(t.PartitionID >> 16), byte(t.PartitionID >> 8), byte(t.PartitionID)}
	d2 := []byte{byte(t.Leader >> 24), byte(t.Leader >> 16), byte(t.Leader >> 8), byte(t.Leader)}
	d3 := []byte{byte(t.Replicas >> 24), byte(t.Replicas >> 16), byte(t.Replicas >> 8), byte(t.Replicas)}
	d4 := []byte{byte(t.ISR >> 24), byte(t.ISR >> 16), byte(t.ISR >> 8), byte(t.ISR)}
}

func (t *ProduceRequest) MarshalBinary() (data []byte, err error) {
	d0 := []byte{byte(t.RequiredAcks >> 8), byte(t.RequiredAcks)}
	d1 := []byte{byte(t.Timeout >> 24), byte(t.Timeout >> 16), byte(t.Timeout >> 8), byte(t.Timeout)}
	// field []
}

func (t *MessageSetInTopic) MarshalBinary() (data []byte, err error) {
	// field []
}

func (t *MessageSetInPartition) MarshalBinary() (data []byte, err error) {
	d0 := []byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}
	d1 := []byte{byte(t.MessageSetSize >> 24), byte(t.MessageSetSize >> 16), byte(t.MessageSetSize >> 8), byte(t.MessageSetSize)}
}

// type ProduceResponse

func (t *OffsetInTopic) MarshalBinary() (data []byte, err error) {
}

// type OffsetInPartitions

func (t *OffsetInPartition) MarshalBinary() (data []byte, err error) {
	d0 := []byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}
	d1 := []byte{byte(t.ErrorCode >> 8), byte(t.ErrorCode)}
	d2 := []byte{byte(t.Offset >> 56), byte(t.Offset >> 48), byte(t.Offset >> 32), byte(t.Offset >> 24), byte(t.Offset >> 16), byte(t.Offset >> 8), byte(t.Offset)}
}

func (t *FetchRequest) MarshalBinary() (data []byte, err error) {
	d0 := []byte{byte(t.ReplicaID >> 24), byte(t.ReplicaID >> 16), byte(t.ReplicaID >> 8), byte(t.ReplicaID)}
	d1 := []byte{byte(t.MaxWaitTime >> 24), byte(t.MaxWaitTime >> 16), byte(t.MaxWaitTime >> 8), byte(t.MaxWaitTime)}
	d2 := []byte{byte(t.MinBytes >> 24), byte(t.MinBytes >> 16), byte(t.MinBytes >> 8), byte(t.MinBytes)}
	// field []
}

func (t *FetchOffsetInTopic) MarshalBinary() (data []byte, err error) {
	// field []
}

func (t *FetchOffsetInPartition) MarshalBinary() (data []byte, err error) {
	d0 := []byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}
	d1 := []byte{byte(t.FetchOffset >> 56), byte(t.FetchOffset >> 48), byte(t.FetchOffset >> 32), byte(t.FetchOffset >> 24), byte(t.FetchOffset >> 16), byte(t.FetchOffset >> 8), byte(t.FetchOffset)}
	d2 := []byte{byte(t.MaxBytes >> 24), byte(t.MaxBytes >> 16), byte(t.MaxBytes >> 8), byte(t.MaxBytes)}
}

// type FetchResponse

func (t *FetchMessageSetInTopic) MarshalBinary() (data []byte, err error) {
	// field []
}

func (t *FetchMessageSetInPartition) MarshalBinary() (data []byte, err error) {
	d0 := []byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}
	d1 := []byte{byte(t.ErrorCode >> 8), byte(t.ErrorCode)}
	d2 := []byte{byte(t.HighwaterMarkOffset >> 56), byte(t.HighwaterMarkOffset >> 48), byte(t.HighwaterMarkOffset >> 32), byte(t.HighwaterMarkOffset >> 24), byte(t.HighwaterMarkOffset >> 16), byte(t.HighwaterMarkOffset >> 8), byte(t.HighwaterMarkOffset)}
	d3 := []byte{byte(t.MessageSetSize >> 24), byte(t.MessageSetSize >> 16), byte(t.MessageSetSize >> 8), byte(t.MessageSetSize)}
}

func (t *OffsetRequest) MarshalBinary() (data []byte, err error) {
	d0 := []byte{byte(t.ReplicaID >> 24), byte(t.ReplicaID >> 16), byte(t.ReplicaID >> 8), byte(t.ReplicaID)}
	// field []
}

func (t *TimeInTopic) MarshalBinary() (data []byte, err error) {
	// field []
}

func (t *TimeInPartition) MarshalBinary() (data []byte, err error) {
	d0 := []byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}
	d1 := []byte{byte(t.Time >> 56), byte(t.Time >> 48), byte(t.Time >> 32), byte(t.Time >> 24), byte(t.Time >> 16), byte(t.Time >> 8), byte(t.Time)}
	d2 := []byte{byte(t.MaxNumberOfOffsets >> 24), byte(t.MaxNumberOfOffsets >> 16), byte(t.MaxNumberOfOffsets >> 8), byte(t.MaxNumberOfOffsets)}
}

// type OffsetResponse

func (t *OffsetsInTopic) MarshalBinary() (data []byte, err error) {
	// field []
}

func (t *OffsetsInPartition) MarshalBinary() (data []byte, err error) {
	d0 := []byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}
	d1 := []byte{byte(t.ErrorCode >> 8), byte(t.ErrorCode)}
	// field []
}

// type ConsumerMetadataRequest

func (t *ConsumerMetadataResponse) MarshalBinary() (data []byte, err error) {
	d0 := []byte{byte(t.ErrorCode >> 8), byte(t.ErrorCode)}
	d1 := []byte{byte(t.CoordinatorID >> 24), byte(t.CoordinatorID >> 16), byte(t.CoordinatorID >> 8), byte(t.CoordinatorID)}
	d3 := []byte{byte(t.CoordinatorPort >> 24), byte(t.CoordinatorPort >> 16), byte(t.CoordinatorPort >> 8), byte(t.CoordinatorPort)}
}

func (t *OffsetCommitRequestV0) MarshalBinary() (data []byte, err error) {
	// field []
}

func (t *OffsetCommitInTopicV0) MarshalBinary() (data []byte, err error) {
	// field []
}

func (t *OffsetCommitInPartitionV0) MarshalBinary() (data []byte, err error) {
	d0 := []byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}
	d1 := []byte{byte(t.Offset >> 56), byte(t.Offset >> 48), byte(t.Offset >> 32), byte(t.Offset >> 24), byte(t.Offset >> 16), byte(t.Offset >> 8), byte(t.Offset)}
}

func (t *OffsetCommitRequestV1) MarshalBinary() (data []byte, err error) {
	d1 := []byte{byte(t.ConsumerGroupGenerationID >> 24), byte(t.ConsumerGroupGenerationID >> 16), byte(t.ConsumerGroupGenerationID >> 8), byte(t.ConsumerGroupGenerationID)}
	// field []
}

func (t *OffsetCommitInTopicV1) MarshalBinary() (data []byte, err error) {
	// field []
}

func (t *OffsetCommitInPartitionV1) MarshalBinary() (data []byte, err error) {
	d0 := []byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}
	d1 := []byte{byte(t.Offset >> 56), byte(t.Offset >> 48), byte(t.Offset >> 32), byte(t.Offset >> 24), byte(t.Offset >> 16), byte(t.Offset >> 8), byte(t.Offset)}
	d2 := []byte{byte(t.TimeStamp >> 56), byte(t.TimeStamp >> 48), byte(t.TimeStamp >> 32), byte(t.TimeStamp >> 24), byte(t.TimeStamp >> 16), byte(t.TimeStamp >> 8), byte(t.TimeStamp)}
}

func (t *OffsetCommitRequestV2) MarshalBinary() (data []byte, err error) {
	d1 := []byte{byte(t.ConsumerGroupGenerationID >> 24), byte(t.ConsumerGroupGenerationID >> 16), byte(t.ConsumerGroupGenerationID >> 8), byte(t.ConsumerGroupGenerationID)}
	d3 := []byte{byte(t.RetentionTime >> 56), byte(t.RetentionTime >> 48), byte(t.RetentionTime >> 32), byte(t.RetentionTime >> 24), byte(t.RetentionTime >> 16), byte(t.RetentionTime >> 8), byte(t.RetentionTime)}
	// field []
}

func (t *OffsetCommitInTopicV2) MarshalBinary() (data []byte, err error) {
	// field []
}

func (t *OffsetCommitInPartitionV2) MarshalBinary() (data []byte, err error) {
	d0 := []byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}
	d1 := []byte{byte(t.Offset >> 56), byte(t.Offset >> 48), byte(t.Offset >> 32), byte(t.Offset >> 24), byte(t.Offset >> 16), byte(t.Offset >> 8), byte(t.Offset)}
}

// type OffsetCommitResponse

func (t *ErrorInTopic) MarshalBinary() (data []byte, err error) {
	// field []
}

func (t *ErrorInPartition) MarshalBinary() (data []byte, err error) {
	d0 := []byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}
	d1 := []byte{byte(t.ErrorCode >> 8), byte(t.ErrorCode)}
}

func (t *OffsetFetchRequest) MarshalBinary() (data []byte, err error) {
	// field []
}

func (t *PartitionInTopic) MarshalBinary() (data []byte, err error) {
	// field []
}

// type OffsetFetchResponse

func (t *OffsetMetadataInTopic) MarshalBinary() (data []byte, err error) {
	// field []
}

func (t *OffsetMetadataInPartition) MarshalBinary() (data []byte, err error) {
	d0 := []byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}
	d1 := []byte{byte(t.Offset >> 56), byte(t.Offset >> 48), byte(t.Offset >> 32), byte(t.Offset >> 24), byte(t.Offset >> 16), byte(t.Offset >> 8), byte(t.Offset)}
	d3 := []byte{byte(t.ErrorCode >> 8), byte(t.ErrorCode)}
}
