package consumer

import (
	"errors"
	"time"

	"h12.me/kafka/cluster"
	"h12.me/kafka/proto"
)

var (
	ErrOffsetNotFound          = errors.New("offset not found")
	ErrFailToFetchOffsetByTime = errors.New("fail to fetch offset by time")
)

type Message struct {
	Key    []byte
	Value  []byte
	Offset int64
}

type Config struct {
	Cluster         cluster.Config
	MaxWaitTime     time.Duration
	MinBytes        int
	MaxBytes        int
	OffsetRetention time.Duration
}

func DefaultConfig(brokers ...string) *Config {
	return &Config{
		Cluster:         *cluster.DefaultConfig(brokers...),
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

func New(config *Config, cl *cluster.C) (*C, error) {
	if cl == nil {
		var err error
		cl, err = cluster.New(&config.Cluster)
		if err != nil {
			return nil, err
		}
	}
	config.Cluster = *cl.Config
	return &C{
		cluster: cl,
		config:  config,
	}, nil
}

func (c *C) SearchOffsetByTime(topic string, partition int32, keyTime time.Time, getTime proto.GetTimeFunc) (int64, error) {
	return (&proto.OffsetByTime{
		Topic:     topic,
		Partition: partition,
		Time:      keyTime,
	}).Search(c.cluster, getTime)
}

func (c *C) Offset(topic string, partition int32, consumerGroup string) (int64, error) {
	return (&proto.Offset{Topic: topic, Partition: partition, Group: consumerGroup}).Fetch(c.cluster)
}

func (c *C) Consume(topic string, partition int32, offset int64) (messages []Message, err error) {
	ms, err := (&proto.Messages{
		Topic:       topic,
		Partition:   partition,
		Offset:      offset,
		MinBytes:    c.config.MinBytes,
		MaxBytes:    c.config.MaxBytes,
		MaxWaitTime: c.config.MaxWaitTime,
	}).Consume(c.cluster)
	if err != nil {
		return nil, err
	}
	for i := range ms {
		m := &ms[i].SizedMessage.CRCMessage.Message
		messages = append(messages, Message{
			Key:    m.Key,
			Value:  m.Value,
			Offset: ms[i].Offset,
		})
	}
	return
}

func (c *C) Commit(topic string, partition int32, consumerGroup string, offset int64) error {
	return (&proto.Offset{
		Topic:     topic,
		Partition: partition,
		Group:     consumerGroup,
		Offset:    offset,
	}).Commit(c.cluster)
}
