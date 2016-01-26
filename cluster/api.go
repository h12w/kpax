package cluster

import (
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
