package cluster

import (
	"time"

	"h12.me/kafka/broker"
)

func (c *C) Commit(commit *broker.OffsetCommit) error {
	coord, err := c.Coordinator(commit.Topic, commit.Group)
	if err != nil {
		return err
	}
	if err := coord.Commit(commit); err != nil {
		if broker.IsNotCoordinator(err) {
			c.CoordinatorIsDown(commit.Group)
		}
		return err
	}
	return nil
}

func (c *C) SegmentOffset(topic string, partition int32, t time.Time) (int64, error) {
	leader, err := c.Leader(topic, partition)
	if err != nil {
		return -1, err
	}
	return leader.SegmentOffset(topic, partition, t)
}
