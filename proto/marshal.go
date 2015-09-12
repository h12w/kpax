package proto

import (
	"io"
)

func (t *RequestOrResponse) Marshal(w io.Writer) error {
	if _, err := w.Write([]byte{byte(t.Size >> 24), byte(t.Size >> 16), byte(t.Size >> 8), byte(t.Size)}); err != nil {
		return err
	}
	if err := t.T.Marshal(w); err != nil {
		return err
	}
	return nil
}

func (t *Request) Marshal(w io.Writer) error {
	if _, err := w.Write([]byte{byte(t.APIKey >> 8), byte(t.APIKey)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.APIVersion >> 8), byte(t.APIVersion)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.CorrelationID >> 24), byte(t.CorrelationID >> 16), byte(t.CorrelationID >> 8), byte(t.CorrelationID)}); err != nil {
		return err
	}
	{
		l := int16(len(t.ClientID))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(t.ClientID)); err != nil {
			return err
		}
	}
	if err := t.RequestMessage.Marshal(w); err != nil {
		return err
	}
	return nil
}

// RequestMessage T

func (t *Response) Marshal(w io.Writer) error {
	if _, err := w.Write([]byte{byte(t.CorrelationID >> 24), byte(t.CorrelationID >> 16), byte(t.CorrelationID >> 8), byte(t.CorrelationID)}); err != nil {
		return err
	}
	if err := t.ResponseMessage.Marshal(w); err != nil {
		return err
	}
	return nil
}

// ResponseMessage T

func (t MessageSet) Marshal(w io.Writer) error {
	{
		l := int32(len(t))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t {
			if err := t[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *SizedMessage) Marshal(w io.Writer) error {
	if _, err := w.Write([]byte{byte(t.Offset >> 56), byte(t.Offset >> 48), byte(t.Offset >> 32), byte(t.Offset >> 24), byte(t.Offset >> 16), byte(t.Offset >> 8), byte(t.Offset)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.MessageSize >> 24), byte(t.MessageSize >> 16), byte(t.MessageSize >> 8), byte(t.MessageSize)}); err != nil {
		return err
	}
	if err := t.Message.Marshal(w); err != nil {
		return err
	}
	return nil
}

func (t *Message) Marshal(w io.Writer) error {
	if _, err := w.Write([]byte{byte(t.CRC >> 24), byte(t.CRC >> 16), byte(t.CRC >> 8), byte(t.CRC)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.MagicByte)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.Attributes)}); err != nil {
		return err
	}
	{
		l := int32(len(t.Key))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write(t.Key); err != nil {
			return err
		}
	}
	{
		l := int32(len(t.Value))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write(t.Value); err != nil {
			return err
		}
	}
	return nil
}

func (t TopicMetadataRequest) Marshal(w io.Writer) error {
	{
		l := int32(len(t))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t {
			{
				l := int16(len(t[i]))
				if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
					return err
				}
				if _, err := w.Write([]byte(t[i])); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (t *MetadataResponse) Marshal(w io.Writer) error {
	{
		l := int32(len(t.Brokers))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t.Brokers {
			if err := t.Brokers[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	{
		l := int32(len(t.TopicMetadatas))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t.TopicMetadatas {
			if err := t.TopicMetadatas[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *Broker) Marshal(w io.Writer) error {
	if _, err := w.Write([]byte{byte(t.NodeID >> 24), byte(t.NodeID >> 16), byte(t.NodeID >> 8), byte(t.NodeID)}); err != nil {
		return err
	}
	{
		l := int16(len(t.Host))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(t.Host)); err != nil {
			return err
		}
	}
	if _, err := w.Write([]byte{byte(t.Port >> 24), byte(t.Port >> 16), byte(t.Port >> 8), byte(t.Port)}); err != nil {
		return err
	}
	return nil
}

func (t *TopicMetadata) Marshal(w io.Writer) error {
	if _, err := w.Write([]byte{byte(t.TopicErrorCode >> 8), byte(t.TopicErrorCode)}); err != nil {
		return err
	}
	{
		l := int16(len(t.TopicName))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(t.TopicName)); err != nil {
			return err
		}
	}
	{
		l := int32(len(t.PartitionMetadatas))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t.PartitionMetadatas {
			if err := t.PartitionMetadatas[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *PartitionMetadata) Marshal(w io.Writer) error {
	if _, err := w.Write([]byte{byte(t.PartitionErrorCode >> 8), byte(t.PartitionErrorCode)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.PartitionID >> 24), byte(t.PartitionID >> 16), byte(t.PartitionID >> 8), byte(t.PartitionID)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.Leader >> 24), byte(t.Leader >> 16), byte(t.Leader >> 8), byte(t.Leader)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.Replicas >> 24), byte(t.Replicas >> 16), byte(t.Replicas >> 8), byte(t.Replicas)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.ISR >> 24), byte(t.ISR >> 16), byte(t.ISR >> 8), byte(t.ISR)}); err != nil {
		return err
	}
	return nil
}

func (t *ProduceRequest) Marshal(w io.Writer) error {
	if _, err := w.Write([]byte{byte(t.RequiredAcks >> 8), byte(t.RequiredAcks)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.Timeout >> 24), byte(t.Timeout >> 16), byte(t.Timeout >> 8), byte(t.Timeout)}); err != nil {
		return err
	}
	{
		l := int32(len(t.MessageSetInTopics))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t.MessageSetInTopics {
			if err := t.MessageSetInTopics[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *MessageSetInTopic) Marshal(w io.Writer) error {
	{
		l := int16(len(t.TopicName))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(t.TopicName)); err != nil {
			return err
		}
	}
	{
		l := int32(len(t.MessageSetInPartitions))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t.MessageSetInPartitions {
			if err := t.MessageSetInPartitions[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *MessageSetInPartition) Marshal(w io.Writer) error {
	if _, err := w.Write([]byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.MessageSetSize >> 24), byte(t.MessageSetSize >> 16), byte(t.MessageSetSize >> 8), byte(t.MessageSetSize)}); err != nil {
		return err
	}
	if err := t.MessageSet.Marshal(w); err != nil {
		return err
	}
	return nil
}

func (t ProduceResponse) Marshal(w io.Writer) error {
	{
		l := int32(len(t))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t {
			if err := t[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *OffsetInTopic) Marshal(w io.Writer) error {
	{
		l := int16(len(t.TopicName))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(t.TopicName)); err != nil {
			return err
		}
	}
	if err := t.OffsetInPartitions.Marshal(w); err != nil {
		return err
	}
	return nil
}

func (t OffsetInPartitions) Marshal(w io.Writer) error {
	{
		l := int32(len(t))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t {
			if err := t[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *OffsetInPartition) Marshal(w io.Writer) error {
	if _, err := w.Write([]byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.ErrorCode >> 8), byte(t.ErrorCode)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.Offset >> 56), byte(t.Offset >> 48), byte(t.Offset >> 32), byte(t.Offset >> 24), byte(t.Offset >> 16), byte(t.Offset >> 8), byte(t.Offset)}); err != nil {
		return err
	}
	return nil
}

func (t *FetchRequest) Marshal(w io.Writer) error {
	if _, err := w.Write([]byte{byte(t.ReplicaID >> 24), byte(t.ReplicaID >> 16), byte(t.ReplicaID >> 8), byte(t.ReplicaID)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.MaxWaitTime >> 24), byte(t.MaxWaitTime >> 16), byte(t.MaxWaitTime >> 8), byte(t.MaxWaitTime)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.MinBytes >> 24), byte(t.MinBytes >> 16), byte(t.MinBytes >> 8), byte(t.MinBytes)}); err != nil {
		return err
	}
	{
		l := int32(len(t.FetchOffsetInTopics))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t.FetchOffsetInTopics {
			if err := t.FetchOffsetInTopics[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *FetchOffsetInTopic) Marshal(w io.Writer) error {
	{
		l := int16(len(t.TopicName))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(t.TopicName)); err != nil {
			return err
		}
	}
	{
		l := int32(len(t.FetchOffsetInPartitions))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t.FetchOffsetInPartitions {
			if err := t.FetchOffsetInPartitions[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *FetchOffsetInPartition) Marshal(w io.Writer) error {
	if _, err := w.Write([]byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.FetchOffset >> 56), byte(t.FetchOffset >> 48), byte(t.FetchOffset >> 32), byte(t.FetchOffset >> 24), byte(t.FetchOffset >> 16), byte(t.FetchOffset >> 8), byte(t.FetchOffset)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.MaxBytes >> 24), byte(t.MaxBytes >> 16), byte(t.MaxBytes >> 8), byte(t.MaxBytes)}); err != nil {
		return err
	}
	return nil
}

func (t FetchResponse) Marshal(w io.Writer) error {
	{
		l := int32(len(t))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t {
			if err := t[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *FetchMessageSetInTopic) Marshal(w io.Writer) error {
	{
		l := int16(len(t.TopicName))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(t.TopicName)); err != nil {
			return err
		}
	}
	{
		l := int32(len(t.FetchMessageSetInPartitions))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t.FetchMessageSetInPartitions {
			if err := t.FetchMessageSetInPartitions[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *FetchMessageSetInPartition) Marshal(w io.Writer) error {
	if _, err := w.Write([]byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.ErrorCode >> 8), byte(t.ErrorCode)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.HighwaterMarkOffset >> 56), byte(t.HighwaterMarkOffset >> 48), byte(t.HighwaterMarkOffset >> 32), byte(t.HighwaterMarkOffset >> 24), byte(t.HighwaterMarkOffset >> 16), byte(t.HighwaterMarkOffset >> 8), byte(t.HighwaterMarkOffset)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.MessageSetSize >> 24), byte(t.MessageSetSize >> 16), byte(t.MessageSetSize >> 8), byte(t.MessageSetSize)}); err != nil {
		return err
	}
	if err := t.MessageSet.Marshal(w); err != nil {
		return err
	}
	return nil
}

func (t *OffsetRequest) Marshal(w io.Writer) error {
	if _, err := w.Write([]byte{byte(t.ReplicaID >> 24), byte(t.ReplicaID >> 16), byte(t.ReplicaID >> 8), byte(t.ReplicaID)}); err != nil {
		return err
	}
	{
		l := int32(len(t.TimeInTopics))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t.TimeInTopics {
			if err := t.TimeInTopics[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *TimeInTopic) Marshal(w io.Writer) error {
	{
		l := int16(len(t.TopicName))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(t.TopicName)); err != nil {
			return err
		}
	}
	{
		l := int32(len(t.TimeInPartitions))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t.TimeInPartitions {
			if err := t.TimeInPartitions[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *TimeInPartition) Marshal(w io.Writer) error {
	if _, err := w.Write([]byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.Time >> 56), byte(t.Time >> 48), byte(t.Time >> 32), byte(t.Time >> 24), byte(t.Time >> 16), byte(t.Time >> 8), byte(t.Time)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.MaxNumberOfOffsets >> 24), byte(t.MaxNumberOfOffsets >> 16), byte(t.MaxNumberOfOffsets >> 8), byte(t.MaxNumberOfOffsets)}); err != nil {
		return err
	}
	return nil
}

func (t OffsetResponse) Marshal(w io.Writer) error {
	{
		l := int32(len(t))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t {
			if err := t[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *OffsetsInTopic) Marshal(w io.Writer) error {
	{
		l := int16(len(t.TopicName))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(t.TopicName)); err != nil {
			return err
		}
	}
	{
		l := int32(len(t.OffsetsInPartitions))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t.OffsetsInPartitions {
			if err := t.OffsetsInPartitions[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *OffsetsInPartition) Marshal(w io.Writer) error {
	if _, err := w.Write([]byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.ErrorCode >> 8), byte(t.ErrorCode)}); err != nil {
		return err
	}
	{
		l := int32(len(t.Offsets))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t.Offsets {
			if _, err := w.Write([]byte{byte(t.Offsets[i] >> 56), byte(t.Offsets[i] >> 48), byte(t.Offsets[i] >> 32), byte(t.Offsets[i] >> 24), byte(t.Offsets[i] >> 16), byte(t.Offsets[i] >> 8), byte(t.Offsets[i])}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t ConsumerMetadataRequest) Marshal(w io.Writer) error {
	{
		l := int16(len(t))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(t)); err != nil {
			return err
		}
	}
	return nil
}

func (t *ConsumerMetadataResponse) Marshal(w io.Writer) error {
	if _, err := w.Write([]byte{byte(t.ErrorCode >> 8), byte(t.ErrorCode)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.CoordinatorID >> 24), byte(t.CoordinatorID >> 16), byte(t.CoordinatorID >> 8), byte(t.CoordinatorID)}); err != nil {
		return err
	}
	{
		l := int16(len(t.CoordinatorHost))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(t.CoordinatorHost)); err != nil {
			return err
		}
	}
	if _, err := w.Write([]byte{byte(t.CoordinatorPort >> 24), byte(t.CoordinatorPort >> 16), byte(t.CoordinatorPort >> 8), byte(t.CoordinatorPort)}); err != nil {
		return err
	}
	return nil
}

func (t *OffsetCommitRequestV0) Marshal(w io.Writer) error {
	{
		l := int16(len(t.ConsumerGroupID))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(t.ConsumerGroupID)); err != nil {
			return err
		}
	}
	{
		l := int32(len(t.OffsetCommitInTopicV0s))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t.OffsetCommitInTopicV0s {
			if err := t.OffsetCommitInTopicV0s[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *OffsetCommitInTopicV0) Marshal(w io.Writer) error {
	{
		l := int16(len(t.TopicName))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(t.TopicName)); err != nil {
			return err
		}
	}
	{
		l := int32(len(t.OffsetCommitInPartitionV0s))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t.OffsetCommitInPartitionV0s {
			if err := t.OffsetCommitInPartitionV0s[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *OffsetCommitInPartitionV0) Marshal(w io.Writer) error {
	if _, err := w.Write([]byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.Offset >> 56), byte(t.Offset >> 48), byte(t.Offset >> 32), byte(t.Offset >> 24), byte(t.Offset >> 16), byte(t.Offset >> 8), byte(t.Offset)}); err != nil {
		return err
	}
	{
		l := int16(len(t.Metadata))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(t.Metadata)); err != nil {
			return err
		}
	}
	return nil
}

func (t *OffsetCommitRequestV1) Marshal(w io.Writer) error {
	{
		l := int16(len(t.ConsumerGroupID))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(t.ConsumerGroupID)); err != nil {
			return err
		}
	}
	if _, err := w.Write([]byte{byte(t.ConsumerGroupGenerationID >> 24), byte(t.ConsumerGroupGenerationID >> 16), byte(t.ConsumerGroupGenerationID >> 8), byte(t.ConsumerGroupGenerationID)}); err != nil {
		return err
	}
	{
		l := int16(len(t.ConsumerID))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(t.ConsumerID)); err != nil {
			return err
		}
	}
	{
		l := int32(len(t.OffsetCommitInTopicV1s))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t.OffsetCommitInTopicV1s {
			if err := t.OffsetCommitInTopicV1s[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *OffsetCommitInTopicV1) Marshal(w io.Writer) error {
	{
		l := int16(len(t.TopicName))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(t.TopicName)); err != nil {
			return err
		}
	}
	{
		l := int32(len(t.OffsetCommitInPartitionV1s))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t.OffsetCommitInPartitionV1s {
			if err := t.OffsetCommitInPartitionV1s[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *OffsetCommitInPartitionV1) Marshal(w io.Writer) error {
	if _, err := w.Write([]byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.Offset >> 56), byte(t.Offset >> 48), byte(t.Offset >> 32), byte(t.Offset >> 24), byte(t.Offset >> 16), byte(t.Offset >> 8), byte(t.Offset)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.TimeStamp >> 56), byte(t.TimeStamp >> 48), byte(t.TimeStamp >> 32), byte(t.TimeStamp >> 24), byte(t.TimeStamp >> 16), byte(t.TimeStamp >> 8), byte(t.TimeStamp)}); err != nil {
		return err
	}
	{
		l := int16(len(t.Metadata))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(t.Metadata)); err != nil {
			return err
		}
	}
	return nil
}

func (t *OffsetCommitRequestV2) Marshal(w io.Writer) error {
	{
		l := int16(len(t.ConsumerGroup))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(t.ConsumerGroup)); err != nil {
			return err
		}
	}
	if _, err := w.Write([]byte{byte(t.ConsumerGroupGenerationID >> 24), byte(t.ConsumerGroupGenerationID >> 16), byte(t.ConsumerGroupGenerationID >> 8), byte(t.ConsumerGroupGenerationID)}); err != nil {
		return err
	}
	{
		l := int16(len(t.ConsumerID))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(t.ConsumerID)); err != nil {
			return err
		}
	}
	if _, err := w.Write([]byte{byte(t.RetentionTime >> 56), byte(t.RetentionTime >> 48), byte(t.RetentionTime >> 32), byte(t.RetentionTime >> 24), byte(t.RetentionTime >> 16), byte(t.RetentionTime >> 8), byte(t.RetentionTime)}); err != nil {
		return err
	}
	{
		l := int32(len(t.OffsetCommitInTopicV2s))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t.OffsetCommitInTopicV2s {
			if err := t.OffsetCommitInTopicV2s[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *OffsetCommitInTopicV2) Marshal(w io.Writer) error {
	{
		l := int16(len(t.TopicName))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(t.TopicName)); err != nil {
			return err
		}
	}
	{
		l := int32(len(t.OffsetCommitInPartitionV2s))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t.OffsetCommitInPartitionV2s {
			if err := t.OffsetCommitInPartitionV2s[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *OffsetCommitInPartitionV2) Marshal(w io.Writer) error {
	if _, err := w.Write([]byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.Offset >> 56), byte(t.Offset >> 48), byte(t.Offset >> 32), byte(t.Offset >> 24), byte(t.Offset >> 16), byte(t.Offset >> 8), byte(t.Offset)}); err != nil {
		return err
	}
	{
		l := int16(len(t.Metadata))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(t.Metadata)); err != nil {
			return err
		}
	}
	return nil
}

func (t OffsetCommitResponse) Marshal(w io.Writer) error {
	{
		l := int32(len(t))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t {
			if err := t[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *ErrorInTopic) Marshal(w io.Writer) error {
	{
		l := int16(len(t.TopicName))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(t.TopicName)); err != nil {
			return err
		}
	}
	{
		l := int32(len(t.ErrorInPartitions))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t.ErrorInPartitions {
			if err := t.ErrorInPartitions[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *ErrorInPartition) Marshal(w io.Writer) error {
	if _, err := w.Write([]byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.ErrorCode >> 8), byte(t.ErrorCode)}); err != nil {
		return err
	}
	return nil
}

func (t *OffsetFetchRequest) Marshal(w io.Writer) error {
	{
		l := int16(len(t.ConsumerGroup))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(t.ConsumerGroup)); err != nil {
			return err
		}
	}
	{
		l := int32(len(t.PartitionInTopics))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t.PartitionInTopics {
			if err := t.PartitionInTopics[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *PartitionInTopic) Marshal(w io.Writer) error {
	{
		l := int16(len(t.TopicName))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(t.TopicName)); err != nil {
			return err
		}
	}
	{
		l := int32(len(t.Partitions))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t.Partitions {
			if _, err := w.Write([]byte{byte(t.Partitions[i] >> 24), byte(t.Partitions[i] >> 16), byte(t.Partitions[i] >> 8), byte(t.Partitions[i])}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t OffsetFetchResponse) Marshal(w io.Writer) error {
	{
		l := int32(len(t))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t {
			if err := t[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *OffsetMetadataInTopic) Marshal(w io.Writer) error {
	{
		l := int16(len(t.TopicName))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(t.TopicName)); err != nil {
			return err
		}
	}
	{
		l := int32(len(t.OffsetMetadataInPartitions))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range t.OffsetMetadataInPartitions {
			if err := t.OffsetMetadataInPartitions[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *OffsetMetadataInPartition) Marshal(w io.Writer) error {
	if _, err := w.Write([]byte{byte(t.Partition >> 24), byte(t.Partition >> 16), byte(t.Partition >> 8), byte(t.Partition)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.Offset >> 56), byte(t.Offset >> 48), byte(t.Offset >> 32), byte(t.Offset >> 24), byte(t.Offset >> 16), byte(t.Offset >> 8), byte(t.Offset)}); err != nil {
		return err
	}
	{
		l := int16(len(t.Metadata))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(t.Metadata)); err != nil {
			return err
		}
	}
	if _, err := w.Write([]byte{byte(t.ErrorCode >> 8), byte(t.ErrorCode)}); err != nil {
		return err
	}
	return nil
}
