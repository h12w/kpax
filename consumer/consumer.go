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
	ErrFailCommitOffset        = errors.New("fail to commit offset")
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

func (c *C) getTime(topic string, partition int32, offset int64, getTime func([]byte) (time.Time, error)) (time.Time, error) {
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

func (c *C) SearchOffsetByTime(topic string, partition int32, keyTime time.Time, getTime func([]byte) (time.Time, error)) (int64, error) {
	earliest, err := c.OffsetByTime(topic, partition, broker.Earliest)
	if err != nil {
		return -1, err
	}
	latest, err := c.OffsetByTime(topic, partition, broker.Latest)
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
			return c.searchOffsetBefore(topic, partition, earliest, mid, latest, keyTime, getTime)
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

func (c *C) OffsetByTime(topic string, partition int32, t time.Time) (int64, error) {
	leader, err := c.cluster.Leader(topic, partition)
	if err != nil {
		return -1, err
	}
	resp, err := leader.OffsetByTime(topic, partition, t)
	if err != nil {
		return -1, err
	}
	for _, t := range *resp {
		if t.TopicName != topic {
			continue
		}
		for _, p := range t.OffsetsInPartitions {
			if p.Partition != partition {
				continue
			}
			if p.ErrorCode.HasError() {
				return -1, p.ErrorCode
			}
			if len(p.Offsets) == 0 {
				return -1, ErrFailToFetchOffsetByTime
			}
			return p.Offsets[0], nil
		}
	}
	return -1, ErrFailToFetchOffsetByTime
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
		log.Debugf("fail to get coordinator %v", err)
		return 0, err
	}
	if err := coord.Do(req, &resp); err != nil {
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
							MaxBytes:    int32(maxBytes),
						},
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
	if err := leader.Do(req, &resp); err != nil {
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
		log.Debugf("fail to get coordinator %v")
		return err
	}
	if err := coord.Do(req, &resp); err != nil {
		if broker.IsNotCoordinator(err) {
			c.cluster.CoordinatorIsDown(consumerGroup)
		}
		log.Debugf("fail to commit %v", err)
		return err
	}
	for i := range resp {
		t := &resp[i]
		if t.TopicName == topic {
			for j := range t.ErrorInPartitions {
				p := &t.ErrorInPartitions[j]
				if p.Partition == partition {
					if p.ErrorCode.HasError() {
						if broker.IsNotCoordinator(err) {
							c.cluster.CoordinatorIsDown(consumerGroup)
						}
						log.Debugf("fail to commit %v", p.ErrorCode)
						return p.ErrorCode
					}
					return nil
				}
			}
		}
	}
	log.Debugf("fail to commit %v", ErrFailCommitOffset)
	return ErrFailCommitOffset
}
