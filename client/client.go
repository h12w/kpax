package client

import (
	"errors"
	"fmt"
	"sync"

	"h12.me/kafka/broker"
	"h12.me/kafka/log"
	"h12.me/kafka/proto"
)

var (
	ErrLeaderNotFound = errors.New("leader not found")
	ErrCoordNotFound  = errors.New("coordinator not found")
	ErrNoBrokerFound  = errors.New("no broker found")
)

type Config struct {
	Brokers      []string
	BrokerConfig broker.Config
	ClientID     string
}

func DefaultConfig(brokers ...string) *Config {
	return &Config{
		Brokers:      brokers,
		BrokerConfig: *broker.DefaultConfig(),
		ClientID:     "h12.me/kafka",
	}
}

type C struct {
	config *Config
	topics *topicPartitions
	pool   *brokerPool
	mu     sync.Mutex
}

func New(config *Config) (*C, error) {
	c := &C{
		config: config,
		topics: newTopicPartitions(),
		pool: newBrokerPool(func(addr string) *broker.B {
			cfg := config.BrokerConfig
			cfg.Addr = addr
			return broker.New(&cfg)
		}),
	}
	for _, addr := range config.Brokers {
		c.pool.AddAddr(addr)
	}
	return c, nil
}

func (c *C) NewRequest(req proto.RequestMessage) *proto.Request {
	return &proto.Request{
		ClientID:       c.config.ClientID,
		RequestMessage: req,
	}
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

func (c *C) Coordinator(topic, consumerGroup string) (*broker.B, error) {
	if coord, err := c.pool.GetCoordinator(consumerGroup); err == nil {
		return coord, nil
	}
	if err := c.updateFromConsumerMetadata(topic, consumerGroup); err != nil {
		return nil, err
	}
	return c.pool.GetCoordinator(consumerGroup)
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

func (c *C) CoordinatorIsDown(consumerGroup string) {
	log.Warnf("coordinator (%s) is down", consumerGroup)
	c.pool.DeleteCoordinator(consumerGroup)
}

func (c *C) updateFromConsumerMetadata(topic, consumerGroup string) error {
	brokers, err := c.pool.Brokers()
	if err != nil {
		return err
	}
	for _, broker := range brokers {
		for i := 0; i < 2; i++ { // twice
			m, err := c.getGroupCoordinator(broker, consumerGroup)
			if err != nil {
				return err
			}
			if m.ErrorCode.HasError() {
				return ErrCoordNotFound
			}
			c.pool.SetCoordinator(consumerGroup, m.Broker.NodeID, m.Broker.Addr())
			return nil
		}
	}
	return nil
}

func (c *C) updateFromTopicMetadata(topic string) error {
	brokers, err := c.pool.Brokers()
	if err != nil {
		return err
	}
	for _, broker := range brokers {
		m, err := c.getTopicMetadata(broker, topic)
		if err != nil {
			return err
		}
		for i := range m.Brokers {
			b := &m.Brokers[i]
			c.pool.Add(b.NodeID, b.Addr())
		}
		for i := range m.TopicMetadatas {
			t := &m.TopicMetadatas[i]
			if t.ErrorCode.HasError() {
				log.Warnf("topic error %v", t.ErrorCode)
			}
			if t.TopicName == topic {
				partitions := make([]int32, len(t.PartitionMetadatas))
				if t.ErrorCode.HasError() {
					return t.ErrorCode
				}
				for i := range t.PartitionMetadatas {
					partition := &t.PartitionMetadatas[i]
					if partition.ErrorCode.HasError() {
						return partition.ErrorCode
					}
					partitions[i] = partition.PartitionID
					if err := c.pool.SetLeader(topic, partition.PartitionID, partition.Leader); err != nil {
						log.Debugf("cannot set leader %s, %d, %d", topic, partition.PartitionID, partition.Leader)
						return ErrLeaderNotFound
					}
				}
				c.topics.addPartitions(topic, partitions)
				return nil
			}
		}
	}
	return fmt.Errorf("topic %s not found", topic)
}

func (c *C) getTopicMetadata(broker *broker.B, topic string) (*proto.TopicMetadataResponse, error) {
	req := c.NewRequest(&proto.TopicMetadataRequest{topic})
	resp := &proto.TopicMetadataResponse{}
	if err := broker.Do(req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *C) getGroupCoordinator(broker *broker.B, consumerGroup string) (*proto.GroupCoordinatorResponse, error) {
	creq := proto.GroupCoordinatorRequest(consumerGroup)
	req := c.NewRequest(&creq)
	resp := &proto.GroupCoordinatorResponse{}
	if err := broker.Do(req, resp); err != nil {
		return nil, err
	}
	if resp.ErrorCode.HasError() {
		return nil, resp.ErrorCode
	}
	return resp, nil
}
