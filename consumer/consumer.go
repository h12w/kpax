package consumer

import (
	"errors"
	"time"

	"h12.me/kpax/model"
	"h12.me/kpax/proto"
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

type C struct {
	MaxWaitTime     time.Duration
	MinBytes        int
	MaxBytes        int
	OffsetRetention time.Duration
	Cluster         model.Cluster
}

func New(cluster model.Cluster) *C {
	return &C{
		Cluster:         cluster,
		MaxWaitTime:     100 * time.Millisecond,
		MinBytes:        1,
		MaxBytes:        1024 * 1024,
		OffsetRetention: 7 * 24 * time.Hour,
	}
}

func (c *C) FetchOffsetByTime(topic string, partition int32, keyTime time.Time) (int64, error) {
	return (&proto.OffsetByTime{
		Topic:     topic,
		Partition: partition,
		Time:      keyTime,
	}).Fetch(c.Cluster)
}

func (c *C) SearchOffsetByTime(topic string, partition int32, keyTime time.Time, getTime proto.GetTimeFunc) (int64, error) {
	return (&proto.OffsetByTime{
		Topic:     topic,
		Partition: partition,
		Time:      keyTime,
	}).Search(c.Cluster, getTime)
}

func (c *C) Offset(topic string, partition int32, consumerGroup string) (int64, error) {
	return (&proto.Offset{Topic: topic, Partition: partition, Group: consumerGroup}).Fetch(c.Cluster)
}

func (c *C) Consume(topic string, partition int32, offset int64) (messages []Message, err error) {
	ms, err := (&proto.Messages{
		Topic:       topic,
		Partition:   partition,
		Offset:      offset,
		MinBytes:    c.MinBytes,
		MaxBytes:    c.MaxBytes,
		MaxWaitTime: c.MaxWaitTime,
	}).Consume(c.Cluster)
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
	}).Commit(c.Cluster)
}
