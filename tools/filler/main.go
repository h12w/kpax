package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"h12.me/kafka/producer"
)

func main() {
	var (
		brokerList string
		topic      string
		count      int
	)
	flag.StringVar(&brokerList, "brokers", "", "broker address")
	flag.StringVar(&topic, "topic", "", "topic")
	flag.IntVar(&count, "count", 1, "message count to send")
	flag.Parse()
	brokers := strings.Split(brokerList, ",")
	pr, err := producer.New(producer.DefaultConfig(brokers...))
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < count; i++ {
		if err := pr.Produce(topic, nil, []byte(strconv.Itoa(i))); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("%d messages sent to %s\n", count, topic)
}
