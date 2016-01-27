package consumer

import (
	"errors"
	"fmt"
	"time"

	"h12.me/kafka/broker"
	"h12.me/kafka/cluster"
	"h12.me/kafka/log"
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
	earliest, err := c.cluster.SegmentOffset(topic, partition, broker.Earliest)
	if err != nil {
		return -1, err
	}
	if keyTime == broker.Earliest {
		return earliest, nil
	}
	latest, err := c.cluster.SegmentOffset(topic, partition, broker.Latest)
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
	req := broker.OffsetFetchRequestV1{
		ConsumerGroup: consumerGroup,
		PartitionInTopics: []broker.PartitionInTopic{
			{
				TopicName:  topic,
				Partitions: []int32{partition},
			},
		},
	}
	resp := broker.OffsetFetchResponse{}
	coord, err := c.cluster.Coordinator(topic, consumerGroup)
	if err != nil {
		log.Debugf("fail to get coordinator %v", err)
		return 0, err
	}
	if err := coord.Do(&req, &resp); err != nil {
		if broker.IsNotCoordinator(err) {
			c.cluster.CoordinatorIsDown(consumerGroup)
		}
		log.Debugf("fail to get offset %v", err)
		return 0, err
	}
	for i := range resp {
		t := &resp[i]
		if t.TopicName == topic {
			for j := range resp[i].OffsetMetadataInPartitions {
				p := &t.OffsetMetadataInPartitions[j]
				if p.ErrorCode.HasError() {
					if broker.IsNotCoordinator(err) {
						c.cluster.CoordinatorIsDown(consumerGroup)
					}
					log.Debugf("fail to get offset %v", p.ErrorCode)
					return 0, fmt.Errorf("fail to get offset for (%s, %d): %v", topic, p.Partition, p.ErrorCode)
				}
				return p.Offset, nil
			}
		}
	}
	log.Debugf("fail to get offset %v", ErrOffsetNotFound)
	return 0, ErrOffsetNotFound
}

func (c *C) Consume(topic string, partition int32, offset int64) (messages []Message, err error) {
	return c.consumeBytes(topic, partition, offset, c.config.MaxBytes)
}

func (c *C) consumeBytes(topic string, partition int32, offset int64, maxBytes int) (messages []Message, err error) {
	req := broker.FetchRequest{
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
						MaxBytes:    int32(maxBytes),
					},
				},
			},
		},
	}
	resp := broker.FetchResponse{}
	leader, err := c.cluster.Leader(topic, partition)
	if err != nil {
		log.Debugf("fail to get leader %v", err)
		return nil, err
	}
	if err := leader.Do(&req, &resp); err != nil {
		if broker.IsNotLeader(err) {
			c.cluster.LeaderIsDown(topic, partition)
		}
		log.Debugf("fail to consume %v", err)
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
			if p.ErrorCode.HasError() {
				if broker.IsNotLeader(p.ErrorCode) {
					c.cluster.LeaderIsDown(topic, partition)
				}
				log.Debugf("fail to consume %v", p.ErrorCode)
				return nil, p.ErrorCode
			}
			ms := p.MessageSet
			ms, err := ms.Flatten()
			if err != nil {
				return nil, err
			}
			for k := range ms {
				m := &ms[k]
				if m.Offset == offset {
					ms = ms[k:]
					break
				}
			}
			if len(ms) == 0 {
				continue
			}
			if ms[0].Offset != offset {
				log.Debugf("offset mismatch")
				return nil, fmt.Errorf("2: OFFSET MISMATCH %d %d", ms[0].Offset, offset)
			}
			for i := range ms {
				m := &ms[i].SizedMessage.CRCMessage.Message
				messages = append(messages, Message{
					Key:    m.Key,
					Value:  m.Value,
					Offset: ms[i].Offset,
				})
			}
		}
	}
	return
}

func (c *C) Commit(topic string, partition int32, consumerGroup string, offset int64) error {
	return c.cluster.Commit(&broker.OffsetCommit{
		Topic:     topic,
		Partition: partition,
		Group:     consumerGroup,
		Offset:    offset,
	})
}
