package client

import (
	"strconv"
	"sync"

	"h12.me/kafka/broker"
	"h12.me/kafka/log"
)

type brokerPool struct {
	addrBroker           map[string]*broker.B
	idAddr               map[int32]string
	topicPartitionLeader map[topicPartition]*broker.B
	groupCoordinator     map[string]*broker.B
	newBroker            func(string) *broker.B
	mu                   sync.Mutex
}

type topicPartition struct {
	topic     string
	partition int32
}

func newBrokerPool(newBroker func(string) *broker.B) *brokerPool {
	return &brokerPool{
		addrBroker:           make(map[string]*broker.B),
		idAddr:               make(map[int32]string),
		topicPartitionLeader: make(map[topicPartition]*broker.B),
		groupCoordinator:     make(map[string]*broker.B),
		newBroker:            newBroker,
	}
}

func (p *brokerPool) Brokers() (map[string]*broker.B, error) {
	if len(p.addrBroker) > 0 {
		return p.addrBroker, nil
	}
	return nil, ErrNoBrokerFound
}

func (p *brokerPool) AddAddr(addr string) *broker.B {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.addAddr(addr)
}

func (p *brokerPool) addAddr(addr string) *broker.B {
	if broker, ok := p.addrBroker[addr]; ok {
		return broker
	}
	broker := p.newBroker(addr)
	p.addrBroker[addr] = broker
	log.Debugf("broker %s added to pool", addr)
	return broker
}

func (p *brokerPool) add(brokerID int32, host string, port int32) *broker.B {
	addr := host + ":" + strconv.Itoa(int(port))
	p.idAddr[brokerID] = addr
	return p.addAddr(addr)
}

func (p *brokerPool) Add(brokerID int32, host string, port int32) *broker.B {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.add(brokerID, host, port)
}

func (p *brokerPool) SetLeader(topic string, partition int32, brokerID int32) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	broker, err := p.find(brokerID)
	if err != nil {
		return err
	}
	p.topicPartitionLeader[topicPartition{topic, partition}] = broker
	return nil
}

func (p *brokerPool) GetLeader(topic string, partition int32) (*broker.B, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	broker, ok := p.topicPartitionLeader[topicPartition{topic, partition}]
	if !ok {
		return nil, ErrLeaderNotFound
	}
	return broker, nil
}

func (p *brokerPool) DeleteLeader(topic string, partition int32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.topicPartitionLeader, topicPartition{topic, partition})
}

func (p *brokerPool) DeleteCoordinator(consumerGroup string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.groupCoordinator, consumerGroup)
}

func (p *brokerPool) SetCoordinator(consumerGroup string, brokerID int32, host string, port int32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	broker := p.add(brokerID, host, port)
	p.groupCoordinator[consumerGroup] = broker
}

func (p *brokerPool) GetCoordinator(consumerGroup string) (*broker.B, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	broker, ok := p.groupCoordinator[consumerGroup]
	if !ok {
		log.Debugf("cannot find coordinator for group %s", consumerGroup)
		return nil, ErrCoordNotFound
	}
	return broker, nil
}

func (p *brokerPool) find(brokerID int32) (*broker.B, error) {
	if addr, ok := p.idAddr[brokerID]; ok {
		if broker, ok := p.addrBroker[addr]; ok {
			return broker, nil
		}
	}
	log.Debugf("cannot find broker %d", brokerID)
	return nil, ErrNoBrokerFound
}

type topicPartitions struct {
	m  map[string][]int32
	mu sync.Mutex
}

func newTopicPartitions() *topicPartitions {
	return &topicPartitions{
		m: make(map[string][]int32),
	}
}

func (tp *topicPartitions) getPartitions(topic string) []int32 {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	return tp.m[topic]
}

func (tp *topicPartitions) addPartitions(topic string, partitions []int32) {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	tp.m[topic] = partitions
}
