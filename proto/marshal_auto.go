package proto

import (
	"hash/crc32"
)

func (t *RequestOrResponse) Marshal(w *Writer) {
	offset := len(w.B)
	w.WriteInt32(t.Size)
	start := len(w.B)
	t.M.Marshal(w)
	w.SetInt32(offset, int32(len(w.B)-start))
}

func (t *RequestOrResponse) Unmarshal(r *Reader) {
	t.Size = r.ReadInt32()
	start := r.Offset
	t.M.Unmarshal(r)
	if r.Err == nil && int(t.Size) != r.Offset-start {
		r.Err = ErrSizeMismatch
	}
}

func (t *Request) Marshal(w *Writer) {
	w.WriteInt16(t.APIKey)
	w.WriteInt16(t.APIVersion)
	w.WriteInt32(t.CorrelationID)
	w.WriteString(t.ClientID)
	t.RequestMessage.Marshal(w)
}

func (t *Request) Unmarshal(r *Reader) {
	t.APIKey = r.ReadInt16()
	t.APIVersion = r.ReadInt16()
	t.CorrelationID = r.ReadInt32()
	t.ClientID = r.ReadString()
	t.RequestMessage.Unmarshal(r)
}

func (t *Response) Marshal(w *Writer) {
	w.WriteInt32(t.CorrelationID)
	t.ResponseMessage.Marshal(w)
}

func (t *Response) Unmarshal(r *Reader) {
	t.CorrelationID = r.ReadInt32()
	t.ResponseMessage.Unmarshal(r)
}

func (t *MessageSet) Marshal(w *Writer) {
	offset := len(w.B)
	w.WriteInt32(0)
	start := len(w.B)
	for i := range *t {
		(*t)[i].Marshal(w)
	}
	w.SetInt32(offset, int32(len(w.B)-start))
}

func (t *MessageSet) Unmarshal(r *Reader) {
	size := int(r.ReadInt32())
	start := r.Offset
	for r.Offset-start < size {
		var m OffsetMessage
		m.Unmarshal(r)
		if r.Err != nil {
			r.Err = nil
			r.Offset = len(r.B)
			return
		}
		*t = append(*t, m)
	}
}

func (t *OffsetMessage) Marshal(w *Writer) {
	w.WriteInt64(t.Offset)
	t.SizedMessage.Marshal(w)
}

func (t *OffsetMessage) Unmarshal(r *Reader) {
	t.Offset = r.ReadInt64()
	t.SizedMessage.Unmarshal(r)
}

func (t *SizedMessage) Marshal(w *Writer) {
	offset := len(w.B)
	w.WriteInt32(t.Size)
	start := len(w.B)
	t.CRCMessage.Marshal(w)
	w.SetInt32(offset, int32(len(w.B)-start))
}

func (t *SizedMessage) Unmarshal(r *Reader) {
	t.Size = r.ReadInt32()
	start := r.Offset
	t.CRCMessage.Unmarshal(r)
	if r.Err == nil && int(t.Size) != r.Offset-start {
		r.Err = ErrSizeMismatch
	}
}

func (t *CRCMessage) Marshal(w *Writer) {
	offset := len(w.B)
	w.WriteUint32(t.CRC)
	start := len(w.B)
	t.Message.Marshal(w)
	w.SetUint32(offset, crc32.ChecksumIEEE(w.B[start:]))
}

func (t *CRCMessage) Unmarshal(r *Reader) {
	t.CRC = r.ReadUint32()
	start := r.Offset
	t.Message.Unmarshal(r)
	if r.Err == nil && t.CRC != crc32.ChecksumIEEE(r.B[start:r.Offset]) {
		r.Err = ErrCRCMismatch
	}
}

func (t *Message) Marshal(w *Writer) {
	w.WriteInt8(t.MagicByte)
	w.WriteInt8(t.Attributes)
	w.WriteBytes(t.Key)
	w.WriteBytes(t.Value)
}

func (t *Message) Unmarshal(r *Reader) {
	t.MagicByte = r.ReadInt8()
	t.Attributes = r.ReadInt8()
	t.Key = r.ReadBytes()
	t.Value = r.ReadBytes()
}

func (t *TopicMetadataRequest) Marshal(w *Writer) {
	w.WriteInt32(int32(len((*t))))
	for i := range *t {
		w.WriteString((*t)[i])
	}
}

func (t *TopicMetadataRequest) Unmarshal(r *Reader) {
	(*t) = make([]string, int(r.ReadInt32()))
	for i := range *t {
		(*t)[i] = r.ReadString()
	}
}

func (t *TopicMetadataResponse) Marshal(w *Writer) {
	w.WriteInt32(int32(len(t.Brokers)))
	for i := range t.Brokers {
		t.Brokers[i].Marshal(w)
	}
	w.WriteInt32(int32(len(t.TopicMetadatas)))
	for i := range t.TopicMetadatas {
		t.TopicMetadatas[i].Marshal(w)
	}
}

func (t *TopicMetadataResponse) Unmarshal(r *Reader) {
	t.Brokers = make([]Broker, int(r.ReadInt32()))
	for i := range t.Brokers {
		t.Brokers[i].Unmarshal(r)
	}
	t.TopicMetadatas = make([]TopicMetadata, int(r.ReadInt32()))
	for i := range t.TopicMetadatas {
		t.TopicMetadatas[i].Unmarshal(r)
	}
}

func (t *Broker) Marshal(w *Writer) {
	w.WriteInt32(t.NodeID)
	w.WriteString(t.Host)
	w.WriteInt32(t.Port)
}

func (t *Broker) Unmarshal(r *Reader) {
	t.NodeID = r.ReadInt32()
	t.Host = r.ReadString()
	t.Port = r.ReadInt32()
}

func (t *TopicMetadata) Marshal(w *Writer) {
	t.ErrorCode.Marshal(w)
	w.WriteString(t.TopicName)
	w.WriteInt32(int32(len(t.PartitionMetadatas)))
	for i := range t.PartitionMetadatas {
		t.PartitionMetadatas[i].Marshal(w)
	}
}

func (t *TopicMetadata) Unmarshal(r *Reader) {
	t.ErrorCode.Unmarshal(r)
	t.TopicName = r.ReadString()
	t.PartitionMetadatas = make([]PartitionMetadata, int(r.ReadInt32()))
	for i := range t.PartitionMetadatas {
		t.PartitionMetadatas[i].Unmarshal(r)
	}
}

func (t *PartitionMetadata) Marshal(w *Writer) {
	t.ErrorCode.Marshal(w)
	w.WriteInt32(t.PartitionID)
	w.WriteInt32(t.Leader)
	w.WriteInt32(int32(len(t.Replicas)))
	for i := range t.Replicas {
		w.WriteInt32(t.Replicas[i])
	}
	w.WriteInt32(int32(len(t.ISR)))
	for i := range t.ISR {
		w.WriteInt32(t.ISR[i])
	}
}

func (t *PartitionMetadata) Unmarshal(r *Reader) {
	t.ErrorCode.Unmarshal(r)
	t.PartitionID = r.ReadInt32()
	t.Leader = r.ReadInt32()
	t.Replicas = make([]int32, int(r.ReadInt32()))
	for i := range t.Replicas {
		t.Replicas[i] = r.ReadInt32()
	}
	t.ISR = make([]int32, int(r.ReadInt32()))
	for i := range t.ISR {
		t.ISR[i] = r.ReadInt32()
	}
}

func (t *ProduceRequest) Marshal(w *Writer) {
	w.WriteInt16(t.RequiredAcks)
	w.WriteInt32(t.Timeout)
	w.WriteInt32(int32(len(t.MessageSetInTopics)))
	for i := range t.MessageSetInTopics {
		t.MessageSetInTopics[i].Marshal(w)
	}
}

func (t *ProduceRequest) Unmarshal(r *Reader) {
	t.RequiredAcks = r.ReadInt16()
	t.Timeout = r.ReadInt32()
	t.MessageSetInTopics = make([]MessageSetInTopic, int(r.ReadInt32()))
	for i := range t.MessageSetInTopics {
		t.MessageSetInTopics[i].Unmarshal(r)
	}
}

func (t *MessageSetInTopic) Marshal(w *Writer) {
	w.WriteString(t.TopicName)
	w.WriteInt32(int32(len(t.MessageSetInPartitions)))
	for i := range t.MessageSetInPartitions {
		t.MessageSetInPartitions[i].Marshal(w)
	}
}

func (t *MessageSetInTopic) Unmarshal(r *Reader) {
	t.TopicName = r.ReadString()
	t.MessageSetInPartitions = make([]MessageSetInPartition, int(r.ReadInt32()))
	for i := range t.MessageSetInPartitions {
		t.MessageSetInPartitions[i].Unmarshal(r)
	}
}

func (t *MessageSetInPartition) Marshal(w *Writer) {
	w.WriteInt32(t.Partition)
	t.MessageSet.Marshal(w)
}

func (t *MessageSetInPartition) Unmarshal(r *Reader) {
	t.Partition = r.ReadInt32()
	t.MessageSet.Unmarshal(r)
}

func (t *ProduceResponse) Marshal(w *Writer) {
	w.WriteInt32(int32(len((*t))))
	for i := range *t {
		(*t)[i].Marshal(w)
	}
}

func (t *ProduceResponse) Unmarshal(r *Reader) {
	(*t) = make([]OffsetInTopic, int(r.ReadInt32()))
	for i := range *t {
		(*t)[i].Unmarshal(r)
	}
}

func (t *OffsetInTopic) Marshal(w *Writer) {
	w.WriteString(t.TopicName)
	w.WriteInt32(int32(len(t.OffsetInPartitions)))
	for i := range t.OffsetInPartitions {
		t.OffsetInPartitions[i].Marshal(w)
	}
}

func (t *OffsetInTopic) Unmarshal(r *Reader) {
	t.TopicName = r.ReadString()
	t.OffsetInPartitions = make([]OffsetInPartition, int(r.ReadInt32()))
	for i := range t.OffsetInPartitions {
		t.OffsetInPartitions[i].Unmarshal(r)
	}
}

func (t *OffsetInPartition) Marshal(w *Writer) {
	w.WriteInt32(t.Partition)
	t.ErrorCode.Marshal(w)
	w.WriteInt64(t.Offset)
}

func (t *OffsetInPartition) Unmarshal(r *Reader) {
	t.Partition = r.ReadInt32()
	t.ErrorCode.Unmarshal(r)
	t.Offset = r.ReadInt64()
}

func (t *FetchRequest) Marshal(w *Writer) {
	w.WriteInt32(t.ReplicaID)
	w.WriteInt32(t.MaxWaitTime)
	w.WriteInt32(t.MinBytes)
	w.WriteInt32(int32(len(t.FetchOffsetInTopics)))
	for i := range t.FetchOffsetInTopics {
		t.FetchOffsetInTopics[i].Marshal(w)
	}
}

func (t *FetchRequest) Unmarshal(r *Reader) {
	t.ReplicaID = r.ReadInt32()
	t.MaxWaitTime = r.ReadInt32()
	t.MinBytes = r.ReadInt32()
	t.FetchOffsetInTopics = make([]FetchOffsetInTopic, int(r.ReadInt32()))
	for i := range t.FetchOffsetInTopics {
		t.FetchOffsetInTopics[i].Unmarshal(r)
	}
}

func (t *FetchOffsetInTopic) Marshal(w *Writer) {
	w.WriteString(t.TopicName)
	w.WriteInt32(int32(len(t.FetchOffsetInPartitions)))
	for i := range t.FetchOffsetInPartitions {
		t.FetchOffsetInPartitions[i].Marshal(w)
	}
}

func (t *FetchOffsetInTopic) Unmarshal(r *Reader) {
	t.TopicName = r.ReadString()
	t.FetchOffsetInPartitions = make([]FetchOffsetInPartition, int(r.ReadInt32()))
	for i := range t.FetchOffsetInPartitions {
		t.FetchOffsetInPartitions[i].Unmarshal(r)
	}
}

func (t *FetchOffsetInPartition) Marshal(w *Writer) {
	w.WriteInt32(t.Partition)
	w.WriteInt64(t.FetchOffset)
	w.WriteInt32(t.MaxBytes)
}

func (t *FetchOffsetInPartition) Unmarshal(r *Reader) {
	t.Partition = r.ReadInt32()
	t.FetchOffset = r.ReadInt64()
	t.MaxBytes = r.ReadInt32()
}

func (t *FetchResponse) Marshal(w *Writer) {
	w.WriteInt32(int32(len((*t))))
	for i := range *t {
		(*t)[i].Marshal(w)
	}
}

func (t *FetchResponse) Unmarshal(r *Reader) {
	(*t) = make([]FetchMessageSetInTopic, int(r.ReadInt32()))
	for i := range *t {
		(*t)[i].Unmarshal(r)
	}
}

func (t *FetchMessageSetInTopic) Marshal(w *Writer) {
	w.WriteString(t.TopicName)
	w.WriteInt32(int32(len(t.FetchMessageSetInPartitions)))
	for i := range t.FetchMessageSetInPartitions {
		t.FetchMessageSetInPartitions[i].Marshal(w)
	}
}

func (t *FetchMessageSetInTopic) Unmarshal(r *Reader) {
	t.TopicName = r.ReadString()
	t.FetchMessageSetInPartitions = make([]FetchMessageSetInPartition, int(r.ReadInt32()))
	for i := range t.FetchMessageSetInPartitions {
		t.FetchMessageSetInPartitions[i].Unmarshal(r)
	}
}

func (t *FetchMessageSetInPartition) Marshal(w *Writer) {
	w.WriteInt32(t.Partition)
	t.ErrorCode.Marshal(w)
	w.WriteInt64(t.HighwaterMarkOffset)
	t.MessageSet.Marshal(w)
}

func (t *FetchMessageSetInPartition) Unmarshal(r *Reader) {
	t.Partition = r.ReadInt32()
	t.ErrorCode.Unmarshal(r)
	t.HighwaterMarkOffset = r.ReadInt64()
	t.MessageSet.Unmarshal(r)
}

func (t *OffsetRequest) Marshal(w *Writer) {
	w.WriteInt32(t.ReplicaID)
	w.WriteInt32(int32(len(t.TimeInTopics)))
	for i := range t.TimeInTopics {
		t.TimeInTopics[i].Marshal(w)
	}
}

func (t *OffsetRequest) Unmarshal(r *Reader) {
	t.ReplicaID = r.ReadInt32()
	t.TimeInTopics = make([]TimeInTopic, int(r.ReadInt32()))
	for i := range t.TimeInTopics {
		t.TimeInTopics[i].Unmarshal(r)
	}
}

func (t *TimeInTopic) Marshal(w *Writer) {
	w.WriteString(t.TopicName)
	w.WriteInt32(int32(len(t.TimeInPartitions)))
	for i := range t.TimeInPartitions {
		t.TimeInPartitions[i].Marshal(w)
	}
}

func (t *TimeInTopic) Unmarshal(r *Reader) {
	t.TopicName = r.ReadString()
	t.TimeInPartitions = make([]TimeInPartition, int(r.ReadInt32()))
	for i := range t.TimeInPartitions {
		t.TimeInPartitions[i].Unmarshal(r)
	}
}

func (t *TimeInPartition) Marshal(w *Writer) {
	w.WriteInt32(t.Partition)
	w.WriteInt64(t.Time)
	w.WriteInt32(t.MaxNumberOfOffsets)
}

func (t *TimeInPartition) Unmarshal(r *Reader) {
	t.Partition = r.ReadInt32()
	t.Time = r.ReadInt64()
	t.MaxNumberOfOffsets = r.ReadInt32()
}

func (t *OffsetResponse) Marshal(w *Writer) {
	w.WriteInt32(int32(len((*t))))
	for i := range *t {
		(*t)[i].Marshal(w)
	}
}

func (t *OffsetResponse) Unmarshal(r *Reader) {
	(*t) = make([]OffsetsInTopic, int(r.ReadInt32()))
	for i := range *t {
		(*t)[i].Unmarshal(r)
	}
}

func (t *OffsetsInTopic) Marshal(w *Writer) {
	w.WriteString(t.TopicName)
	w.WriteInt32(int32(len(t.OffsetsInPartitions)))
	for i := range t.OffsetsInPartitions {
		t.OffsetsInPartitions[i].Marshal(w)
	}
}

func (t *OffsetsInTopic) Unmarshal(r *Reader) {
	t.TopicName = r.ReadString()
	t.OffsetsInPartitions = make([]OffsetsInPartition, int(r.ReadInt32()))
	for i := range t.OffsetsInPartitions {
		t.OffsetsInPartitions[i].Unmarshal(r)
	}
}

func (t *OffsetsInPartition) Marshal(w *Writer) {
	w.WriteInt32(t.Partition)
	t.ErrorCode.Marshal(w)
	w.WriteInt32(int32(len(t.Offsets)))
	for i := range t.Offsets {
		w.WriteInt64(t.Offsets[i])
	}
}

func (t *OffsetsInPartition) Unmarshal(r *Reader) {
	t.Partition = r.ReadInt32()
	t.ErrorCode.Unmarshal(r)
	t.Offsets = make([]int64, int(r.ReadInt32()))
	for i := range t.Offsets {
		t.Offsets[i] = r.ReadInt64()
	}
}

func (t *GroupCoordinatorRequest) Marshal(w *Writer) {
	w.WriteString(string((*t)))
}

func (t *GroupCoordinatorRequest) Unmarshal(r *Reader) {
	(*t) = GroupCoordinatorRequest(r.ReadString())
}

func (t *GroupCoordinatorResponse) Marshal(w *Writer) {
	t.ErrorCode.Marshal(w)
	t.Broker.Marshal(w)
}

func (t *GroupCoordinatorResponse) Unmarshal(r *Reader) {
	t.ErrorCode.Unmarshal(r)
	t.Broker.Unmarshal(r)
}

func (t *OffsetCommitRequestV0) Marshal(w *Writer) {
	w.WriteString(t.ConsumerGroupID)
	w.WriteInt32(int32(len(t.OffsetCommitInTopicV0s)))
	for i := range t.OffsetCommitInTopicV0s {
		t.OffsetCommitInTopicV0s[i].Marshal(w)
	}
}

func (t *OffsetCommitRequestV0) Unmarshal(r *Reader) {
	t.ConsumerGroupID = r.ReadString()
	t.OffsetCommitInTopicV0s = make([]OffsetCommitInTopicV0, int(r.ReadInt32()))
	for i := range t.OffsetCommitInTopicV0s {
		t.OffsetCommitInTopicV0s[i].Unmarshal(r)
	}
}

func (t *OffsetCommitInTopicV0) Marshal(w *Writer) {
	w.WriteString(t.TopicName)
	w.WriteInt32(int32(len(t.OffsetCommitInPartitionV0s)))
	for i := range t.OffsetCommitInPartitionV0s {
		t.OffsetCommitInPartitionV0s[i].Marshal(w)
	}
}

func (t *OffsetCommitInTopicV0) Unmarshal(r *Reader) {
	t.TopicName = r.ReadString()
	t.OffsetCommitInPartitionV0s = make([]OffsetCommitInPartitionV0, int(r.ReadInt32()))
	for i := range t.OffsetCommitInPartitionV0s {
		t.OffsetCommitInPartitionV0s[i].Unmarshal(r)
	}
}

func (t *OffsetCommitInPartitionV0) Marshal(w *Writer) {
	w.WriteInt32(t.Partition)
	w.WriteInt64(t.Offset)
	w.WriteString(t.Metadata)
}

func (t *OffsetCommitInPartitionV0) Unmarshal(r *Reader) {
	t.Partition = r.ReadInt32()
	t.Offset = r.ReadInt64()
	t.Metadata = r.ReadString()
}

func (t *OffsetCommitRequestV1) Marshal(w *Writer) {
	w.WriteString(t.ConsumerGroupID)
	w.WriteInt32(t.ConsumerGroupGenerationID)
	w.WriteString(t.ConsumerID)
	w.WriteInt32(int32(len(t.OffsetCommitInTopicV1s)))
	for i := range t.OffsetCommitInTopicV1s {
		t.OffsetCommitInTopicV1s[i].Marshal(w)
	}
}

func (t *OffsetCommitRequestV1) Unmarshal(r *Reader) {
	t.ConsumerGroupID = r.ReadString()
	t.ConsumerGroupGenerationID = r.ReadInt32()
	t.ConsumerID = r.ReadString()
	t.OffsetCommitInTopicV1s = make([]OffsetCommitInTopicV1, int(r.ReadInt32()))
	for i := range t.OffsetCommitInTopicV1s {
		t.OffsetCommitInTopicV1s[i].Unmarshal(r)
	}
}

func (t *OffsetCommitInTopicV1) Marshal(w *Writer) {
	w.WriteString(t.TopicName)
	w.WriteInt32(int32(len(t.OffsetCommitInPartitionV1s)))
	for i := range t.OffsetCommitInPartitionV1s {
		t.OffsetCommitInPartitionV1s[i].Marshal(w)
	}
}

func (t *OffsetCommitInTopicV1) Unmarshal(r *Reader) {
	t.TopicName = r.ReadString()
	t.OffsetCommitInPartitionV1s = make([]OffsetCommitInPartitionV1, int(r.ReadInt32()))
	for i := range t.OffsetCommitInPartitionV1s {
		t.OffsetCommitInPartitionV1s[i].Unmarshal(r)
	}
}

func (t *OffsetCommitInPartitionV1) Marshal(w *Writer) {
	w.WriteInt32(t.Partition)
	w.WriteInt64(t.Offset)
	w.WriteInt64(t.TimeStamp)
	w.WriteString(t.Metadata)
}

func (t *OffsetCommitInPartitionV1) Unmarshal(r *Reader) {
	t.Partition = r.ReadInt32()
	t.Offset = r.ReadInt64()
	t.TimeStamp = r.ReadInt64()
	t.Metadata = r.ReadString()
}

func (t *OffsetCommitRequestV2) Marshal(w *Writer) {
	w.WriteString(t.ConsumerGroup)
	w.WriteInt32(t.ConsumerGroupGenerationID)
	w.WriteString(t.ConsumerID)
	w.WriteInt64(t.RetentionTime)
	w.WriteInt32(int32(len(t.OffsetCommitInTopicV2s)))
	for i := range t.OffsetCommitInTopicV2s {
		t.OffsetCommitInTopicV2s[i].Marshal(w)
	}
}

func (t *OffsetCommitRequestV2) Unmarshal(r *Reader) {
	t.ConsumerGroup = r.ReadString()
	t.ConsumerGroupGenerationID = r.ReadInt32()
	t.ConsumerID = r.ReadString()
	t.RetentionTime = r.ReadInt64()
	t.OffsetCommitInTopicV2s = make([]OffsetCommitInTopicV2, int(r.ReadInt32()))
	for i := range t.OffsetCommitInTopicV2s {
		t.OffsetCommitInTopicV2s[i].Unmarshal(r)
	}
}

func (t *OffsetCommitInTopicV2) Marshal(w *Writer) {
	w.WriteString(t.TopicName)
	w.WriteInt32(int32(len(t.OffsetCommitInPartitionV2s)))
	for i := range t.OffsetCommitInPartitionV2s {
		t.OffsetCommitInPartitionV2s[i].Marshal(w)
	}
}

func (t *OffsetCommitInTopicV2) Unmarshal(r *Reader) {
	t.TopicName = r.ReadString()
	t.OffsetCommitInPartitionV2s = make([]OffsetCommitInPartitionV2, int(r.ReadInt32()))
	for i := range t.OffsetCommitInPartitionV2s {
		t.OffsetCommitInPartitionV2s[i].Unmarshal(r)
	}
}

func (t *OffsetCommitInPartitionV2) Marshal(w *Writer) {
	w.WriteInt32(t.Partition)
	w.WriteInt64(t.Offset)
	w.WriteString(t.Metadata)
}

func (t *OffsetCommitInPartitionV2) Unmarshal(r *Reader) {
	t.Partition = r.ReadInt32()
	t.Offset = r.ReadInt64()
	t.Metadata = r.ReadString()
}

func (t *OffsetCommitResponse) Marshal(w *Writer) {
	w.WriteInt32(int32(len((*t))))
	for i := range *t {
		(*t)[i].Marshal(w)
	}
}

func (t *OffsetCommitResponse) Unmarshal(r *Reader) {
	(*t) = make([]ErrorInTopic, int(r.ReadInt32()))
	for i := range *t {
		(*t)[i].Unmarshal(r)
	}
}

func (t *ErrorInTopic) Marshal(w *Writer) {
	w.WriteString(t.TopicName)
	w.WriteInt32(int32(len(t.ErrorInPartitions)))
	for i := range t.ErrorInPartitions {
		t.ErrorInPartitions[i].Marshal(w)
	}
}

func (t *ErrorInTopic) Unmarshal(r *Reader) {
	t.TopicName = r.ReadString()
	t.ErrorInPartitions = make([]ErrorInPartition, int(r.ReadInt32()))
	for i := range t.ErrorInPartitions {
		t.ErrorInPartitions[i].Unmarshal(r)
	}
}

func (t *ErrorInPartition) Marshal(w *Writer) {
	w.WriteInt32(t.Partition)
	t.ErrorCode.Marshal(w)
}

func (t *ErrorInPartition) Unmarshal(r *Reader) {
	t.Partition = r.ReadInt32()
	t.ErrorCode.Unmarshal(r)
}

func (t *OffsetFetchRequestV0) Marshal(w *Writer) {
	w.WriteString(t.ConsumerGroup)
	w.WriteInt32(int32(len(t.PartitionInTopics)))
	for i := range t.PartitionInTopics {
		t.PartitionInTopics[i].Marshal(w)
	}
}

func (t *OffsetFetchRequestV0) Unmarshal(r *Reader) {
	t.ConsumerGroup = r.ReadString()
	t.PartitionInTopics = make([]PartitionInTopic, int(r.ReadInt32()))
	for i := range t.PartitionInTopics {
		t.PartitionInTopics[i].Unmarshal(r)
	}
}

func (t *PartitionInTopic) Marshal(w *Writer) {
	w.WriteString(t.TopicName)
	w.WriteInt32(int32(len(t.Partitions)))
	for i := range t.Partitions {
		w.WriteInt32(t.Partitions[i])
	}
}

func (t *PartitionInTopic) Unmarshal(r *Reader) {
	t.TopicName = r.ReadString()
	t.Partitions = make([]int32, int(r.ReadInt32()))
	for i := range t.Partitions {
		t.Partitions[i] = r.ReadInt32()
	}
}

func (t *OffsetFetchRequestV1) Marshal(w *Writer) {
	w.WriteString(t.ConsumerGroup)
	w.WriteInt32(int32(len(t.PartitionInTopics)))
	for i := range t.PartitionInTopics {
		t.PartitionInTopics[i].Marshal(w)
	}
}

func (t *OffsetFetchRequestV1) Unmarshal(r *Reader) {
	t.ConsumerGroup = r.ReadString()
	t.PartitionInTopics = make([]PartitionInTopic, int(r.ReadInt32()))
	for i := range t.PartitionInTopics {
		t.PartitionInTopics[i].Unmarshal(r)
	}
}

func (t *OffsetFetchResponse) Marshal(w *Writer) {
	w.WriteInt32(int32(len((*t))))
	for i := range *t {
		(*t)[i].Marshal(w)
	}
}

func (t *OffsetFetchResponse) Unmarshal(r *Reader) {
	(*t) = make([]OffsetMetadataInTopic, int(r.ReadInt32()))
	for i := range *t {
		(*t)[i].Unmarshal(r)
	}
}

func (t *OffsetMetadataInTopic) Marshal(w *Writer) {
	w.WriteString(t.TopicName)
	w.WriteInt32(int32(len(t.OffsetMetadataInPartitions)))
	for i := range t.OffsetMetadataInPartitions {
		t.OffsetMetadataInPartitions[i].Marshal(w)
	}
}

func (t *OffsetMetadataInTopic) Unmarshal(r *Reader) {
	t.TopicName = r.ReadString()
	t.OffsetMetadataInPartitions = make([]OffsetMetadataInPartition, int(r.ReadInt32()))
	for i := range t.OffsetMetadataInPartitions {
		t.OffsetMetadataInPartitions[i].Unmarshal(r)
	}
}

func (t *OffsetMetadataInPartition) Marshal(w *Writer) {
	w.WriteInt32(t.Partition)
	w.WriteInt64(t.Offset)
	w.WriteString(t.Metadata)
	t.ErrorCode.Marshal(w)
}

func (t *OffsetMetadataInPartition) Unmarshal(r *Reader) {
	t.Partition = r.ReadInt32()
	t.Offset = r.ReadInt64()
	t.Metadata = r.ReadString()
	t.ErrorCode.Unmarshal(r)
}

func (t *JoinGroupRequest) Marshal(w *Writer) {
	w.WriteString(t.GroupID)
	w.WriteInt32(t.SessionTimeout)
	w.WriteString(t.MemberID)
	w.WriteString(t.ProtocolType)
	t.GroupProtocols.Marshal(w)
}

func (t *JoinGroupRequest) Unmarshal(r *Reader) {
	t.GroupID = r.ReadString()
	t.SessionTimeout = r.ReadInt32()
	t.MemberID = r.ReadString()
	t.ProtocolType = r.ReadString()
	t.GroupProtocols.Unmarshal(r)
}

func (t *GroupProtocols) Marshal(w *Writer) {
	w.WriteInt32(int32(len((*t))))
	for i := range *t {
		(*t)[i].Marshal(w)
	}
}

func (t *GroupProtocols) Unmarshal(r *Reader) {
	(*t) = make([]GroupProtocol, int(r.ReadInt32()))
	for i := range *t {
		(*t)[i].Unmarshal(r)
	}
}

func (t *GroupProtocol) Marshal(w *Writer) {
	w.WriteString(t.ProtocolName)
	t.ProtocolMetadata.Marshal(w)
}

func (t *GroupProtocol) Unmarshal(r *Reader) {
	t.ProtocolName = r.ReadString()
	t.ProtocolMetadata.Unmarshal(r)
}

func (t *ProtocolMetadata) Marshal(w *Writer) {
	w.WriteInt16(t.Version)
	t.Subscription.Marshal(w)
	w.WriteBytes(t.UserData)
}

func (t *ProtocolMetadata) Unmarshal(r *Reader) {
	t.Version = r.ReadInt16()
	t.Subscription.Unmarshal(r)
	t.UserData = r.ReadBytes()
}

func (t *Subscription) Marshal(w *Writer) {
	w.WriteInt32(int32(len((*t))))
	for i := range *t {
		w.WriteString((*t)[i])
	}
}

func (t *Subscription) Unmarshal(r *Reader) {
	(*t) = make([]string, int(r.ReadInt32()))
	for i := range *t {
		(*t)[i] = r.ReadString()
	}
}

func (t *JoinGroupResponse) Marshal(w *Writer) {
	t.ErrorCode.Marshal(w)
	w.WriteInt32(t.GenerationID)
	w.WriteString(t.GroupProtocolName)
	w.WriteString(t.LeaderID)
	w.WriteString(t.MemberID)
	t.MemberWithMetas.Marshal(w)
}

func (t *JoinGroupResponse) Unmarshal(r *Reader) {
	t.ErrorCode.Unmarshal(r)
	t.GenerationID = r.ReadInt32()
	t.GroupProtocolName = r.ReadString()
	t.LeaderID = r.ReadString()
	t.MemberID = r.ReadString()
	t.MemberWithMetas.Unmarshal(r)
}

func (t *MemberWithMetas) Marshal(w *Writer) {
	w.WriteInt32(int32(len((*t))))
	for i := range *t {
		(*t)[i].Marshal(w)
	}
}

func (t *MemberWithMetas) Unmarshal(r *Reader) {
	(*t) = make([]MemberWithMeta, int(r.ReadInt32()))
	for i := range *t {
		(*t)[i].Unmarshal(r)
	}
}

func (t *MemberWithMeta) Marshal(w *Writer) {
	w.WriteString(t.MemberID)
	w.WriteBytes(t.MemberMetadata)
}

func (t *MemberWithMeta) Unmarshal(r *Reader) {
	t.MemberID = r.ReadString()
	t.MemberMetadata = r.ReadBytes()
}

func (t *SyncGroupRequest) Marshal(w *Writer) {
	w.WriteString(t.GroupID)
	w.WriteInt32(t.GenerationID)
	w.WriteString(t.MemberID)
	t.GroupAssignments.Marshal(w)
}

func (t *SyncGroupRequest) Unmarshal(r *Reader) {
	t.GroupID = r.ReadString()
	t.GenerationID = r.ReadInt32()
	t.MemberID = r.ReadString()
	t.GroupAssignments.Unmarshal(r)
}

func (t *GroupAssignments) Marshal(w *Writer) {
	w.WriteInt32(int32(len((*t))))
	for i := range *t {
		(*t)[i].Marshal(w)
	}
}

func (t *GroupAssignments) Unmarshal(r *Reader) {
	(*t) = make([]GroupAssignment, int(r.ReadInt32()))
	for i := range *t {
		(*t)[i].Unmarshal(r)
	}
}

func (t *GroupAssignment) Marshal(w *Writer) {
	w.WriteString(t.MemberID)
	t.MemberAssignment.Marshal(w)
}

func (t *GroupAssignment) Unmarshal(r *Reader) {
	t.MemberID = r.ReadString()
	t.MemberAssignment.Unmarshal(r)
}

func (t *MemberAssignment) Marshal(w *Writer) {
	w.WriteInt16(t.Version)
	t.PartitionAssignments.Marshal(w)
}

func (t *MemberAssignment) Unmarshal(r *Reader) {
	t.Version = r.ReadInt16()
	t.PartitionAssignments.Unmarshal(r)
}

func (t *PartitionAssignments) Marshal(w *Writer) {
	w.WriteInt32(int32(len((*t))))
	for i := range *t {
		(*t)[i].Marshal(w)
	}
}

func (t *PartitionAssignments) Unmarshal(r *Reader) {
	(*t) = make([]PartitionAssignment, int(r.ReadInt32()))
	for i := range *t {
		(*t)[i].Unmarshal(r)
	}
}

func (t *PartitionAssignment) Marshal(w *Writer) {
	w.WriteString(t.Topic)
	w.WriteInt32(int32(len(t.Partitions)))
	for i := range t.Partitions {
		w.WriteInt32(t.Partitions[i])
	}
}

func (t *PartitionAssignment) Unmarshal(r *Reader) {
	t.Topic = r.ReadString()
	t.Partitions = make([]int32, int(r.ReadInt32()))
	for i := range t.Partitions {
		t.Partitions[i] = r.ReadInt32()
	}
}

func (t *SyncGroupResponse) Marshal(w *Writer) {
	t.ErrorCode.Marshal(w)
	t.MemberAssignment.Marshal(w)
}

func (t *SyncGroupResponse) Unmarshal(r *Reader) {
	t.ErrorCode.Unmarshal(r)
	t.MemberAssignment.Unmarshal(r)
}

func (t *HeartbeatRequest) Marshal(w *Writer) {
	w.WriteString(t.GroupID)
	w.WriteInt32(t.GenerationID)
	w.WriteString(t.MemberID)
}

func (t *HeartbeatRequest) Unmarshal(r *Reader) {
	t.GroupID = r.ReadString()
	t.GenerationID = r.ReadInt32()
	t.MemberID = r.ReadString()
}

func (t *HeartbeatResponse) Marshal(w *Writer) {
	(*t).Marshal(w)
}

func (t *HeartbeatResponse) Unmarshal(r *Reader) {
	(*t).Unmarshal(r)
}

func (t *LeaveGroupRequest) Marshal(w *Writer) {
	w.WriteString(t.GroupID)
	w.WriteString(t.MemberID)
}

func (t *LeaveGroupRequest) Unmarshal(r *Reader) {
	t.GroupID = r.ReadString()
	t.MemberID = r.ReadString()
}

func (t *LeaveGroupResponse) Marshal(w *Writer) {
	(*t).Marshal(w)
}

func (t *LeaveGroupResponse) Unmarshal(r *Reader) {
	(*t).Unmarshal(r)
}

func (t *ListGroupsRequest) Marshal(w *Writer) {
	// no fields for type ListGroupsRequest, {struct  [] map[]}
}

func (t *ListGroupsRequest) Unmarshal(r *Reader) {
	// no fields for type ListGroupsRequest, {struct  [] map[]}
}

func (t *ListGroupsResponse) Marshal(w *Writer) {
	t.ErrorCode.Marshal(w)
	t.Groups.Marshal(w)
}

func (t *ListGroupsResponse) Unmarshal(r *Reader) {
	t.ErrorCode.Unmarshal(r)
	t.Groups.Unmarshal(r)
}

func (t *Groups) Marshal(w *Writer) {
	w.WriteInt32(int32(len((*t))))
	for i := range *t {
		(*t)[i].Marshal(w)
	}
}

func (t *Groups) Unmarshal(r *Reader) {
	(*t) = make([]Group, int(r.ReadInt32()))
	for i := range *t {
		(*t)[i].Unmarshal(r)
	}
}

func (t *Group) Marshal(w *Writer) {
	w.WriteString(t.GroupID)
	w.WriteString(t.ProtocolType)
}

func (t *Group) Unmarshal(r *Reader) {
	t.GroupID = r.ReadString()
	t.ProtocolType = r.ReadString()
}

func (t *DescribeGroupsRequest) Marshal(w *Writer) {
	w.WriteInt32(int32(len((*t))))
	for i := range *t {
		w.WriteString((*t)[i])
	}
}

func (t *DescribeGroupsRequest) Unmarshal(r *Reader) {
	(*t) = make([]string, int(r.ReadInt32()))
	for i := range *t {
		(*t)[i] = r.ReadString()
	}
}

func (t *DescribeGroupsResponse) Marshal(w *Writer) {
	w.WriteInt32(int32(len((*t))))
	for i := range *t {
		(*t)[i].Marshal(w)
	}
}

func (t *DescribeGroupsResponse) Unmarshal(r *Reader) {
	(*t) = make([]GroupDescription, int(r.ReadInt32()))
	for i := range *t {
		(*t)[i].Unmarshal(r)
	}
}

func (t *GroupDescription) Marshal(w *Writer) {
	t.ErrorCode.Marshal(w)
	w.WriteString(t.GroupID)
	w.WriteString(t.State)
	w.WriteString(t.ProtocolType)
	w.WriteString(t.Protocol)
	t.Members.Marshal(w)
}

func (t *GroupDescription) Unmarshal(r *Reader) {
	t.ErrorCode.Unmarshal(r)
	t.GroupID = r.ReadString()
	t.State = r.ReadString()
	t.ProtocolType = r.ReadString()
	t.Protocol = r.ReadString()
	t.Members.Unmarshal(r)
}

func (t *Members) Marshal(w *Writer) {
	w.WriteInt32(int32(len((*t))))
	for i := range *t {
		(*t)[i].Marshal(w)
	}
}

func (t *Members) Unmarshal(r *Reader) {
	(*t) = make([]Member, int(r.ReadInt32()))
	for i := range *t {
		(*t)[i].Unmarshal(r)
	}
}

func (t *Member) Marshal(w *Writer) {
	w.WriteString(t.MemberID)
	w.WriteString(t.ClientID)
	w.WriteString(t.ClientHost)
	w.WriteBytes(t.MemberMetadata)
	t.MemberAssignment.Marshal(w)
}

func (t *Member) Unmarshal(r *Reader) {
	t.MemberID = r.ReadString()
	t.ClientID = r.ReadString()
	t.ClientHost = r.ReadString()
	t.MemberMetadata = r.ReadBytes()
	t.MemberAssignment.Unmarshal(r)
}

func (t *ErrorCode) Marshal(w *Writer) {
	w.WriteInt16(int16((*t)))
}

func (t *ErrorCode) Unmarshal(r *Reader) {
	(*t) = ErrorCode(r.ReadInt16())
}
