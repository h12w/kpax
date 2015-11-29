package proto

import (
	"io"
	"strconv"

	"h12.me/wipro"
)

type RequestMessage interface {
	wipro.M
	APIKey() int16
	APIVersion() int16
}

func (*ProduceRequest) APIKey() int16          { return 0 }
func (*FetchRequest) APIKey() int16            { return 1 }
func (*OffsetRequest) APIKey() int16           { return 2 }
func (*TopicMetadataRequest) APIKey() int16    { return 3 }
func (*OffsetCommitRequestV0) APIKey() int16   { return 8 }
func (*OffsetCommitRequestV1) APIKey() int16   { return 8 }
func (*OffsetCommitRequestV2) APIKey() int16   { return 8 }
func (*OffsetFetchRequestV0) APIKey() int16    { return 9 }
func (*OffsetFetchRequestV1) APIKey() int16    { return 9 }
func (*GroupCoordinatorRequest) APIKey() int16 { return 10 }
func (*JoinGroupRequest) APIKey() int16        { return 11 }
func (*HeartbeatRequest) APIKey() int16        { return 12 }
func (*LeaveGroupRequest) APIKey() int16       { return 13 }
func (*SyncGroupRequest) APIKey() int16        { return 14 }
func (*DescribeGroupsRequest) APIKey() int16   { return 15 }
func (*ListGroupsRequest) APIKey() int16       { return 16 }

func (*ProduceRequest) APIVersion() int16          { return 0 }
func (*FetchRequest) APIVersion() int16            { return 0 }
func (*OffsetRequest) APIVersion() int16           { return 0 }
func (*TopicMetadataRequest) APIVersion() int16    { return 0 }
func (*OffsetCommitRequestV0) APIVersion() int16   { return 0 }
func (*OffsetCommitRequestV1) APIVersion() int16   { return 1 }
func (*OffsetCommitRequestV2) APIVersion() int16   { return 2 }
func (*OffsetFetchRequestV0) APIVersion() int16    { return 0 }
func (*OffsetFetchRequestV1) APIVersion() int16    { return 1 }
func (*GroupCoordinatorRequest) APIVersion() int16 { return 0 }
func (*JoinGroupRequest) APIVersion() int16        { return 0 }
func (*HeartbeatRequest) APIVersion() int16        { return 0 }
func (*LeaveGroupRequest) APIVersion() int16       { return 0 }
func (*SyncGroupRequest) APIVersion() int16        { return 0 }
func (*DescribeGroupsRequest) APIVersion() int16   { return 0 }
func (*ListGroupsRequest) APIVersion() int16       { return 0 }

func (b *Broker) Addr() string {
	return b.Host + ":" + strconv.Itoa(int(b.Port))
}

func (req *Request) Send(conn io.Writer) error {
	return wipro.Send(&RequestOrResponse{M: req}, conn)
}

func (resp *Response) Receive(conn io.Reader) error {
	return wipro.Receive(conn, &RequestOrResponse{M: resp})
}
