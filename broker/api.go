package broker

import (
	"fmt"
	"time"
)

var (
	Earliest = time.Time{}
	Latest   = time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC)
)

func (b *B) TopicMetadata(topics ...string) (*TopicMetadataResponse, error) {
	req := TopicMetadataRequest(topics)
	resp := TopicMetadataResponse{}
	if err := b.Do(&req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (b *B) GroupCoordinator(group string) (*Broker, error) {
	req := GroupCoordinatorRequest(group)
	resp := GroupCoordinatorResponse{}
	if err := b.Do(&req, &resp); err != nil {
		return nil, err
	}
	if resp.ErrorCode.HasError() {
		return nil, resp.ErrorCode
	}
	return &resp.Broker, nil
}

type Produce struct {
	Topic        string
	Partition    int32
	MessageSet   MessageSet
	RequiredAcks int16
	AckTimeout   time.Duration
}

func (p *Produce) Exec(b *B) error {
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
	if err := b.Do(&req, &resp); err != nil {
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

type Consume struct {
	Topic       string
	Partition   int32
	Offset      int64
	MinBytes    int
	MaxBytes    int
	MaxWaitTime time.Duration
}

func (fr *Consume) Exec(c *B) (messages MessageSet, err error) {
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
	if err := c.Do(&req, &resp); err != nil {
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

func (b *B) SegmentOffset(topic string, partition int32, t time.Time) (int64, error) {
	var milliSec int64
	switch t {
	case Latest:
		milliSec = -1
	case Earliest:
		milliSec = -2
	default:
		milliSec = t.UnixNano() / 1000000
	}
	req := OffsetRequest{
		ReplicaID: -1,
		TimeInTopics: []TimeInTopic{
			{
				TopicName: topic,
				TimeInPartitions: []TimeInPartition{
					{
						Partition:          partition,
						Time:               milliSec,
						MaxNumberOfOffsets: 1,
					},
				},
			},
		},
	}
	resp := OffsetResponse{}
	if err := b.Do(&req, &resp); err != nil {
		return -1, err
	}
	for _, t := range resp {
		if t.TopicName != topic {
			continue
		}
		for _, p := range t.OffsetsInPartitions {
			if p.Partition != partition {
				continue
			}
			if p.ErrorCode.HasError() {
				return -1, p.ErrorCode
			}
			if len(p.Offsets) == 0 {
				return -1, fmt.Errorf("failt to fetch offset for %s, %d", topic, partition)
			}
			return p.Offsets[0], nil
		}
	}
	return -1, fmt.Errorf("failt to fetch offset for %s, %d", topic, partition)
}

type OffsetCommit struct {
	Topic     string
	Partition int32
	Group     string
	Offset    int64
	Retention time.Duration
}

func (commit *OffsetCommit) Exec(b *B) error {
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
	if err := b.Do(&req, &resp); err != nil {
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
