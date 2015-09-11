package proto

import (
	"bytes"
)

func (t *RequestOrResponse) Marshal() []byte {
	d0 := []byte{byte(t.Size >> 24), byte(t.Size >> 16), byte(t.Size >> 8), byte(t.Size)}
	d1 := t.Marshal()
	return bytes.Join([][]byte{d0, d1}, nil)
}

func (t *Request) Marshal() []byte {
	d0 := []byte{byte(t.APIKey >> 8), byte(t.APIKey)}
	d1 := []byte{byte(t.APIVersion >> 8), byte(t.APIVersion)}
	d2 := []byte{byte(t.CorrelationID >> 24), byte(t.CorrelationID >> 16), byte(t.CorrelationID >> 8), byte(t.CorrelationID)}
	d3Len := strlen(t.ClientID)
	d3 := append([]byte{byte(d3Len >> 8), byte(d3Len)}, []byte(t.ClientID)...)
	d4 := t.Marshal()
	return bytes.Join([][]byte{d0, d1, d2, d3, d4}, nil)
}

// type RequestMessage

func (t *Response) Marshal() []byte {
	d0 := []byte{byte(t.CorrelationID >> 24), byte(t.CorrelationID >> 16), byte(t.CorrelationID >> 8), byte(t.CorrelationID)}
	d1 := t.Marshal()
	return bytes.Join([][]byte{d0, d1}, nil)
}

// type ResponseMessage

// type MessageSet

func (t *SizedMessage) Marshal() []byte {
	d0 := []byte{byte(t.Offset >> 56), byte(t.Offset >> 48), byte(t.Offset >> 32), byte(t.Offset >> 24), byte(t.Offset >> 16), byte(t.Offset >> 8), byte(t.Offset)}
	d1 := []byte{byte(t.MessageSize >> 24), byte(t.MessageSize >> 16), byte(t.MessageSize >> 8), byte(t.MessageSize)}
	d2 := t.Marshal()
	return bytes.Join([][]byte{d0, d1, d2}, nil)
}

func (t *Message) Marshal() []byte {
	d0 := []byte{byte(t.CRC >> 24), byte(t.CRC >> 16), byte(t.CRC >> 8), byte(t.CRC)}
	d1 := []byte{byte(t.MagicByte)}
	d2 := []byte{byte(t.Attributes)}
	d3 := t.Marshal()
	d4 := t.Marshal()
	return bytes.Join([][]byte{d0, d1, d2, d3, d4}, nil)
}

// type TopicMetadataRequest

func (t *MetadataResponse) Marshal() []byte {
	// field []
	// field []
	return bytes.Join([][]byte{d0, d1}, nil)
}

func (t *Broker) Marshal() []byte {
	d0 := []byte{byte(t.NodeID >> 24), byte(t.NodeID >> 16), byte(t.NodeID >> 8), byte(t.NodeID)}
	d1Len := strlen(t.Host)
	d1 := append([]byte{byte(d1Len >> 8), byte(d1Len)}, []byte(t.Host)...)
	d2 := []byte{byte(t.Port >> 24), byte(t.Port >> 16), byte(t.Port >> 8), byte(t.Port)}
	return bytes.Join([][]byte{d0, d1, d2}, nil)
}

func (t *TopicMetadata) Marshal() []byte {
	d0 := []byte{byte(t.TopicErrorCode >> 8), byte(t.TopicErrorCode)}
	d1Len := strlen(t.TopicName)
	d1 := append([]byte{byte(d1Len >> 8), byte(d1Len)}, []byte(t.TopicName)...)
	// field []
	return bytes.Join([][]byte{d0, d1, d2}, nil)
}

func (t *PartitionMetadata) Marshal() []byte {
	d0 := []byte{byte(t.PartitionErrorCode >> 8), byte(t.PartitionErrorCode)}
	d1 := []byte{byte(t.PartitionID >> 24), byte(t.PartitionID >> 16), byte(t.PartitionID >> 8), byte(t.PartitionID)}
	d2 := []byte{byte(t.Leader >> 24), byte(t.Leader >> 16), byte(t.Leader >> 8), byte(t.Leader)}
	d3 := []byte{byte(t.Replicas >> 24), byte(t.Replicas >> 16), byte(t.Replicas >> 8), byte(t.Replicas)}
	d4 := []byte{byte(t.ISR >> 24), byte(t.ISR >> 16), byte(t.ISR >> 8), byte(t.ISR)}
	return bytes.Join([][]byte{d0, d1, d2, d3, d4}, nil)
}

func (t *ProduceRequest) Marshal() []byte {
	d0 := []byte{byte(t.RequiredAcks >> 8), byte(t.RequiredAcks)}
	d1 := []byte{byte(t.Timeout >> 24), byte(t.Timeout >> 16), byte(t.Timeout >> 8), byte(t.Timeout)}
	// field []
	return bytes.Join([][]byte{d0, d1, d2}, nil)
}

func (t *MessageSetInTopic) Marshal() []byte {
	d0Len := strlen(t.TopicName)
	d0 := append([]byte{byte(d0Len >> 8), byte(d0Len)}, []byte(t.TopicName)...)
	// field []
	return bytes.Join([][]byte{d0, d1}, nil)
}

func (t *MessageSetInPartition) Marshal() []byte {
	d0 := []byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}
	d1 := []byte{byte(t.MessageSetSize >> 24), byte(t.MessageSetSize >> 16), byte(t.MessageSetSize >> 8), byte(t.MessageSetSize)}
	d2 := t.Marshal()
	return bytes.Join([][]byte{d0, d1, d2}, nil)
}

// type ProduceResponse

func (t *OffsetInTopic) Marshal() []byte {
	d0Len := strlen(t.TopicName)
	d0 := append([]byte{byte(d0Len >> 8), byte(d0Len)}, []byte(t.TopicName)...)
	d1 := t.Marshal()
	return bytes.Join([][]byte{d0, d1}, nil)
}

// type OffsetInPartitions

func (t *OffsetInPartition) Marshal() []byte {
	d0 := []byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}
	d1 := []byte{byte(t.ErrorCode >> 8), byte(t.ErrorCode)}
	d2 := []byte{byte(t.Offset >> 56), byte(t.Offset >> 48), byte(t.Offset >> 32), byte(t.Offset >> 24), byte(t.Offset >> 16), byte(t.Offset >> 8), byte(t.Offset)}
	return bytes.Join([][]byte{d0, d1, d2}, nil)
}

func (t *FetchRequest) Marshal() []byte {
	d0 := []byte{byte(t.ReplicaID >> 24), byte(t.ReplicaID >> 16), byte(t.ReplicaID >> 8), byte(t.ReplicaID)}
	d1 := []byte{byte(t.MaxWaitTime >> 24), byte(t.MaxWaitTime >> 16), byte(t.MaxWaitTime >> 8), byte(t.MaxWaitTime)}
	d2 := []byte{byte(t.MinBytes >> 24), byte(t.MinBytes >> 16), byte(t.MinBytes >> 8), byte(t.MinBytes)}
	// field []
	return bytes.Join([][]byte{d0, d1, d2, d3}, nil)
}

func (t *FetchOffsetInTopic) Marshal() []byte {
	d0Len := strlen(t.TopicName)
	d0 := append([]byte{byte(d0Len >> 8), byte(d0Len)}, []byte(t.TopicName)...)
	// field []
	return bytes.Join([][]byte{d0, d1}, nil)
}

func (t *FetchOffsetInPartition) Marshal() []byte {
	d0 := []byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}
	d1 := []byte{byte(t.FetchOffset >> 56), byte(t.FetchOffset >> 48), byte(t.FetchOffset >> 32), byte(t.FetchOffset >> 24), byte(t.FetchOffset >> 16), byte(t.FetchOffset >> 8), byte(t.FetchOffset)}
	d2 := []byte{byte(t.MaxBytes >> 24), byte(t.MaxBytes >> 16), byte(t.MaxBytes >> 8), byte(t.MaxBytes)}
	return bytes.Join([][]byte{d0, d1, d2}, nil)
}

// type FetchResponse

func (t *FetchMessageSetInTopic) Marshal() []byte {
	d0Len := strlen(t.TopicName)
	d0 := append([]byte{byte(d0Len >> 8), byte(d0Len)}, []byte(t.TopicName)...)
	// field []
	return bytes.Join([][]byte{d0, d1}, nil)
}

func (t *FetchMessageSetInPartition) Marshal() []byte {
	d0 := []byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}
	d1 := []byte{byte(t.ErrorCode >> 8), byte(t.ErrorCode)}
	d2 := []byte{byte(t.HighwaterMarkOffset >> 56), byte(t.HighwaterMarkOffset >> 48), byte(t.HighwaterMarkOffset >> 32), byte(t.HighwaterMarkOffset >> 24), byte(t.HighwaterMarkOffset >> 16), byte(t.HighwaterMarkOffset >> 8), byte(t.HighwaterMarkOffset)}
	d3 := []byte{byte(t.MessageSetSize >> 24), byte(t.MessageSetSize >> 16), byte(t.MessageSetSize >> 8), byte(t.MessageSetSize)}
	d4 := t.Marshal()
	return bytes.Join([][]byte{d0, d1, d2, d3, d4}, nil)
}

func (t *OffsetRequest) Marshal() []byte {
	d0 := []byte{byte(t.ReplicaID >> 24), byte(t.ReplicaID >> 16), byte(t.ReplicaID >> 8), byte(t.ReplicaID)}
	// field []
	return bytes.Join([][]byte{d0, d1}, nil)
}

func (t *TimeInTopic) Marshal() []byte {
	d0Len := strlen(t.TopicName)
	d0 := append([]byte{byte(d0Len >> 8), byte(d0Len)}, []byte(t.TopicName)...)
	// field []
	return bytes.Join([][]byte{d0, d1}, nil)
}

func (t *TimeInPartition) Marshal() []byte {
	d0 := []byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}
	d1 := []byte{byte(t.Time >> 56), byte(t.Time >> 48), byte(t.Time >> 32), byte(t.Time >> 24), byte(t.Time >> 16), byte(t.Time >> 8), byte(t.Time)}
	d2 := []byte{byte(t.MaxNumberOfOffsets >> 24), byte(t.MaxNumberOfOffsets >> 16), byte(t.MaxNumberOfOffsets >> 8), byte(t.MaxNumberOfOffsets)}
	return bytes.Join([][]byte{d0, d1, d2}, nil)
}

// type OffsetResponse

func (t *OffsetsInTopic) Marshal() []byte {
	d0Len := strlen(t.TopicName)
	d0 := append([]byte{byte(d0Len >> 8), byte(d0Len)}, []byte(t.TopicName)...)
	// field []
	return bytes.Join([][]byte{d0, d1}, nil)
}

func (t *OffsetsInPartition) Marshal() []byte {
	d0 := []byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}
	d1 := []byte{byte(t.ErrorCode >> 8), byte(t.ErrorCode)}
	// field []
	return bytes.Join([][]byte{d0, d1, d2}, nil)
}

// type ConsumerMetadataRequest

func (t *ConsumerMetadataResponse) Marshal() []byte {
	d0 := []byte{byte(t.ErrorCode >> 8), byte(t.ErrorCode)}
	d1 := []byte{byte(t.CoordinatorID >> 24), byte(t.CoordinatorID >> 16), byte(t.CoordinatorID >> 8), byte(t.CoordinatorID)}
	d2Len := strlen(t.CoordinatorHost)
	d2 := append([]byte{byte(d2Len >> 8), byte(d2Len)}, []byte(t.CoordinatorHost)...)
	d3 := []byte{byte(t.CoordinatorPort >> 24), byte(t.CoordinatorPort >> 16), byte(t.CoordinatorPort >> 8), byte(t.CoordinatorPort)}
	return bytes.Join([][]byte{d0, d1, d2, d3}, nil)
}

func (t *OffsetCommitRequestV0) Marshal() []byte {
	d0Len := strlen(t.ConsumerGroupID)
	d0 := append([]byte{byte(d0Len >> 8), byte(d0Len)}, []byte(t.ConsumerGroupID)...)
	// field []
	return bytes.Join([][]byte{d0, d1}, nil)
}

func (t *OffsetCommitInTopicV0) Marshal() []byte {
	d0Len := strlen(t.TopicName)
	d0 := append([]byte{byte(d0Len >> 8), byte(d0Len)}, []byte(t.TopicName)...)
	// field []
	return bytes.Join([][]byte{d0, d1}, nil)
}

func (t *OffsetCommitInPartitionV0) Marshal() []byte {
	d0 := []byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}
	d1 := []byte{byte(t.Offset >> 56), byte(t.Offset >> 48), byte(t.Offset >> 32), byte(t.Offset >> 24), byte(t.Offset >> 16), byte(t.Offset >> 8), byte(t.Offset)}
	d2Len := strlen(t.Metadata)
	d2 := append([]byte{byte(d2Len >> 8), byte(d2Len)}, []byte(t.Metadata)...)
	return bytes.Join([][]byte{d0, d1, d2}, nil)
}

func (t *OffsetCommitRequestV1) Marshal() []byte {
	d0Len := strlen(t.ConsumerGroupID)
	d0 := append([]byte{byte(d0Len >> 8), byte(d0Len)}, []byte(t.ConsumerGroupID)...)
	d1 := []byte{byte(t.ConsumerGroupGenerationID >> 24), byte(t.ConsumerGroupGenerationID >> 16), byte(t.ConsumerGroupGenerationID >> 8), byte(t.ConsumerGroupGenerationID)}
	d2Len := strlen(t.ConsumerID)
	d2 := append([]byte{byte(d2Len >> 8), byte(d2Len)}, []byte(t.ConsumerID)...)
	// field []
	return bytes.Join([][]byte{d0, d1, d2, d3}, nil)
}

func (t *OffsetCommitInTopicV1) Marshal() []byte {
	d0Len := strlen(t.TopicName)
	d0 := append([]byte{byte(d0Len >> 8), byte(d0Len)}, []byte(t.TopicName)...)
	// field []
	return bytes.Join([][]byte{d0, d1}, nil)
}

func (t *OffsetCommitInPartitionV1) Marshal() []byte {
	d0 := []byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}
	d1 := []byte{byte(t.Offset >> 56), byte(t.Offset >> 48), byte(t.Offset >> 32), byte(t.Offset >> 24), byte(t.Offset >> 16), byte(t.Offset >> 8), byte(t.Offset)}
	d2 := []byte{byte(t.TimeStamp >> 56), byte(t.TimeStamp >> 48), byte(t.TimeStamp >> 32), byte(t.TimeStamp >> 24), byte(t.TimeStamp >> 16), byte(t.TimeStamp >> 8), byte(t.TimeStamp)}
	d3Len := strlen(t.Metadata)
	d3 := append([]byte{byte(d3Len >> 8), byte(d3Len)}, []byte(t.Metadata)...)
	return bytes.Join([][]byte{d0, d1, d2, d3}, nil)
}

func (t *OffsetCommitRequestV2) Marshal() []byte {
	d0Len := strlen(t.ConsumerGroup)
	d0 := append([]byte{byte(d0Len >> 8), byte(d0Len)}, []byte(t.ConsumerGroup)...)
	d1 := []byte{byte(t.ConsumerGroupGenerationID >> 24), byte(t.ConsumerGroupGenerationID >> 16), byte(t.ConsumerGroupGenerationID >> 8), byte(t.ConsumerGroupGenerationID)}
	d2Len := strlen(t.ConsumerID)
	d2 := append([]byte{byte(d2Len >> 8), byte(d2Len)}, []byte(t.ConsumerID)...)
	d3 := []byte{byte(t.RetentionTime >> 56), byte(t.RetentionTime >> 48), byte(t.RetentionTime >> 32), byte(t.RetentionTime >> 24), byte(t.RetentionTime >> 16), byte(t.RetentionTime >> 8), byte(t.RetentionTime)}
	// field []
	return bytes.Join([][]byte{d0, d1, d2, d3, d4}, nil)
}

func (t *OffsetCommitInTopicV2) Marshal() []byte {
	d0Len := strlen(t.TopicName)
	d0 := append([]byte{byte(d0Len >> 8), byte(d0Len)}, []byte(t.TopicName)...)
	// field []
	return bytes.Join([][]byte{d0, d1}, nil)
}

func (t *OffsetCommitInPartitionV2) Marshal() []byte {
	d0 := []byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}
	d1 := []byte{byte(t.Offset >> 56), byte(t.Offset >> 48), byte(t.Offset >> 32), byte(t.Offset >> 24), byte(t.Offset >> 16), byte(t.Offset >> 8), byte(t.Offset)}
	d2Len := strlen(t.Metadata)
	d2 := append([]byte{byte(d2Len >> 8), byte(d2Len)}, []byte(t.Metadata)...)
	return bytes.Join([][]byte{d0, d1, d2}, nil)
}

// type OffsetCommitResponse

func (t *ErrorInTopic) Marshal() []byte {
	d0Len := strlen(t.TopicName)
	d0 := append([]byte{byte(d0Len >> 8), byte(d0Len)}, []byte(t.TopicName)...)
	// field []
	return bytes.Join([][]byte{d0, d1}, nil)
}

func (t *ErrorInPartition) Marshal() []byte {
	d0 := []byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}
	d1 := []byte{byte(t.ErrorCode >> 8), byte(t.ErrorCode)}
	return bytes.Join([][]byte{d0, d1}, nil)
}

func (t *OffsetFetchRequest) Marshal() []byte {
	d0Len := strlen(t.ConsumerGroup)
	d0 := append([]byte{byte(d0Len >> 8), byte(d0Len)}, []byte(t.ConsumerGroup)...)
	// field []
	return bytes.Join([][]byte{d0, d1}, nil)
}

func (t *PartitionInTopic) Marshal() []byte {
	d0Len := strlen(t.TopicName)
	d0 := append([]byte{byte(d0Len >> 8), byte(d0Len)}, []byte(t.TopicName)...)
	// field []
	return bytes.Join([][]byte{d0, d1}, nil)
}

// type OffsetFetchResponse

func (t *OffsetMetadataInTopic) Marshal() []byte {
	d0Len := strlen(t.TopicName)
	d0 := append([]byte{byte(d0Len >> 8), byte(d0Len)}, []byte(t.TopicName)...)
	// field []
	return bytes.Join([][]byte{d0, d1}, nil)
}

func (t *OffsetMetadataInPartition) Marshal() []byte {
	d0 := []byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}
	d1 := []byte{byte(t.Offset >> 56), byte(t.Offset >> 48), byte(t.Offset >> 32), byte(t.Offset >> 24), byte(t.Offset >> 16), byte(t.Offset >> 8), byte(t.Offset)}
	d2Len := strlen(t.Metadata)
	d2 := append([]byte{byte(d2Len >> 8), byte(d2Len)}, []byte(t.Metadata)...)
	d3 := []byte{byte(t.ErrorCode >> 8), byte(t.ErrorCode)}
	return bytes.Join([][]byte{d0, d1, d2, d3}, nil)
}
