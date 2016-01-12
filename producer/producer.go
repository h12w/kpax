package producer

import (
	"errors"
	"math/rand"
	"time"

	"h12.me/kafka/broker"
	"h12.me/kafka/cluster"
	"h12.me/kafka/log"
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
}

func DefaultConfig(brokers ...string) *Config {
	return &Config{
		Cluster:            *cluster.DefaultConfig(brokers...),
		LeaderRecoveryTime: 60 * time.Second,
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
			break
		}

		leader, err := p.cluster.Leader(topic, partition)
		if err != nil {
			partitioner.Skip(partition)
			continue
		}
		resp, err := leader.Produce(topic, partition, messageSet)
		if err != nil {
			if broker.IsNotLeader(err) {
				p.cluster.LeaderIsDown(topic, partition)
				continue nextPartition
			}
			return err
		}
		for i := range *resp {
			t := &(*resp)[i]
			if t.TopicName != topic {
				continue
			}
			for j := range t.OffsetInPartitions {
				pres := &t.OffsetInPartitions[j]
				if pres.Partition != partition {
					continue
				}
				if pres.ErrorCode.HasError() {
					if broker.IsNotLeader(err) {
						p.cluster.LeaderIsDown(topic, partition)
						continue nextPartition
					}
					return pres.ErrorCode
				}
			}
		}
		return nil
	}
	log.Warnf("fail to find a usable partition")
	return ErrProduceFailed
}

func getMessageSet(key, value []byte) []broker.OffsetMessage {
	return []broker.OffsetMessage{
		{
			SizedMessage: broker.SizedMessage{CRCMessage: broker.CRCMessage{
				Message: broker.Message{
					Key:   key,
					Value: value,
				},
			}}},
	}
}

func (p *P) ProduceWithPartition(topic string, partition int32, key, value []byte) error {
	leader, err := p.cluster.Leader(topic, partition)
	if err != nil {
		return err
	}
	messageSet := getMessageSet(key, value)
	resp, err := leader.Produce(topic, partition, messageSet)
	if err != nil {
		if broker.IsNotLeader(err) {
			p.cluster.LeaderIsDown(topic, partition)
		}
		return err
	}
	for i := range *resp {
		t := &(*resp)[i]
		if t.TopicName != topic {
			continue
		}
		for j := range t.OffsetInPartitions {
			p := &t.OffsetInPartitions[j]
			if p.Partition != partition {
				continue
			}
			if p.ErrorCode.HasError() {
				return p.ErrorCode
			}
		}
	}
	return nil
}
