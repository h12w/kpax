package cluster

import (
	"time"

	"h12.me/kafka/proto"
)

func (c *C) Commit(offset *proto.Offset) error {
	coord, err := c.Coordinator(offset.Topic, offset.Group)
	if err != nil {
		return err
	}
	if err := offset.Commit(coord); err != nil {
		if proto.IsNotCoordinator(err) {
			c.CoordinatorIsDown(offset.Group)
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
	offset, err := (&proto.SegmentOffset{Topic: topic, Partition: partition, Time: t}).Fetch(leader)
	if err != nil {
		if proto.IsNotLeader(err) {
			c.LeaderIsDown(topic, partition)
		}
		return -1, err
	}
	return offset, nil
}
