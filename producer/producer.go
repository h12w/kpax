package producer

import (
	"errors"
	"math/rand"
	"time"

	"h12.me/kafka/broker"
	"h12.me/kafka/client"
	"h12.me/kafka/log"
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
	Client             client.Config
	RequiredAcks       int16
	Timeout            time.Duration
	LeaderRecoveryTime time.Duration
}

func DefaultConfig(brokers ...string) *Config {
	return &Config{
		Client:             *client.DefaultConfig(brokers...),
		RequiredAcks:       proto.AckLocal,
		Timeout:            10 * time.Second,
		LeaderRecoveryTime: 60 * time.Second,
	}
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
	messageSet := getMessageSet(key, value)
	for i := 0; i < partitioner.Count(); i++ {
		partition, err := partitioner.Partition(key)
		if err != nil {
			p.topicPartitioner.Delete(topic)
			break
		}

		leader, err := p.client.Leader(topic, partition)
		if err != nil {
			partitioner.Skip(partition)
			continue
		}
		req := p.getProducerRequest(topic, partition, messageSet)

		err = p.doSentMessage(leader, req)
		if err == proto.ErrConn {
			p.client.LeaderIsDown(topic, partition)
		}
		return err
	}
	log.Warnf("fail to find a usable partition")
	return ErrProduceFailed
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

func (p *P) getProducerRequest(topic string, partition int32, messageSet []proto.OffsetMessage) *proto.Request {
	return p.client.NewRequest(&proto.ProduceRequest{
		RequiredAcks: p.config.RequiredAcks,
		Timeout:      int32(p.config.Timeout / time.Millisecond),
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
}

func (p *P) ProduceWithPartition(topic string, partition int32, key, value []byte) error {
	leader, err := p.client.Leader(topic, partition)
	if err != nil {
		return err
	}
	messageSet := getMessageSet(key, value)
	req := p.getProducerRequest(topic, partition, messageSet)
	err = p.doSentMessage(leader, req)
	if err == proto.ErrConn {
		p.client.LeaderIsDown(topic, partition)
	}
	return err
}

func (p *P) doSentMessage(leader *broker.B, req *proto.Request) error {
	resp := proto.ProduceResponse{}
	if err := leader.Do(req, &resp); err != nil {
		return err
	}
	for i := range resp {
		for j := range resp[i].OffsetInPartitions {
			if errCode := resp[i].OffsetInPartitions[j].ErrorCode; errCode != 0 {
				log.Warnf("produce error, code %d", errCode)
				continue
			}
		}
	}
	return nil
}
