package broker

import (
	"time"
)

var (
	Earliest = time.Time{}
	Latest   = time.Time{}.Add(time.Nanosecond)
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

func (b *B) OffsetByTime(topic string, partition int32, t time.Time) (*OffsetResponse, error) {
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
		return nil, err
	}
	return respMsg, nil
}
