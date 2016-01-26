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
	reqMsg := TopicMetadataRequest(topics)
	req := &Request{
		RequestMessage: &reqMsg,
	}
	respMsg := &TopicMetadataResponse{}
	if err := b.Do(req, respMsg); err != nil {
		return nil, err
	}
	return respMsg, nil
}

func (b *B) GroupCoordinator(group string) (*GroupCoordinatorResponse, error) {
	reqMsg := GroupCoordinatorRequest(group)
	req := &Request{
		RequestMessage: &reqMsg,
	}
	respMsg := &GroupCoordinatorResponse{}
	if err := b.Do(req, respMsg); err != nil {
		return nil, err
	}
	if respMsg.ErrorCode.HasError() {
		return nil, respMsg.ErrorCode
	}
	return respMsg, nil
}

// TODO: produce multiple topics
func (b *B) Produce(topic string, partition int32, messageSet MessageSet) (*ProduceResponse, error) {
	cfg := &b.config.Producer
	req := &Request{
		RequestMessage: &ProduceRequest{
			RequiredAcks: cfg.RequiredAcks,
			Timeout:      int32(cfg.Timeout / time.Millisecond),
			MessageSetInTopics: []MessageSetInTopic{
				{
					TopicName: topic,
					MessageSetInPartitions: []MessageSetInPartition{
						{
							Partition:  partition,
							MessageSet: messageSet,
						},
					},
				},
			},
		},
	}

	respMsg := &ProduceResponse{}
	if err := b.Do(req, respMsg); err != nil {
		return nil, err
	}
	return respMsg, nil
}

func (b *B) OffsetByTime(topic string, partition int32, t time.Time) (int64, error) {
	var milliSec int64
	switch t {
	case Latest:
		milliSec = -1
	case Earliest:
		milliSec = -2
	default:
		milliSec = t.UnixNano() / 1000000
	}
	req := &Request{
		RequestMessage: &OffsetRequest{
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
		},
	}
	respMsg := &OffsetResponse{}
	if err := b.Do(req, respMsg); err != nil {
		return -1, err
	}
	for _, t := range *respMsg {
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

func (b *B) Commit(commit *OffsetCommit) error {
	req := &Request{RequestMessage: &OffsetCommitRequestV1{
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
	}}
	resp := OffsetCommitResponse{}
	if err := b.Do(req, &resp); err != nil {
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
