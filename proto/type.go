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
func (r *GroupCoordinatorRequest) APIKey() int16 { return 10 }
func (r *JoinGroupRequest) APIKey() int16        { return 11 }
func (r *HeartbeatRequest) APIKey() int16        { return 12 }
func (r *LeaveGroupRequest) APIKey() int16       { return 13 }
func (r *SyncGroupRequest) APIKey() int16        { return 14 }
func (r *DescribeGroupsRequest) APIKey() int16   { return 15 }
func (r *ListGroupsRequest) APIKey() int16       { return 16 }

func (r *ProduceRequest) APIVersion() int16          { return 0 }
func (r *FetchRequest) APIVersion() int16            { return 0 }
func (r *OffsetRequest) APIVersion() int16           { return 0 }
func (r *TopicMetadataRequest) APIVersion() int16    { return 0 }
func (r *OffsetCommitRequestV0) APIVersion() int16   { return 0 }
func (r *OffsetCommitRequestV1) APIVersion() int16   { return 1 }
func (r *OffsetCommitRequestV2) APIVersion() int16   { return 2 }
func (r *OffsetFetchRequestV0) APIVersion() int16    { return 0 }
func (r *OffsetFetchRequestV1) APIVersion() int16    { return 1 }
func (r *GroupCoordinatorRequest) APIVersion() int16 { return 0 }
func (r *JoinGroupRequest) APIVersion() int16        { return 0 }
func (r *HeartbeatRequest) APIVersion() int16        { return 0 }
func (r *LeaveGroupRequest) APIVersion() int16       { return 0 }
func (r *SyncGroupRequest) APIVersion() int16        { return 0 }
func (r *DescribeGroupsRequest) APIVersion() int16   { return 0 }
func (r *ListGroupsRequest) APIVersion() int16       { return 0 }
