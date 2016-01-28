package cluster

import (
	"errors"
	"fmt"
	"sync"

	"h12.me/kafka/common"
	"h12.me/kafka/log"
	"h12.me/kafka/proto"
)

var (
	ErrLeaderNotFound = errors.New("leader not found")
	ErrCoordNotFound  = errors.New("coordinator not found")
	ErrNoBrokerFound  = errors.New("no broker found")
)

type (
	C struct {
		topics *topicPartitions
		pool   *brokerPool
		mu     sync.Mutex
	}
	NewBrokerFunc func(addr string) common.Broker
)

func New(newBroker NewBrokerFunc, brokers []string) (*C, error) {
	c := &C{
		topics: newTopicPartitions(),
		pool:   newBrokerPool(newBroker),
	}
	for _, addr := range brokers {
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

func (c *C) Coordinator(group string) (common.Broker, error) {
	if coord, err := c.pool.GetCoordinator(group); err == nil {
		return coord, nil
	}
	if err := c.updateCoordinator(group); err != nil {
		return nil, err
	}
	return c.pool.GetCoordinator(group)
}

func (c *C) Leader(topic string, partition int32) (common.Broker, error) {
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

func (c *C) updateCoordinator(group string) error {
	brokers, err := c.pool.Brokers()
	if err != nil {
		return err
	}
	var merr MultiError
	for _, broker := range brokers {
		coord, err := proto.GroupCoordinator(group).Fetch(broker)
		if err != nil {
			merr = append(merr, err)
			continue
		}
		c.pool.SetCoordinator(group, coord.NodeID, coord.Addr())
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
		m, err := proto.Metadata{topic}.Fetch(broker)
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
			if t.HasError() {
				merr = append(merr, t.ErrorCode)
				continue
			}
			if t.TopicName == topic {
				partitions := make([]int32, len(t.PartitionMetadatas))
				if t.HasError() {
					merr = append(merr, t.ErrorCode)
					continue
				}
				for i := range t.PartitionMetadatas {
					partition := &t.PartitionMetadatas[i]
					if partition.HasError() {
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
