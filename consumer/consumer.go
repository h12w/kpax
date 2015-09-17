package consumer

import (
	"errors"
	"fmt"
	"time"

	"h12.me/kafka/client"
	"h12.me/kafka/proto"
)

var (
	ErrOffsetNotFound   = errors.New("offset not found")
	ErrFailCommitOffset = errors.New("fail to commit offset")
)

type Config struct {
	Client      client.Config
	MaxWaitTime int32
	MinBytes    int32
	MaxBytes    int32
}

type C struct {
	client *client.C
	config *Config
}

func New(config *Config) (*C, error) {
	client, err := client.New(&config.Client)
	if err != nil {
		return nil, err
	}
	return &C{
		client: client,
		config: config,
	}, nil
}

func (c *C) Offset(topic string, partition int32, consumerGroup string) (int64, error) {
	req := c.client.NewRequest(&proto.OffsetFetchRequestV1{
		ConsumerGroup: consumerGroup,
		PartitionInTopics: []proto.PartitionInTopic{
			{
				TopicName:  topic,
				Partitions: []int32{partition},
			},
		},
	})
	resp := proto.OffsetFetchResponse{}
	broker, err := c.client.Coordinator(topic, consumerGroup)
	if err != nil {
		return 0, err
	}
	if err := broker.Do(req, &proto.Response{ResponseMessage: &resp}); err != nil {
		return 0, err
	}
	for i := range resp {
		t := &resp[i]
		if t.TopicName == topic {
			for j := range resp[i].OffsetMetadataInPartitions {
				p := &t.OffsetMetadataInPartitions[j]
				if p.ErrorCode != 0 {
					return 0, fmt.Errorf("offset not exist %d", p.ErrorCode)
				}
				return p.Offset, nil
			}
		}
	}
	return 0, ErrOffsetNotFound
}

func (c *C) Consume(topic string, partition int32, offset int64) (values [][]byte, err error) {
	req := c.client.NewRequest(&proto.FetchRequest{
		ReplicaID:   -1,
		MaxWaitTime: c.config.MaxWaitTime,
		MinBytes:    c.config.MinBytes,
		FetchOffsetInTopics: []proto.FetchOffsetInTopic{
			{
				TopicName: topic,
				FetchOffsetInPartitions: []proto.FetchOffsetInPartition{
					{
						Partition:   partition,
						FetchOffset: offset,
						MaxBytes:    c.config.MaxBytes,
					},
				},
			},
		},
	})
	resp := proto.FetchResponse{}
	broker, err := c.client.Leader(topic, partition)
	if err != nil {
		return nil, err
	}
	if err := broker.Do(req, &proto.Response{ResponseMessage: &resp}); err != nil {
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
	req := c.client.NewRequest(&proto.OffsetCommitRequestV1{
		ConsumerGroupID: consumerGroup,
		OffsetCommitInTopicV1s: []proto.OffsetCommitInTopicV1{
			{
				TopicName: topic,
				OffsetCommitInPartitionV1s: []proto.OffsetCommitInPartitionV1{
					{
						Partition: partition,
						Offset:    offset,
						TimeStamp: time.Now().Unix(),
					},
				},
			},
		},
	})
	resp := proto.OffsetCommitResponse{}
	broker, err := c.client.Coordinator(topic, consumerGroup)
	if err != nil {
		return err
	}
	if err := broker.Do(req, &proto.Response{ResponseMessage: &resp}); err != nil {
		return err
	}
	for i := range resp {
		t := &resp[i]
		if t.TopicName == topic {
			for j := range t.ErrorInPartitions {
				p := &t.ErrorInPartitions[j]
				if p.Partition == partition {
					if p.ErrorCode != 0 {
						return fmt.Errorf("fail to commit %d", p.ErrorCode)
					}
					return nil
				}
			}
		}
	}
	return ErrFailCommitOffset
}
