package proto

type RequestMessage interface {
	M
	APIKey() int16
	APIVersion() int16
}

func (r *ProduceRequest) APIKey() int16          { return 0 }
func (r *FetchRequest) APIKey() int16            { return 1 }
func (r *OffsetRequest) APIKey() int16           { return 2 }
func (r *TopicMetadataRequest) APIKey() int16    { return 3 }
func (r *OffsetCommitRequestV0) APIKey() int16   { return 8 }
func (r *OffsetCommitRequestV1) APIKey() int16   { return 8 }
func (r *OffsetCommitRequestV2) APIKey() int16   { return 8 }
func (r *OffsetFetchRequestV0) APIKey() int16    { return 9 }
func (r *OffsetFetchRequestV1) APIKey() int16    { return 9 }
func (r *OffsetFetchRequestV2) APIKey() int16    { return 9 }
func (r *ConsumerMetadataRequest) APIKey() int16 { return 10 }

func (r *ProduceRequest) APIVersion() int16          { return 0 }
func (r *FetchRequest) APIVersion() int16            { return 0 }
func (r *OffsetRequest) APIVersion() int16           { return 0 }
func (r *TopicMetadataRequest) APIVersion() int16    { return 0 }
func (r *OffsetCommitRequestV0) APIVersion() int16   { return 0 }
func (r *OffsetCommitRequestV1) APIVersion() int16   { return 1 }
func (r *OffsetCommitRequestV2) APIVersion() int16   { return 0 }
func (r *OffsetFetchRequestV0) APIVersion() int16    { return 0 }
func (r *OffsetFetchRequestV1) APIVersion() int16    { return 1 }
func (r *OffsetFetchRequestV2) APIVersion() int16    { return 2 }
func (r *ConsumerMetadataRequest) APIVersion() int16 { return 0 }
