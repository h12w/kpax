package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"h12.me/kafka/broker"
	"h12.me/kafka/proto"
)

const (
	clientID = "h12.me/kafka/kafpro"
)

type Config struct {
	Broker  string
	Leader  LeaderConfig
	Coord   CoordConfig
	Offset  OffsetConfig
	Consume ConsumeConfig
}

type CoordConfig struct {
	GroupName string
}

type OffsetConfig struct {
	GroupName string
	Topic     string
	Partition int
}

type ConsumeConfig struct {
	Topic     string
	Partition int
	Offset    int
}

type LeaderConfig struct {
	Topic string
}

func main() {
	var cfg Config
	flag.StringVar(&cfg.Broker, "broker", "", "broker address")

	// get subcommand
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	subCmd := os.Args[1]
	os.Args = append(os.Args[0:1], os.Args[2:]...)

	switch subCmd {
	case "leader":
		flag.StringVar(&cfg.Leader.Topic, "topic", "", "topic name")
		flag.Parse()
		br := broker.New(broker.DefaultConfig().WithAddr(cfg.Broker))
		if err := leader(br, &cfg.Leader); err != nil {
			log.Fatal(err)
		}
	case "coord":
		flag.StringVar(&cfg.Coord.GroupName, "group", "", "group name")
		flag.Parse()
		br := broker.New(broker.DefaultConfig().WithAddr(cfg.Broker))
		if err := coord(br, &cfg.Coord); err != nil {
			log.Fatal(err)
		}
	case "offset":
		flag.StringVar(&cfg.Offset.GroupName, "group", "", "group name")
		flag.StringVar(&cfg.Offset.Topic, "topic", "", "topic name")
		flag.IntVar(&cfg.Offset.Partition, "partition", 0, "partition")
		flag.Parse()
		br := broker.New(broker.DefaultConfig().WithAddr(cfg.Broker))
		if err := offset(br, &cfg.Offset); err != nil {
			log.Fatal(err)
		}
	case "consume":
		flag.StringVar(&cfg.Consume.Topic, "topic", "", "topic name")
		flag.IntVar(&cfg.Consume.Partition, "partition", 0, "partition")
		flag.IntVar(&cfg.Consume.Offset, "offset", 0, "offset")
		flag.Parse()
		br := broker.New(broker.DefaultConfig().WithAddr(cfg.Broker))
		if err := consume(br, &cfg.Consume); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("invalid subcommand %s", subCmd)
	}
}

func usage() {
	fmt.Println(`
kafpro is a command line tool for querying Kafka wire API

Usage:

	kafpro command [arguments]

The commands are:

	leader   TopicMetadataRequest
	consume  FetchRequest
	offset   OffsetFetchRequestV1
	coord    GroupCoordinatorRequest

`)
	flag.PrintDefaults()
}

func leader(br *broker.B, cfg *LeaderConfig) error {
	req := &proto.Request{
		ClientID:       clientID,
		RequestMessage: &proto.TopicMetadataRequest{cfg.Topic},
	}
	resp := &proto.TopicMetadataResponse{}
	if err := br.Do(req, resp); err != nil {
		return err
	}
	fmt.Println(toJSON(&resp))
	return nil
}

func coord(br *broker.B, coord *CoordConfig) error {
	reqMsg := proto.GroupCoordinatorRequest(coord.GroupName)
	req := &proto.Request{
		ClientID:       clientID,
		RequestMessage: &reqMsg,
	}
	resp := &proto.GroupCoordinatorResponse{}
	if err := br.Do(req, resp); err != nil {
		return err
	}
	fmt.Println(toJSON(&resp))
	if resp.ErrorCode.HasError() {
		return resp.ErrorCode
	}
	return nil
}

func offset(br *broker.B, cfg *OffsetConfig) error {
	req := &proto.Request{
		ClientID: clientID,
		RequestMessage: &proto.OffsetFetchRequestV1{
			ConsumerGroup: cfg.GroupName,
			PartitionInTopics: []proto.PartitionInTopic{
				{
					TopicName:  cfg.Topic,
					Partitions: []int32{int32(cfg.Partition)},
				},
			},
		},
	}
	resp := proto.OffsetFetchResponse{}
	if err := br.Do(req, &resp); err != nil {
		return err
	}
	fmt.Println(toJSON(&resp))
	for i := range resp {
		t := &resp[i]
		if t.TopicName == cfg.Topic {
			for j := range resp[i].OffsetMetadataInPartitions {
				p := &t.OffsetMetadataInPartitions[j]
				if p.Partition == int32(cfg.Partition) {
					if p.ErrorCode.HasError() {
						return p.ErrorCode
					}
				}
			}
		}
	}
	return nil
}

func consume(br *broker.B, cfg *ConsumeConfig) error {
	req := &proto.Request{
		ClientID: clientID,
		RequestMessage: &proto.FetchRequest{
			ReplicaID:   -1,
			MaxWaitTime: int32(time.Second / time.Millisecond),
			MinBytes:    int32(1024),
			FetchOffsetInTopics: []proto.FetchOffsetInTopic{
				{
					TopicName: cfg.Topic,
					FetchOffsetInPartitions: []proto.FetchOffsetInPartition{
						{
							Partition:   int32(cfg.Partition),
							FetchOffset: int64(cfg.Offset),
							MaxBytes:    int32(1000),
						},
					},
				},
			},
		},
	}
	resp := proto.FetchResponse{}
	if err := br.Do(req, &resp); err != nil {
		return err
	}
	fmt.Println(toJSON(resp))
	return nil
}

func toJSON(v interface{}) string {
	buf, _ := json.MarshalIndent(v, "", "\t")
	return string(buf)
}
