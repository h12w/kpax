package proto

import (
	"fmt"
	"time"

	"h12.me/kafka/common"
)

var (
	Earliest = time.Time{}
	Latest   = time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC)
)

type client struct {
	id   string
	doer common.Broker
}

func (r *Response) ID() int32     { return r.CorrelationID }
func (r *Request) ID() int32      { return r.CorrelationID }
func (r *Request) SetID(id int32) { r.CorrelationID = id }

func (c client) Do(req RequestMessage, resp ResponseMessage) error {
	return c.doer.Do(
		&Request{
			ClientID:       c.id,
			RequestMessage: req,
		},
		&Response{
			ResponseMessage: resp,
		},
	)
}

const clientID = "h12.me/kafka"

type Metadata []string

func (m Metadata) Fetch(b common.Broker) (*TopicMetadataResponse, error) {
	req := TopicMetadataRequest([]string(m))
	resp := TopicMetadataResponse{}
	if err := (client{clientID, b}).Do(&req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type GroupCoordinator string

func (group GroupCoordinator) Fetch(b common.Broker) (*Broker, error) {
	req := GroupCoordinatorRequest(group)
	resp := GroupCoordinatorResponse{}
	if err := (client{clientID, b}).Do(&req, &resp); err != nil {
		return nil, err
	}
	if resp.ErrorCode.HasError() {
		return nil, resp.ErrorCode
	}
	return &resp.Broker, nil
}

type Payload struct {
	Topic        string
	Partition    int32
	MessageSet   MessageSet
	RequiredAcks int16
	AckTimeout   time.Duration
}

func (p *Payload) Produce(c common.Cluster) error {
	leader, err := c.Leader(p.Topic, p.Partition)
	if err != nil {
		return err
	}
	if err := p.DoProduce(leader); err != nil {
		if IsNotLeader(err) {
			c.LeaderIsDown(p.Topic, p.Partition)
		}
		return err
	}
	return nil
}

func (p *Payload) DoProduce(b common.Broker) error {
	req := ProduceRequest{
		RequiredAcks: p.RequiredAcks,
		Timeout:      int32(p.AckTimeout / time.Millisecond),
		MessageSetInTopics: []MessageSetInTopic{
			{
				TopicName: p.Topic,
				MessageSetInPartitions: []MessageSetInPartition{
					{
						Partition:  p.Partition,
						MessageSet: p.MessageSet,
					},
				},
			},
		},
	}

	resp := ProduceResponse{}
	if err := (client{clientID, b}).Do(&req, &resp); err != nil {
		return err
	}
	for i := range resp {
		t := &resp[i]
		if t.TopicName != p.Topic {
			continue
		}
		for j := range t.OffsetInPartitions {
			pres := &t.OffsetInPartitions[j]
			if pres.Partition != p.Partition {
				continue
			}
			if pres.ErrorCode.HasError() {
				return pres.ErrorCode
			}
			return nil
		}
	}
	return fmt.Errorf("fail to produce to %s, %d", p.Topic, p.Partition)
}

type Messages struct {
	Topic       string
	Partition   int32
	Offset      int64
	MinBytes    int
	MaxBytes    int
	MaxWaitTime time.Duration
}

func (m *Messages) Consume(c common.Cluster) (MessageSet, error) {
	leader, err := c.Leader(m.Topic, m.Partition)
	if err != nil {
		return nil, err
	}
	ms, err := m.DoConsume(leader)
	if err != nil {
		if IsNotLeader(err) {
			c.LeaderIsDown(m.Topic, m.Partition)
		}
		return nil, err
	}
	return ms, nil
}

func (fr *Messages) DoConsume(c common.Broker) (messages MessageSet, err error) {
	req := FetchRequest{
		ReplicaID:   -1,
		MaxWaitTime: int32(fr.MaxWaitTime / time.Millisecond),
		MinBytes:    int32(fr.MinBytes),
		FetchOffsetInTopics: []FetchOffsetInTopic{
			{
				TopicName: fr.Topic,
				FetchOffsetInPartitions: []FetchOffsetInPartition{
					{
						Partition:   fr.Partition,
						FetchOffset: fr.Offset,
						MaxBytes:    int32(fr.MaxBytes),
					},
				},
			},
		},
	}
	resp := FetchResponse{}
	if err := (client{clientID, c}).Do(&req, &resp); err != nil {
		return nil, err
	}
	for i := range resp {
		t := &resp[i]
		if t.TopicName != fr.Topic {
			continue
		}
		for j := range t.FetchMessageSetInPartitions {
			p := &t.FetchMessageSetInPartitions[j]
			if p.Partition != fr.Partition {
				continue
			}
			if p.ErrorCode.HasError() {
				return nil, p.ErrorCode
			}
			ms := p.MessageSet
			ms, err := ms.Flatten()
			if err != nil {
				return nil, err
			}
			for k := range ms {
				m := &ms[k]
				if m.Offset == fr.Offset {
					ms = ms[k:]
					break
				}
			}
			if len(ms) == 0 {
				continue
			}
			if ms[0].Offset != fr.Offset {
				return nil, fmt.Errorf("2: OFFSET MISMATCH %d %d", ms[0].Offset, fr.Offset)
			}
			return ms, nil
		}
	}
	return nil, nil
}

type Offset struct {
	Topic     string
	Partition int32
	Group     string
	Offset    int64
	Retention time.Duration
}

func (o *Offset) Commit(c common.Cluster) error {
	coord, err := c.Coordinator(o.Group)
	if err != nil {
		return err
	}
	if err := o.DoCommit(coord); err != nil {
		if IsNotCoordinator(err) {
			c.CoordinatorIsDown(o.Group)
		}
		return err
	}
	return nil
}

func (commit *Offset) DoCommit(b common.Broker) error {
	req := OffsetCommitRequestV1{
		ConsumerGroupID: commit.Group,
		OffsetCommitInTopicV1s: []OffsetCommitInTopicV1{
			{
				TopicName: commit.Topic,
				OffsetCommitInPartitionV1s: []OffsetCommitInPartitionV1{
					{
						Partition: commit.Partition,
						Offset:    commit.Offset,
						// TimeStamp in milliseconds
						TimeStamp: time.Now().Add(commit.Retention).Unix() * 1000,
					},
				},
			},
		},
	}
	resp := OffsetCommitResponse{}
	if err := (client{clientID, b}).Do(&req, &resp); err != nil {
		return err
	}
	for i := range resp {
		t := &resp[i]
		if t.TopicName == commit.Topic {
			for j := range t.ErrorInPartitions {
				p := &t.ErrorInPartitions[j]
				if p.Partition == commit.Partition {
					if p.ErrorCode.HasError() {
						return p.ErrorCode
					}
					return nil
				}
			}
		}
	}
	return fmt.Errorf("fail to commit offset: %v", commit)
}

func (o *Offset) Fetch(c common.Cluster) (int64, error) {
	coord, err := c.Coordinator(o.Group)
	if err != nil {
		return -1, err
	}
	offset, err := o.DoFetch(coord)
	if err != nil {
		if IsNotCoordinator(err) {
			c.CoordinatorIsDown(o.Group)
		}
		return -1, err
	}
	return offset, nil
}

func (o *Offset) DoFetch(b common.Broker) (int64, error) {
	req := OffsetFetchRequestV1{
		ConsumerGroup: o.Group,
		PartitionInTopics: []PartitionInTopic{
			{
				TopicName:  o.Topic,
				Partitions: []int32{o.Partition},
			},
		},
	}
	resp := OffsetFetchResponse{}
	if err := (client{clientID, b}).Do(&req, &resp); err != nil {
		return -1, err
	}
	for i := range resp {
		t := &resp[i]
		if t.TopicName == o.Topic {
			for j := range resp[i].OffsetMetadataInPartitions {
				p := &t.OffsetMetadataInPartitions[j]
				if p.ErrorCode.HasError() {
					return -1, fmt.Errorf("fail to get offset for (%s, %d): %v", o.Topic, o.Partition, p.ErrorCode)
				}
				return p.Offset, nil
			}
		}
	}

	return -1, fmt.Errorf("fail to get offset for (%s, %d)", o.Topic, o.Partition)
}

type SegmentOffset struct {
	Topic     string
	Partition int32
	Time      time.Time
}

func (o *SegmentOffset) Fetch(c common.Cluster) (int64, error) {
	leader, err := c.Leader(o.Topic, o.Partition)
	if err != nil {
		return -1, err
	}
	offset, err := o.DoFetch(leader)
	if err != nil {
		if IsNotLeader(err) {
			c.LeaderIsDown(o.Topic, o.Partition)
		}
		return -1, err
	}
	return offset, nil
}

func (o *SegmentOffset) DoFetch(b common.Broker) (int64, error) {
	var milliSec int64
	switch o.Time {
	case Latest:
		milliSec = -1
	case Earliest:
		milliSec = -2
	default:
		milliSec = o.Time.UnixNano() / 1000000
	}
	req := OffsetRequest{
		ReplicaID: -1,
		TimeInTopics: []TimeInTopic{
			{
				TopicName: o.Topic,
				TimeInPartitions: []TimeInPartition{
					{
						Partition:          o.Partition,
						Time:               milliSec,
						MaxNumberOfOffsets: 1,
					},
				},
			},
		},
	}
	resp := OffsetResponse{}
	if err := (client{clientID, b}).Do(&req, &resp); err != nil {
		return -1, err
	}
	for _, t := range resp {
		if t.TopicName != o.Topic {
			continue
		}
		for _, p := range t.OffsetsInPartitions {
			if p.Partition != o.Partition {
				continue
			}
			if p.ErrorCode.HasError() {
				return -1, p.ErrorCode
			}
			if len(p.Offsets) == 0 {
				return -1, fmt.Errorf("failt to fetch offset for %s, %d", o.Topic, o.Partition)
			}
			return p.Offsets[0], nil
		}
	}
	return -1, fmt.Errorf("failt to fetch offset for %s, %d", o.Topic, o.Partition)
}
