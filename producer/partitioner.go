package producer

import (
	"sync"
	"time"

	"h12.me/kafka/log"
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
	partitions   []int32
	skipList     map[int32]time.Time
	recoveryTime time.Duration
	i            int
	mu           sync.Mutex
}

func newPartitioner(partitions []int32) *partitioner {
	return &partitioner{
		partitions: partitions,
		skipList:   make(map[int32]time.Time),
	}
}

func (p *partitioner) Skip(partition int32) {
	p.mu.Lock()
	log.Warnf("partition %d skipped", partition)
	p.skipList[partition] = time.Now().Add(p.recoveryTime)
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
		if expireTime, ok := p.skipList[partition]; ok {
			if time.Now().Before(expireTime) {
				continue
			}
			delete(p.skipList, partition)
		}
		return partition, nil
	}
	return -1, ErrNoValidPartition
}
