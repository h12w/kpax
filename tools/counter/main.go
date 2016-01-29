package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"h12.me/kpax/broker"
	"h12.me/kpax/consumer"
)

func main() {
	var (
		brokerList     string
		topic          string
		partitionRange string
	)
	flag.StringVar(&brokerList, "brokers", "", "broker address")
	flag.StringVar(&topic, "topic", "", "topic")
	flag.StringVar(&partitionRange, "partitions", "", "partition range")
	flag.Parse()
	brokers := strings.Split(brokerList, ",")
	partitions := parsePartitionRange(partitionRange)
	consumerConfig := consumer.DefaultConfig(brokers...)
	cr, err := consumer.New(consumerConfig)
	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup
	counter := counter{
		topic:    topic,
		consumer: cr,
	}
	for partition := int32(partitions[0]); partition <= int32(partitions[1]); partition++ {
		wg.Add(1)
		go func(partition int32) {
			err := counter.consume(partition)
			if err != nil {
				log.Fatal(err)
			}
		}(partition)
	}
	go func() {
		lastCount := int64(0)
		for {
			count := counter.count
			if lastCount != count {
				fmt.Println(count)
				lastCount = count
			}
			time.Sleep(time.Second)
		}
	}()
	wg.Wait()
}

type counter struct {
	count    int64
	topic    string
	consumer *consumer.C
}

func (c *counter) consume(partition int32) error {
	offset, err := c.consumer.OffsetByTime(c.topic, partition, broker.Earliest)
	if err != nil {
		return err
	}
	for {
		messages, err := c.consumer.Consume(c.topic, partition, offset)
		if err != nil {
			return err
		}
		if len(messages) > 0 {
			offset = messages[len(messages)-1].Offset + 1
			c.add(len(messages))
		}
	}
}

func (c *counter) add(delta int) {
	atomic.AddInt64(&c.count, int64(delta))
}

func parsePartitionRange(s string) [2]int {
	ss := strings.Split(s, "-")
	if len(ss) < 2 {
		log.Fatal("wrong range " + s)
	}
	from, err := strconv.Atoi(ss[0])
	if err != nil {
		log.Fatal(err)
	}
	to, err := strconv.Atoi(ss[1])
	if err != nil {
		log.Fatal(err)
	}
	return [2]int{from, to}
}
