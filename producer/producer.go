package producer

import (
	"errors"
	"math/rand"
	"time"

	"h12.me/kafka/client"
	"h12.me/kafka/proto"
)

var (
	ErrProduceFailed    = errors.New("produce failed")
	ErrNoValidPartition = errors.New("no valid partition")
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Config struct {
	Client       client.Config
	RequiredAcks int16
	Timeout      int32
}

type P struct {
	client           *client.C
	config           *Config
	topicPartitioner *topicPartitioner
}

func New(config *Config) (*P, error) {
	client, err := client.New(&config.Client)
	if err != nil {
		return nil, err
	}
	return &P{
		client:           client,
		config:           config,
		topicPartitioner: newTopicPartitioner(),
	}, nil
}

func (p *P) Produce(topic string, key, value []byte) error {
	partitioner := p.topicPartitioner.Get(topic)
	if partitioner == nil {
		partitions, err := p.client.Partitions(topic)
		if err != nil {
			return err
		}
		partitioner = p.topicPartitioner.Add(topic, partitions)
	}
	messageSet := []proto.OffsetMessage{
		{
			SizedMessage: proto.SizedMessage{CRCMessage: proto.CRCMessage{
				Message: proto.Message{
					Key:   key,
					Value: value,
				},
			}}},
	}
	for i := 0; i < partitioner.Count(); i++ {
		partition, err := partitioner.Partition(key)
		if err != nil {
			p.topicPartitioner.Delete(topic)
			break // all partitions down, failed permanantly
		}
		leader, err := p.client.Leader(topic, partition)
		if err != nil {
			partitioner.Skip(partition)
			// TODO: add log here
			continue
		}
		req := p.client.NewRequest(&proto.ProduceRequest{
			RequiredAcks: p.config.RequiredAcks,
			Timeout:      p.config.Timeout,
			MessageSetInTopics: []proto.MessageSetInTopic{
				{
					TopicName: topic,
					MessageSetInPartitions: []proto.MessageSetInPartition{
						{
							Partition:  partition,
							MessageSet: messageSet,
						},
					},
				},
			},
		})
		resp := proto.ProduceResponse{}
		if err := leader.Do(req, &resp); err != nil {
			partitioner.Skip(partition)
			p.client.LeaderIsDown(topic, partition)
			return err
		}
		for i := range resp {
			for j := range resp[i].OffsetInPartitions {
				if errCode := resp[i].OffsetInPartitions[j].ErrorCode; errCode != 0 {
					// TODO: add log here
					continue
				}
			}
		}
		return nil
	}
	return ErrProduceFailed
}
