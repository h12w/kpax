package producer

import (
	"sync"
)

type topicPartitioner struct {
	m  map[string]*partitioner
	mu sync.Mutex
}

func newTopicPartitioner() *topicPartitioner {
	return &topicPartitioner{
		m: make(map[string]*partitioner),
	}
}

func (tp *topicPartitioner) Get(topic string) *partitioner {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	return tp.m[topic]
}

func (tp *topicPartitioner) Add(topic string, partitions []int32) *partitioner {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	partitioner := tp.m[topic]
	if partitioner == nil {
		partitioner = newPartitioner(partitions)
		tp.m[topic] = partitioner
	}
	return partitioner
}

func (tp *topicPartitioner) Delete(topic string) {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	delete(tp.m, topic)
}

type partitioner struct {
	partitions []int32
	skipList   map[int32]bool
	i          int
	mu         sync.Mutex
}

func newPartitioner(partitions []int32) *partitioner {
	return &partitioner{
		partitions: partitions,
		skipList:   make(map[int32]bool),
	}
}

func (p *partitioner) Skip(partition int32) {
	p.mu.Lock()
	p.skipList[partition] = true
	p.mu.Unlock()
}

func (p *partitioner) Count() int {
	return len(p.partitions)
}

func (p *partitioner) Partition([]byte) (int32, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i := 0; i < len(p.partitions); i++ {
		partition := p.partitions[p.i]
		p.i++
		if p.i == len(p.partitions) {
			p.i = 0
		}
		if !p.skipList[partition] {
			return partition, nil
		}
	}
	return -1, ErrNoValidPartition
}
