package consumer

import (
	"errors"
	"fmt"
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

type GetTimeFunc func([]byte) (time.Time, error)

func (c *C) getTime(topic string, partition int32, offset int64, getTime GetTimeFunc) (time.Time, error) {
	const maxMessageSize = 10000
	messages, err := c.consumeBytes(topic, partition, offset, maxMessageSize)
	if err != nil {
		return time.Time{}, err
	}
	if len(messages) == 0 {
		return time.Time{}, fmt.Errorf("no messages at topic %s, partition %d, offset %d", topic, partition, offset)
	}
	if messages[0].Offset != offset {
		return time.Time{}, fmt.Errorf("OFFSET MISMATCH!!! %d, %d", messages[0].Offset, offset)
	}
	return getTime(messages[0].Value)
}

func (c *C) SearchOffsetByTime(topic string, partition int32, keyTime time.Time, getTime GetTimeFunc) (int64, error) {
	earliest, err := (&proto.OffsetByTime{
		Topic:     topic,
		Partition: partition,
		Time:      proto.Earliest,
	}).Fetch(c.cluster)
	if err != nil {
		return -1, err
	}
	if keyTime == proto.Earliest {
		return earliest, nil
	}
	latest, err := (&proto.OffsetByTime{
		Topic:     topic,
		Partition: partition,
		Time:      proto.Latest,
	}).Fetch(c.cluster)
	if err != nil {
		return -1, err
	}

	min, max := earliest, latest
	mid := min + (max-min)/2
	for min <= max {
		mid = min + (max-min)/2
		midTime, err := c.getTime(topic, partition, mid, getTime)
		if err != nil {
			return -1, err
		}
		if midTime.Equal(keyTime) {
			break
		} else if midTime.Before(keyTime) {
			min = mid + 1
		} else {
			max = mid - 1
		}
	}
	return c.searchOffsetBefore(topic, partition, earliest, mid, latest, keyTime, getTime)
}

func (c *C) searchOffsetBefore(topic string, partition int32, min, mid, max int64, keyTime time.Time, getTime func([]byte) (time.Time, error)) (int64, error) {
	const maxJitter = 1000 // time may be interleaved in a small range
	midTime, err := c.getTime(topic, partition, mid, getTime)
	if err != nil {
		return -1, err
	}
	if midTime.Before(keyTime) {
		mid -= maxJitter
	} else {
		mid -= 2 * maxJitter // we have passed the point, go back with double jitter
	}
	if mid < min {
		mid = min
	}
	for offset := mid; offset <= max; offset++ {
		messages, err := c.Consume(topic, partition, offset)
		if err != nil {
			return -1, err
		}
		if len(messages) == 0 {
			return -1, fmt.Errorf("fail to search offset: zero message count")
		}
		for _, message := range messages {
			t, err := getTime(message.Value)
			if err != nil {
				return -1, err
			}
			if t.After(keyTime) || t.Equal(keyTime) {
				return offset, nil
			}
			offset = message.Offset
		}
	}
	return mid, nil
}

func (c *C) Offset(topic string, partition int32, consumerGroup string) (int64, error) {
	return (&proto.Offset{Topic: topic, Partition: partition, Group: consumerGroup}).Fetch(c.cluster)
}

func (c *C) Consume(topic string, partition int32, offset int64) (messages []Message, err error) {
	return c.consumeBytes(topic, partition, offset, c.config.MaxBytes)
}

func (c *C) consumeBytes(topic string, partition int32, offset int64, maxBytes int) (messages []Message, err error) {
	ms, err := (&proto.Messages{
		Topic:       topic,
		Partition:   partition,
		Offset:      offset,
		MinBytes:    c.config.MinBytes,
		MaxBytes:    maxBytes,
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
