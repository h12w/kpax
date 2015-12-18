package consumer

import (
	"errors"
	"fmt"
	"time"

	"h12.me/kafka/broker"
	"h12.me/kafka/cluster"
)

var (
	ErrOffsetNotFound   = errors.New("offset not found")
	ErrFailCommitOffset = errors.New("fail to commit offset")
)

type Config struct {
	Client          cluster.Config
	MaxWaitTime     time.Duration
	MinBytes        int
	MaxBytes        int
	OffsetRetention time.Duration
}

func DefaultConfig(brokers ...string) *Config {
	return &Config{
		Client:          *cluster.DefaultConfig(brokers...),
		MaxWaitTime:     100 * time.Millisecond,
		MinBytes:        1,
		MaxBytes:        1024 * 1024,
		OffsetRetention: 7 * 24 * time.Hour,
	}
}

type C struct {
	cluster *cluster.C
	config  *Config
}

func New(config *Config) (*C, error) {
	cluster, err := cluster.New(&config.Client)
	if err != nil {
		return nil, err
	}
	return &C{
		cluster: cluster,
		config:  config,
	}, nil
}

func (c *C) Offset(topic string, partition int32, consumerGroup string) (int64, error) {
	req := &broker.Request{
		RequestMessage: &broker.OffsetFetchRequestV1{
			ConsumerGroup: consumerGroup,
			PartitionInTopics: []broker.PartitionInTopic{
				{
					TopicName:  topic,
					Partitions: []int32{partition},
				},
			},
		},
	}
	resp := broker.OffsetFetchResponse{}
	coord, err := c.cluster.Coordinator(topic, consumerGroup)
	if err != nil {
		return 0, err
	}
	if err := coord.Do(req, &resp); err != nil {
		return 0, err
	}
	for i := range resp {
		t := &resp[i]
		if t.TopicName == topic {
			for j := range resp[i].OffsetMetadataInPartitions {
				p := &t.OffsetMetadataInPartitions[j]
				if p.ErrorCode.HasError() {
					return 0, fmt.Errorf("fail to get offset for (%s, %d): %v", topic, p.Partition, p.ErrorCode)
				}
				return p.Offset, nil
			}
		}
	}
	return 0, ErrOffsetNotFound
}

func (c *C) Consume(topic string, partition int32, offset int64) (values [][]byte, err error) {
	req := &broker.Request{
		RequestMessage: &broker.FetchRequest{
			ReplicaID:   -1,
			MaxWaitTime: int32(c.config.MaxWaitTime / time.Millisecond),
			MinBytes:    int32(c.config.MinBytes),
			FetchOffsetInTopics: []broker.FetchOffsetInTopic{
				{
					TopicName: topic,
					FetchOffsetInPartitions: []broker.FetchOffsetInPartition{
						{
							Partition:   partition,
							FetchOffset: offset,
							MaxBytes:    int32(c.config.MaxBytes),
						},
					},
				},
			},
		},
	}
	resp := broker.FetchResponse{}
	leader, err := c.cluster.Leader(topic, partition)
	if err != nil {
		return nil, err
	}
	if err := leader.Do(req, &resp); err != nil {
		if broker.IsNotLeader(err) {
			c.cluster.LeaderIsDown(topic, partition)
		}
		return nil, err
	}
	for i := range resp {
		t := &resp[i]
		if t.TopicName != topic {
			continue
		}
		for j := range t.FetchMessageSetInPartitions {
			p := &t.FetchMessageSetInPartitions[j]
			if p.Partition != partition {
				continue
			}
			start := 0
			if p.ErrorCode.HasError() {
				return nil, p.ErrorCode
			}
			for k := range p.MessageSet {
				m := &p.MessageSet[k]
				if m.Offset == offset {
					start = k
					break
				}
			}
			for k := range p.MessageSet[start:] {
				m := &p.MessageSet[k]
				values = append(values, m.SizedMessage.CRCMessage.Message.Value)
			}
		}
	}
	return
}

func (c *C) Commit(topic string, partition int32, consumerGroup string, offset int64) error {
	req := &broker.Request{RequestMessage: &broker.OffsetCommitRequestV1{
		ConsumerGroupID: consumerGroup,
		OffsetCommitInTopicV1s: []broker.OffsetCommitInTopicV1{
			{
				TopicName: topic,
				OffsetCommitInPartitionV1s: []broker.OffsetCommitInPartitionV1{
					{
						Partition: partition,
						Offset:    offset,
						// TimeStamp in milliseconds
						TimeStamp: time.Now().Add(c.config.OffsetRetention).Unix() * 1000,
					},
				},
			},
		},
	},
	}
	resp := broker.OffsetCommitResponse{}
	coord, err := c.cluster.Coordinator(topic, consumerGroup)
	if err != nil {
		return err
	}
	if err := coord.Do(req, &resp); err != nil {
		if broker.IsNotCoordinator(err) {
			c.cluster.CoordinatorIsDown(consumerGroup)
		}
		return err
	}
	for i := range resp {
		t := &resp[i]
		if t.TopicName == topic {
			for j := range t.ErrorInPartitions {
				p := &t.ErrorInPartitions[j]
				if p.Partition == partition {
					if p.ErrorCode.HasError() {
						return p.ErrorCode
					}
					return nil
				}
			}
		}
	}
	return ErrFailCommitOffset
}
