package client

import (
	"h12.me/kafka/broker"
)

type brokerPool struct {
	m map[string]*broker.B
}

func newBrokerPool() *brokerPool {
	return &brokerPool{
		m: make(map[string]*broker.B),
	}
}
