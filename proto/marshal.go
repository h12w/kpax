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

func (t *RequestOrResponse) Unmarshal(r io.Reader) error {
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.Size = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		if err := t.T.Unmarshal(r); err != nil {
			return err
		}
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

func (t *Request) Unmarshal(r io.Reader) error {
	{
		var b [2]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.APIKey = int16(b[0])<<8 | int16(b[1])
	}
	{
		var b [2]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.APIVersion = int16(b[0])<<8 | int16(b[1])
	}
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.CorrelationID = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		var lb [2]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int16(lb[0])<<8 | int16(lb[1])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.ClientID = string(b)
	}
	{
		if err := t.RequestMessage.Unmarshal(r); err != nil {
			return err
		}
	}
	return nil
}

func (t *Response) Marshal(w io.Writer) error {
	if _, err := w.Write([]byte{byte(t.CorrelationID >> 24), byte(t.CorrelationID >> 16), byte(t.CorrelationID >> 8), byte(t.CorrelationID)}); err != nil {
		return err
	}
	if err := t.ResponseMessage.Marshal(w); err != nil {
		return err
	}
	return nil
}

func (t *Response) Unmarshal(r io.Reader) error {
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.CorrelationID = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		if err := t.ResponseMessage.Unmarshal(r); err != nil {
			return err
		}
	}
	return nil
}

func (t *MessageSet) Marshal(w io.Writer) error {
	{
		l := int32(len((*t)))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range *t {
			if err := (*t)[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *MessageSet) Unmarshal(r io.Reader) error {
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		(*t) = make([]SizedMessage, int(l))
		for i := range *t {
			if err := (*t)[i].Unmarshal(r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *SizedMessage) Marshal(w io.Writer) error {
	if _, err := w.Write([]byte{byte(t.Offset >> 56), byte(t.Offset >> 48), byte(t.Offset >> 40), byte(t.Offset) >> 32, byte(t.Offset >> 24), byte(t.Offset >> 16), byte(t.Offset >> 8), byte(t.Offset)}); err != nil {
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

func (t *SizedMessage) Unmarshal(r io.Reader) error {
	{
		var b [8]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.Offset = int64(b[0])<<56 | int64(b[1])<<48 | int64(b[2])<<40 | int64(b[3])<<32 | int64(b[4])<<24 | int64(b[5])<<16 | int64(b[6])<<8 | int64(b[7])
	}
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.MessageSize = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		if err := t.Message.Unmarshal(r); err != nil {
			return err
		}
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

func (t *Message) Unmarshal(r io.Reader) error {
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.CRC = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		var b [1]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.MagicByte = int8(b[0])
	}
	{
		var b [1]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.Attributes = int8(b[0])
	}
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.Key = b
	}
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.Value = b
	}
	return nil
}

func (t *TopicMetadataRequest) Marshal(w io.Writer) error {
	{
		l := int32(len((*t)))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range *t {
			{
				l := int16(len((*t)[i]))
				if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
					return err
				}
				if _, err := w.Write([]byte((*t)[i])); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (t *TopicMetadataRequest) Unmarshal(r io.Reader) error {
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		(*t) = make([]string, int(l))
		for i := range *t {
			{
				var lb [2]byte
				if _, err := r.Read(lb[:]); err != nil {
					return err
				}
				l := int16(lb[0])<<8 | int16(lb[1])
				b := make([]byte, l)
				if _, err := r.Read(b); err != nil {
					return err
				}
				(*t)[i] = string(b)
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

func (t *MetadataResponse) Unmarshal(r io.Reader) error {
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		t.Brokers = make([]Broker, int(l))
		for i := range t.Brokers {
			if err := t.Brokers[i].Unmarshal(r); err != nil {
				return err
			}
		}
	}
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		t.TopicMetadatas = make([]TopicMetadata, int(l))
		for i := range t.TopicMetadatas {
			if err := t.TopicMetadatas[i].Unmarshal(r); err != nil {
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

func (t *Broker) Unmarshal(r io.Reader) error {
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.NodeID = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		var lb [2]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int16(lb[0])<<8 | int16(lb[1])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.Host = string(b)
	}
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.Port = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
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

func (t *TopicMetadata) Unmarshal(r io.Reader) error {
	{
		var b [2]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.TopicErrorCode = int16(b[0])<<8 | int16(b[1])
	}
	{
		var lb [2]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int16(lb[0])<<8 | int16(lb[1])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.TopicName = string(b)
	}
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		t.PartitionMetadatas = make([]PartitionMetadata, int(l))
		for i := range t.PartitionMetadatas {
			if err := t.PartitionMetadatas[i].Unmarshal(r); err != nil {
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

func (t *PartitionMetadata) Unmarshal(r io.Reader) error {
	{
		var b [2]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.PartitionErrorCode = int16(b[0])<<8 | int16(b[1])
	}
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.PartitionID = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.Leader = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.Replicas = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.ISR = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
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

func (t *ProduceRequest) Unmarshal(r io.Reader) error {
	{
		var b [2]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.RequiredAcks = int16(b[0])<<8 | int16(b[1])
	}
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.Timeout = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		t.MessageSetInTopics = make([]MessageSetInTopic, int(l))
		for i := range t.MessageSetInTopics {
			if err := t.MessageSetInTopics[i].Unmarshal(r); err != nil {
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

func (t *MessageSetInTopic) Unmarshal(r io.Reader) error {
	{
		var lb [2]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int16(lb[0])<<8 | int16(lb[1])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.TopicName = string(b)
	}
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		t.MessageSetInPartitions = make([]MessageSetInPartition, int(l))
		for i := range t.MessageSetInPartitions {
			if err := t.MessageSetInPartitions[i].Unmarshal(r); err != nil {
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

func (t *MessageSetInPartition) Unmarshal(r io.Reader) error {
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.Partition = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.MessageSetSize = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		if err := t.MessageSet.Unmarshal(r); err != nil {
			return err
		}
	}
	return nil
}

func (t *ProduceResponse) Marshal(w io.Writer) error {
	{
		l := int32(len((*t)))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range *t {
			if err := (*t)[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *ProduceResponse) Unmarshal(r io.Reader) error {
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		(*t) = make([]OffsetInTopic, int(l))
		for i := range *t {
			if err := (*t)[i].Unmarshal(r); err != nil {
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

func (t *OffsetInTopic) Unmarshal(r io.Reader) error {
	{
		var lb [2]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int16(lb[0])<<8 | int16(lb[1])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.TopicName = string(b)
	}
	{
		if err := t.OffsetInPartitions.Unmarshal(r); err != nil {
			return err
		}
	}
	return nil
}

func (t *OffsetInPartitions) Marshal(w io.Writer) error {
	{
		l := int32(len((*t)))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range *t {
			if err := (*t)[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *OffsetInPartitions) Unmarshal(r io.Reader) error {
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		(*t) = make([]OffsetInPartition, int(l))
		for i := range *t {
			if err := (*t)[i].Unmarshal(r); err != nil {
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
	if _, err := w.Write([]byte{byte(t.Offset >> 56), byte(t.Offset >> 48), byte(t.Offset >> 40), byte(t.Offset) >> 32, byte(t.Offset >> 24), byte(t.Offset >> 16), byte(t.Offset >> 8), byte(t.Offset)}); err != nil {
		return err
	}
	return nil
}

func (t *OffsetInPartition) Unmarshal(r io.Reader) error {
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.Partition = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		var b [2]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.ErrorCode = int16(b[0])<<8 | int16(b[1])
	}
	{
		var b [8]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.Offset = int64(b[0])<<56 | int64(b[1])<<48 | int64(b[2])<<40 | int64(b[3])<<32 | int64(b[4])<<24 | int64(b[5])<<16 | int64(b[6])<<8 | int64(b[7])
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

func (t *FetchRequest) Unmarshal(r io.Reader) error {
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.ReplicaID = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.MaxWaitTime = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.MinBytes = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		t.FetchOffsetInTopics = make([]FetchOffsetInTopic, int(l))
		for i := range t.FetchOffsetInTopics {
			if err := t.FetchOffsetInTopics[i].Unmarshal(r); err != nil {
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

func (t *FetchOffsetInTopic) Unmarshal(r io.Reader) error {
	{
		var lb [2]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int16(lb[0])<<8 | int16(lb[1])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.TopicName = string(b)
	}
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		t.FetchOffsetInPartitions = make([]FetchOffsetInPartition, int(l))
		for i := range t.FetchOffsetInPartitions {
			if err := t.FetchOffsetInPartitions[i].Unmarshal(r); err != nil {
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
	if _, err := w.Write([]byte{byte(t.FetchOffset >> 56), byte(t.FetchOffset >> 48), byte(t.FetchOffset >> 40), byte(t.FetchOffset) >> 32, byte(t.FetchOffset >> 24), byte(t.FetchOffset >> 16), byte(t.FetchOffset >> 8), byte(t.FetchOffset)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.MaxBytes >> 24), byte(t.MaxBytes >> 16), byte(t.MaxBytes >> 8), byte(t.MaxBytes)}); err != nil {
		return err
	}
	return nil
}

func (t *FetchOffsetInPartition) Unmarshal(r io.Reader) error {
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.Partition = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		var b [8]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.FetchOffset = int64(b[0])<<56 | int64(b[1])<<48 | int64(b[2])<<40 | int64(b[3])<<32 | int64(b[4])<<24 | int64(b[5])<<16 | int64(b[6])<<8 | int64(b[7])
	}
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.MaxBytes = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	return nil
}

func (t *FetchResponse) Marshal(w io.Writer) error {
	{
		l := int32(len((*t)))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range *t {
			if err := (*t)[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *FetchResponse) Unmarshal(r io.Reader) error {
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		(*t) = make([]FetchMessageSetInTopic, int(l))
		for i := range *t {
			if err := (*t)[i].Unmarshal(r); err != nil {
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

func (t *FetchMessageSetInTopic) Unmarshal(r io.Reader) error {
	{
		var lb [2]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int16(lb[0])<<8 | int16(lb[1])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.TopicName = string(b)
	}
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		t.FetchMessageSetInPartitions = make([]FetchMessageSetInPartition, int(l))
		for i := range t.FetchMessageSetInPartitions {
			if err := t.FetchMessageSetInPartitions[i].Unmarshal(r); err != nil {
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
	if _, err := w.Write([]byte{byte(t.HighwaterMarkOffset >> 56), byte(t.HighwaterMarkOffset >> 48), byte(t.HighwaterMarkOffset >> 40), byte(t.HighwaterMarkOffset) >> 32, byte(t.HighwaterMarkOffset >> 24), byte(t.HighwaterMarkOffset >> 16), byte(t.HighwaterMarkOffset >> 8), byte(t.HighwaterMarkOffset)}); err != nil {
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

func (t *FetchMessageSetInPartition) Unmarshal(r io.Reader) error {
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.Partition = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		var b [2]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.ErrorCode = int16(b[0])<<8 | int16(b[1])
	}
	{
		var b [8]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.HighwaterMarkOffset = int64(b[0])<<56 | int64(b[1])<<48 | int64(b[2])<<40 | int64(b[3])<<32 | int64(b[4])<<24 | int64(b[5])<<16 | int64(b[6])<<8 | int64(b[7])
	}
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.MessageSetSize = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		if err := t.MessageSet.Unmarshal(r); err != nil {
			return err
		}
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

func (t *OffsetRequest) Unmarshal(r io.Reader) error {
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.ReplicaID = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		t.TimeInTopics = make([]TimeInTopic, int(l))
		for i := range t.TimeInTopics {
			if err := t.TimeInTopics[i].Unmarshal(r); err != nil {
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

func (t *TimeInTopic) Unmarshal(r io.Reader) error {
	{
		var lb [2]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int16(lb[0])<<8 | int16(lb[1])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.TopicName = string(b)
	}
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		t.TimeInPartitions = make([]TimeInPartition, int(l))
		for i := range t.TimeInPartitions {
			if err := t.TimeInPartitions[i].Unmarshal(r); err != nil {
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
	if _, err := w.Write([]byte{byte(t.Time >> 56), byte(t.Time >> 48), byte(t.Time >> 40), byte(t.Time) >> 32, byte(t.Time >> 24), byte(t.Time >> 16), byte(t.Time >> 8), byte(t.Time)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.MaxNumberOfOffsets >> 24), byte(t.MaxNumberOfOffsets >> 16), byte(t.MaxNumberOfOffsets >> 8), byte(t.MaxNumberOfOffsets)}); err != nil {
		return err
	}
	return nil
}

func (t *TimeInPartition) Unmarshal(r io.Reader) error {
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.Partition = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		var b [8]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.Time = int64(b[0])<<56 | int64(b[1])<<48 | int64(b[2])<<40 | int64(b[3])<<32 | int64(b[4])<<24 | int64(b[5])<<16 | int64(b[6])<<8 | int64(b[7])
	}
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.MaxNumberOfOffsets = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	return nil
}

func (t *OffsetResponse) Marshal(w io.Writer) error {
	{
		l := int32(len((*t)))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range *t {
			if err := (*t)[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *OffsetResponse) Unmarshal(r io.Reader) error {
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		(*t) = make([]OffsetsInTopic, int(l))
		for i := range *t {
			if err := (*t)[i].Unmarshal(r); err != nil {
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

func (t *OffsetsInTopic) Unmarshal(r io.Reader) error {
	{
		var lb [2]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int16(lb[0])<<8 | int16(lb[1])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.TopicName = string(b)
	}
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		t.OffsetsInPartitions = make([]OffsetsInPartition, int(l))
		for i := range t.OffsetsInPartitions {
			if err := t.OffsetsInPartitions[i].Unmarshal(r); err != nil {
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
			if _, err := w.Write([]byte{byte(t.Offsets[i] >> 56), byte(t.Offsets[i] >> 48), byte(t.Offsets[i] >> 40), byte(t.Offsets[i]) >> 32, byte(t.Offsets[i] >> 24), byte(t.Offsets[i] >> 16), byte(t.Offsets[i] >> 8), byte(t.Offsets[i])}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *OffsetsInPartition) Unmarshal(r io.Reader) error {
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.Partition = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		var b [2]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.ErrorCode = int16(b[0])<<8 | int16(b[1])
	}
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		t.Offsets = make([]int64, int(l))
		for i := range t.Offsets {
			{
				var b [8]byte
				if _, err := r.Read(b[:]); err != nil {
					return err
				}
				t.Offsets[i] = int64(b[0])<<56 | int64(b[1])<<48 | int64(b[2])<<40 | int64(b[3])<<32 | int64(b[4])<<24 | int64(b[5])<<16 | int64(b[6])<<8 | int64(b[7])
			}
		}
	}
	return nil
}

func (t *ConsumerMetadataRequest) Marshal(w io.Writer) error {
	{
		l := int16(len((*t)))
		if _, err := w.Write([]byte{byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		if _, err := w.Write([]byte((*t))); err != nil {
			return err
		}
	}
	return nil
}

func (t *ConsumerMetadataRequest) Unmarshal(r io.Reader) error {
	// type ConsumerMetadataRequest, { string []}
	{
		var lb [2]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int16(lb[0])<<8 | int16(lb[1])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		(*t) = ConsumerMetadataRequest(b)
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

func (t *ConsumerMetadataResponse) Unmarshal(r io.Reader) error {
	{
		var b [2]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.ErrorCode = int16(b[0])<<8 | int16(b[1])
	}
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.CoordinatorID = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		var lb [2]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int16(lb[0])<<8 | int16(lb[1])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.CoordinatorHost = string(b)
	}
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.CoordinatorPort = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
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

func (t *OffsetCommitRequestV0) Unmarshal(r io.Reader) error {
	{
		var lb [2]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int16(lb[0])<<8 | int16(lb[1])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.ConsumerGroupID = string(b)
	}
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		t.OffsetCommitInTopicV0s = make([]OffsetCommitInTopicV0, int(l))
		for i := range t.OffsetCommitInTopicV0s {
			if err := t.OffsetCommitInTopicV0s[i].Unmarshal(r); err != nil {
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

func (t *OffsetCommitInTopicV0) Unmarshal(r io.Reader) error {
	{
		var lb [2]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int16(lb[0])<<8 | int16(lb[1])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.TopicName = string(b)
	}
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		t.OffsetCommitInPartitionV0s = make([]OffsetCommitInPartitionV0, int(l))
		for i := range t.OffsetCommitInPartitionV0s {
			if err := t.OffsetCommitInPartitionV0s[i].Unmarshal(r); err != nil {
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
	if _, err := w.Write([]byte{byte(t.Offset >> 56), byte(t.Offset >> 48), byte(t.Offset >> 40), byte(t.Offset) >> 32, byte(t.Offset >> 24), byte(t.Offset >> 16), byte(t.Offset >> 8), byte(t.Offset)}); err != nil {
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

func (t *OffsetCommitInPartitionV0) Unmarshal(r io.Reader) error {
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.Partition = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		var b [8]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.Offset = int64(b[0])<<56 | int64(b[1])<<48 | int64(b[2])<<40 | int64(b[3])<<32 | int64(b[4])<<24 | int64(b[5])<<16 | int64(b[6])<<8 | int64(b[7])
	}
	{
		var lb [2]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int16(lb[0])<<8 | int16(lb[1])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.Metadata = string(b)
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

func (t *OffsetCommitRequestV1) Unmarshal(r io.Reader) error {
	{
		var lb [2]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int16(lb[0])<<8 | int16(lb[1])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.ConsumerGroupID = string(b)
	}
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.ConsumerGroupGenerationID = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		var lb [2]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int16(lb[0])<<8 | int16(lb[1])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.ConsumerID = string(b)
	}
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		t.OffsetCommitInTopicV1s = make([]OffsetCommitInTopicV1, int(l))
		for i := range t.OffsetCommitInTopicV1s {
			if err := t.OffsetCommitInTopicV1s[i].Unmarshal(r); err != nil {
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

func (t *OffsetCommitInTopicV1) Unmarshal(r io.Reader) error {
	{
		var lb [2]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int16(lb[0])<<8 | int16(lb[1])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.TopicName = string(b)
	}
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		t.OffsetCommitInPartitionV1s = make([]OffsetCommitInPartitionV1, int(l))
		for i := range t.OffsetCommitInPartitionV1s {
			if err := t.OffsetCommitInPartitionV1s[i].Unmarshal(r); err != nil {
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
	if _, err := w.Write([]byte{byte(t.Offset >> 56), byte(t.Offset >> 48), byte(t.Offset >> 40), byte(t.Offset) >> 32, byte(t.Offset >> 24), byte(t.Offset >> 16), byte(t.Offset >> 8), byte(t.Offset)}); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(t.TimeStamp >> 56), byte(t.TimeStamp >> 48), byte(t.TimeStamp >> 40), byte(t.TimeStamp) >> 32, byte(t.TimeStamp >> 24), byte(t.TimeStamp >> 16), byte(t.TimeStamp >> 8), byte(t.TimeStamp)}); err != nil {
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

func (t *OffsetCommitInPartitionV1) Unmarshal(r io.Reader) error {
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.Partition = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		var b [8]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.Offset = int64(b[0])<<56 | int64(b[1])<<48 | int64(b[2])<<40 | int64(b[3])<<32 | int64(b[4])<<24 | int64(b[5])<<16 | int64(b[6])<<8 | int64(b[7])
	}
	{
		var b [8]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.TimeStamp = int64(b[0])<<56 | int64(b[1])<<48 | int64(b[2])<<40 | int64(b[3])<<32 | int64(b[4])<<24 | int64(b[5])<<16 | int64(b[6])<<8 | int64(b[7])
	}
	{
		var lb [2]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int16(lb[0])<<8 | int16(lb[1])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.Metadata = string(b)
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
	if _, err := w.Write([]byte{byte(t.RetentionTime >> 56), byte(t.RetentionTime >> 48), byte(t.RetentionTime >> 40), byte(t.RetentionTime) >> 32, byte(t.RetentionTime >> 24), byte(t.RetentionTime >> 16), byte(t.RetentionTime >> 8), byte(t.RetentionTime)}); err != nil {
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

func (t *OffsetCommitRequestV2) Unmarshal(r io.Reader) error {
	{
		var lb [2]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int16(lb[0])<<8 | int16(lb[1])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.ConsumerGroup = string(b)
	}
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.ConsumerGroupGenerationID = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		var lb [2]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int16(lb[0])<<8 | int16(lb[1])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.ConsumerID = string(b)
	}
	{
		var b [8]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.RetentionTime = int64(b[0])<<56 | int64(b[1])<<48 | int64(b[2])<<40 | int64(b[3])<<32 | int64(b[4])<<24 | int64(b[5])<<16 | int64(b[6])<<8 | int64(b[7])
	}
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		t.OffsetCommitInTopicV2s = make([]OffsetCommitInTopicV2, int(l))
		for i := range t.OffsetCommitInTopicV2s {
			if err := t.OffsetCommitInTopicV2s[i].Unmarshal(r); err != nil {
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

func (t *OffsetCommitInTopicV2) Unmarshal(r io.Reader) error {
	{
		var lb [2]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int16(lb[0])<<8 | int16(lb[1])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.TopicName = string(b)
	}
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		t.OffsetCommitInPartitionV2s = make([]OffsetCommitInPartitionV2, int(l))
		for i := range t.OffsetCommitInPartitionV2s {
			if err := t.OffsetCommitInPartitionV2s[i].Unmarshal(r); err != nil {
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
	if _, err := w.Write([]byte{byte(t.Offset >> 56), byte(t.Offset >> 48), byte(t.Offset >> 40), byte(t.Offset) >> 32, byte(t.Offset >> 24), byte(t.Offset >> 16), byte(t.Offset >> 8), byte(t.Offset)}); err != nil {
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

func (t *OffsetCommitInPartitionV2) Unmarshal(r io.Reader) error {
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.Partition = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		var b [8]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.Offset = int64(b[0])<<56 | int64(b[1])<<48 | int64(b[2])<<40 | int64(b[3])<<32 | int64(b[4])<<24 | int64(b[5])<<16 | int64(b[6])<<8 | int64(b[7])
	}
	{
		var lb [2]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int16(lb[0])<<8 | int16(lb[1])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.Metadata = string(b)
	}
	return nil
}

func (t *OffsetCommitResponse) Marshal(w io.Writer) error {
	{
		l := int32(len((*t)))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range *t {
			if err := (*t)[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *OffsetCommitResponse) Unmarshal(r io.Reader) error {
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		(*t) = make([]ErrorInTopic, int(l))
		for i := range *t {
			if err := (*t)[i].Unmarshal(r); err != nil {
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

func (t *ErrorInTopic) Unmarshal(r io.Reader) error {
	{
		var lb [2]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int16(lb[0])<<8 | int16(lb[1])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.TopicName = string(b)
	}
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		t.ErrorInPartitions = make([]ErrorInPartition, int(l))
		for i := range t.ErrorInPartitions {
			if err := t.ErrorInPartitions[i].Unmarshal(r); err != nil {
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

func (t *ErrorInPartition) Unmarshal(r io.Reader) error {
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.Partition = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		var b [2]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.ErrorCode = int16(b[0])<<8 | int16(b[1])
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

func (t *OffsetFetchRequest) Unmarshal(r io.Reader) error {
	{
		var lb [2]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int16(lb[0])<<8 | int16(lb[1])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.ConsumerGroup = string(b)
	}
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		t.PartitionInTopics = make([]PartitionInTopic, int(l))
		for i := range t.PartitionInTopics {
			if err := t.PartitionInTopics[i].Unmarshal(r); err != nil {
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

func (t *PartitionInTopic) Unmarshal(r io.Reader) error {
	{
		var lb [2]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int16(lb[0])<<8 | int16(lb[1])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.TopicName = string(b)
	}
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		t.Partitions = make([]int32, int(l))
		for i := range t.Partitions {
			{
				var b [4]byte
				if _, err := r.Read(b[:]); err != nil {
					return err
				}
				t.Partitions[i] = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
			}
		}
	}
	return nil
}

func (t *OffsetFetchResponse) Marshal(w io.Writer) error {
	{
		l := int32(len((*t)))
		if _, err := w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}); err != nil {
			return err
		}
		for i := range *t {
			if err := (*t)[i].Marshal(w); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *OffsetFetchResponse) Unmarshal(r io.Reader) error {
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		(*t) = make([]OffsetMetadataInTopic, int(l))
		for i := range *t {
			if err := (*t)[i].Unmarshal(r); err != nil {
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

func (t *OffsetMetadataInTopic) Unmarshal(r io.Reader) error {
	{
		var lb [2]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int16(lb[0])<<8 | int16(lb[1])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.TopicName = string(b)
	}
	{
		var lb [4]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int32(lb[0])<<24 | int32(lb[1])<<16 | int32(lb[2])<<8 | int32(lb[3])
		t.OffsetMetadataInPartitions = make([]OffsetMetadataInPartition, int(l))
		for i := range t.OffsetMetadataInPartitions {
			if err := t.OffsetMetadataInPartitions[i].Unmarshal(r); err != nil {
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
	if _, err := w.Write([]byte{byte(t.Offset >> 56), byte(t.Offset >> 48), byte(t.Offset >> 40), byte(t.Offset) >> 32, byte(t.Offset >> 24), byte(t.Offset >> 16), byte(t.Offset >> 8), byte(t.Offset)}); err != nil {
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

func (t *OffsetMetadataInPartition) Unmarshal(r io.Reader) error {
	{
		var b [4]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.Partition = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	}
	{
		var b [8]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.Offset = int64(b[0])<<56 | int64(b[1])<<48 | int64(b[2])<<40 | int64(b[3])<<32 | int64(b[4])<<24 | int64(b[5])<<16 | int64(b[6])<<8 | int64(b[7])
	}
	{
		var lb [2]byte
		if _, err := r.Read(lb[:]); err != nil {
			return err
		}
		l := int16(lb[0])<<8 | int16(lb[1])
		b := make([]byte, l)
		if _, err := r.Read(b); err != nil {
			return err
		}
		t.Metadata = string(b)
	}
	{
		var b [2]byte
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		t.ErrorCode = int16(b[0])<<8 | int16(b[1])
	}
	return nil
}
