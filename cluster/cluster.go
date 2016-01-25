package cluster

import (
	"errors"
	"fmt"
	"sync"

	"h12.me/kafka/broker"
	"h12.me/kafka/log"
)

var (
	ErrLeaderNotFound = errors.New("leader not found")
	ErrCoordNotFound  = errors.New("coordinator not found")
	ErrNoBrokerFound  = errors.New("no broker found")
)

type Config struct {
	Brokers      []string
	BrokerConfig broker.Config
}

func DefaultConfig(brokers ...string) *Config {
	return &Config{
		Brokers:      brokers,
		BrokerConfig: *broker.DefaultConfig(),
	}
}

type C struct {
	*Config
	topics *topicPartitions
	pool   *brokerPool
	mu     sync.Mutex
}

func New(config *Config) (*C, error) {
	c := &C{
		Config: config,
		topics: newTopicPartitions(),
		pool: newBrokerPool(func(addr string) *broker.B {
			cfg := config.BrokerConfig
			cfg.Connection.Addr = addr
			return broker.New(&cfg)
		}),
	}
	for _, addr := range config.Brokers {
		c.pool.AddAddr(addr)
	}
	return c, nil
}

func (c *C) Partitions(topic string) ([]int32, error) {
	partitions := c.topics.getPartitions(topic)
	if len(partitions) > 0 {
		return partitions, nil
	}
	if err := c.updateFromTopicMetadata(topic); err != nil {
		return nil, err
	}
	partitions = c.topics.getPartitions(topic)
	if len(partitions) > 0 {
		return partitions, nil
	}
	return nil, fmt.Errorf("topic %s not found", topic)
}

func (c *C) Coordinator(topic, group string) (*broker.B, error) {
	if coord, err := c.pool.GetCoordinator(group); err == nil {
		return coord, nil
	}
	if err := c.updateCoordinator(topic, group); err != nil {
		return nil, err
	}
	return c.pool.GetCoordinator(group)
}

func (c *C) Leader(topic string, partition int32) (*broker.B, error) {
	if leader, err := c.pool.GetLeader(topic, partition); err == nil {
		return leader, nil
	}
	if err := c.updateFromTopicMetadata(topic); err != nil {
		return nil, err
	}
	return c.pool.GetLeader(topic, partition)
}

func (c *C) LeaderIsDown(topic string, partition int32) {
	log.Warnf("leader (%s,%d) is down", topic, partition)
	c.pool.DeleteLeader(topic, partition)
}

func (c *C) CoordinatorIsDown(group string) {
	log.Warnf("coordinator (%s) is down", group)
	c.pool.DeleteCoordinator(group)
}

func (c *C) updateCoordinator(topic, group string) error {
	brokers, err := c.pool.Brokers()
	if err != nil {
		return err
	}
	var merr MultiError
	for _, broker := range brokers {
		m, err := broker.GroupCoordinator(group)
		if err != nil {
			merr = append(merr, err)
			continue
		}
		c.pool.SetCoordinator(group, m.Broker.NodeID, m.Broker.Addr())
		return nil
	}
	return merr
}

func (c *C) updateFromTopicMetadata(topic string) error {
	brokers, err := c.pool.Brokers()
	if err != nil {
		return err
	}
	var merr MultiError
	for _, broker := range brokers {
		m, err := broker.TopicMetadata(topic)
		if err != nil {
			merr = append(merr, err)
			continue
		}
		for i := range m.Brokers {
			b := &m.Brokers[i]
			c.pool.Add(b.NodeID, b.Addr())
		}
		for i := range m.TopicMetadatas {
			t := &m.TopicMetadatas[i]
			if t.ErrorCode.HasError() {
				merr = append(merr, t.ErrorCode)
				continue
			}
			if t.TopicName == topic {
				partitions := make([]int32, len(t.PartitionMetadatas))
				if t.ErrorCode.HasError() {
					merr = append(merr, t.ErrorCode)
					continue
				}
				for i := range t.PartitionMetadatas {
					partition := &t.PartitionMetadatas[i]
					if partition.ErrorCode.HasError() {
						merr = append(merr, partition.ErrorCode)
						continue
					}
					partitions[i] = partition.PartitionID
					if err := c.pool.SetLeader(topic, partition.PartitionID, partition.Leader); err != nil {
						merr = append(merr, err)
						continue
					}
				}
				c.topics.addPartitions(topic, partitions)
				return nil
			}
		}
	}
	return merr
}
