package producer

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"h12.me/kafka/cluster"
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
	Cluster            cluster.Config
	LeaderRecoveryTime time.Duration
	RequiredAcks       int16
	AckTimeout         time.Duration
}

func DefaultConfig(brokers ...string) *Config {
	return &Config{
		Cluster:            *cluster.DefaultConfig(brokers...),
		LeaderRecoveryTime: 60 * time.Second,
		RequiredAcks:       proto.AckLocal,
		AckTimeout:         10 * time.Second,
	}
}

type P struct {
	cluster          *cluster.C
	config           *Config
	topicPartitioner *topicPartitioner
}

func New(config *Config) (*P, error) {
	cluster, err := cluster.New(&config.Cluster)
	if err != nil {
		return nil, err
	}
	return &P{
		cluster:          cluster,
		config:           config,
		topicPartitioner: newTopicPartitioner(),
	}, nil
}

func (p *P) Produce(topic string, key, value []byte) error {
	partitioner := p.topicPartitioner.Get(topic)
	if partitioner == nil {
		partitions, err := p.cluster.Partitions(topic)
		if err != nil {
			return err
		}
		partitioner = p.topicPartitioner.Add(topic, partitions)
	}
	messageSet := getMessageSet(key, value)
nextPartition:
	for i := 0; i < partitioner.Count(); i++ {
		partition, err := partitioner.Partition(key)
		if err != nil {
			p.topicPartitioner.Delete(topic)
			return err
		}
		if err := (&proto.Payload{
			Topic:        topic,
			Partition:    partition,
			MessageSet:   messageSet,
			RequiredAcks: p.config.RequiredAcks,
			AckTimeout:   p.config.AckTimeout,
		}).Produce(p.cluster); err != nil {
			partitioner.Skip(partition)
			continue nextPartition
		}
		return nil
	}
	return fmt.Errorf("fail to produce to all partitions in %s", topic)
}

func (p *P) ProduceWithPartition(topic string, partition int32, key, value []byte) error {
	messageSet := getMessageSet(key, value)
	return (&proto.Payload{
		Topic:        topic,
		Partition:    partition,
		MessageSet:   messageSet,
		RequiredAcks: p.config.RequiredAcks,
		AckTimeout:   p.config.AckTimeout,
	}).Produce(p.cluster)
}

func getMessageSet(key, value []byte) []proto.OffsetMessage {
	return []proto.OffsetMessage{
		{
			SizedMessage: proto.SizedMessage{CRCMessage: proto.CRCMessage{
				Message: proto.Message{
					Key:   key,
					Value: value,
				},
			}}},
	}
}
