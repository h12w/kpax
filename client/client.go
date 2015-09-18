package client

import (
	"errors"
	"sync"

	"h12.me/kafka/broker"
	"h12.me/kafka/proto"
)

var (
	ErrTopicNotFound  = errors.New("topic not found")
	ErrLeaderNotFound = errors.New("leader not found")
	ErrCoordNotFound  = errors.New("coordinator not found")
	ErrNoBrokerFound  = errors.New("no broker found")
)

type Config struct {
	Brokers      []string
	BrokerConfig broker.Config
	ClientID     string
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
		APIKey:         req.APIKey(),
		APIVersion:     req.APIVersion(),
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
	return nil, ErrTopicNotFound
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

func (c *C) updateFromConsumerMetadata(topic, consumerGroup string) error {
	m, err := c.getConsumerMetadata(consumerGroup)
	if err != nil {
		return err
	}
	if m.ErrorCode != 0 {
		return ErrCoordNotFound
	}
	c.pool.SetCoordinator(consumerGroup, m.CoordinatorID, m.CoordinatorHost, m.CoordinatorPort)
	return nil
}

func (c *C) updateFromTopicMetadata(topic string) error {
	m, err := c.getTopicMetadata(topic)
	if err != nil {
		return err
	}
	for i := range m.Brokers {
		b := &m.Brokers[i]
		c.pool.Add(b.NodeID, b.Host, b.Port)
	}
	for i := range m.TopicMetadatas {
		t := &m.TopicMetadatas[i]
		if t.TopicName == topic {
			partitions := make([]int32, len(t.PartitionMetadatas))
			for i := range t.PartitionMetadatas {
				partition := &t.PartitionMetadatas[i]
				partitions[i] = partition.PartitionID
				if err := c.pool.SetLeader(topic, partition.PartitionID, partition.Leader); err != nil {
					return ErrLeaderNotFound
				}
			}
			c.topics.addPartitions(topic, partitions)
			return nil
		}
	}
	return ErrTopicNotFound
}

func (c *C) getTopicMetadata(topic string) (*proto.TopicMetadataResponse, error) {
	var err error
	req := c.NewRequest(&proto.TopicMetadataRequest{topic})
	resp := &proto.TopicMetadataResponse{}
	brokers, err := c.pool.Brokers()
	if err != nil {
		return nil, err
	}
	for _, broker := range brokers {
		err = broker.Do(req, resp)
		if err == nil {
			return resp, nil
		}
	}
	return nil, err
}

func (c *C) getConsumerMetadata(consumerGroup string) (*proto.ConsumerMetadataResponse, error) {
	var err error
	creq := proto.ConsumerMetadataRequest(consumerGroup)
	req := c.NewRequest(&creq)
	resp := &proto.ConsumerMetadataResponse{}
	brokers, err := c.pool.Brokers()
	if err != nil {
		return nil, err
	}
	for _, broker := range brokers {
		err = broker.Do(req, resp)
		if err == nil {
			return resp, nil
		}
	}
	return nil, err
}
