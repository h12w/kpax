package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"h12.me/kafka/broker"
)

const (
	clientID = "h12.me/kafka/kafpro"
)

type Config struct {
	Broker  string
	Meta    MetaConfig
	Coord   CoordConfig
	Offset  OffsetConfig
	Commit  CommitConfig
	Time    TimeConfig
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

type TimeConfig struct {
	Topic     string
	Partition int
	Time      string
}

type ConsumeConfig struct {
	Topic     string
	Partition int
	Offset    int
}

type MetaConfig struct {
	Topics []string
}

type CommitConfig struct {
	GroupName string
	Topic     string
	Partition int
	Offset    int
	Retention int // millisecond
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
	case "meta":
		var topicsArg string
		flag.StringVar(&topicsArg, "topics", "", "topic names seperated by comma")
		flag.Parse()
		br := broker.New(broker.DefaultConfig().WithAddr(cfg.Broker))
		cfg.Meta.Topics = strings.Split(topicsArg, ",")
		if err := meta(br, &cfg.Meta); err != nil {
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
	case "commit":
		flag.StringVar(&cfg.Commit.GroupName, "group", "", "group name")
		flag.StringVar(&cfg.Commit.Topic, "topic", "", "topic name")
		flag.IntVar(&cfg.Commit.Partition, "partition", 0, "partition")
		flag.IntVar(&cfg.Commit.Offset, "offset", 0, "offset")
		flag.IntVar(&cfg.Commit.Retention, "retention", 0, "retention")
		flag.Parse()
		br := broker.New(broker.DefaultConfig().WithAddr(cfg.Broker))
		if err := commit(br, &cfg.Commit); err != nil {
			log.Fatal(err)
		}
	case "time":
		flag.StringVar(&cfg.Time.Topic, "topic", "", "topic name")
		flag.IntVar(&cfg.Time.Partition, "partition", 0, "partition")
		flag.StringVar(&cfg.Time.Time, "time", "", "time")
		flag.Parse()
		br := broker.New(broker.DefaultConfig().WithAddr(cfg.Broker))
		if err := timeOffset(br, &cfg.Time); err != nil {
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

	meta     TopicMetadataRequest
	consume  FetchRequest
	time     OffsetRequest
	offset   OffsetFetchRequestV1
	commit   OffsetCommitRequestV1
	coord    GroupCoordinatorRequest

`)
	flag.PrintDefaults()
}

func meta(br *broker.B, cfg *MetaConfig) error {
	resp, err := br.TopicMetadata(cfg.Topics...)
	if err != nil {
		return err
	}
	fmt.Println(toJSON(resp))
	return nil
}

func coord(br *broker.B, coord *CoordConfig) error {
	reqMsg := broker.GroupCoordinatorRequest(coord.GroupName)
	req := &broker.Request{
		ClientID:       clientID,
		RequestMessage: &reqMsg,
	}
	resp := &broker.GroupCoordinatorResponse{}
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
	req := &broker.Request{
		ClientID: clientID,
		RequestMessage: &broker.OffsetFetchRequestV1{
			ConsumerGroup: cfg.GroupName,
			PartitionInTopics: []broker.PartitionInTopic{
				{
					TopicName:  cfg.Topic,
					Partitions: []int32{int32(cfg.Partition)},
				},
			},
		},
	}
	resp := broker.OffsetFetchResponse{}
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

func commit(br *broker.B, cfg *CommitConfig) error {
	req := &broker.Request{
		ClientID: clientID,
		RequestMessage: &broker.OffsetCommitRequestV1{
			ConsumerGroupID: cfg.GroupName,
			OffsetCommitInTopicV1s: []broker.OffsetCommitInTopicV1{
				{
					TopicName: cfg.Topic,
					OffsetCommitInPartitionV1s: []broker.OffsetCommitInPartitionV1{
						{
							Partition: int32(cfg.Partition),
							Offset:    int64(cfg.Offset),
							// TimeStamp in milliseconds
							TimeStamp: time.Now().Add(time.Duration(cfg.Retention)*time.Millisecond).Unix() * 1000,
						},
					},
				},
			},
		},
	}
	resp := broker.OffsetCommitResponse{}
	if err := br.Do(req, &resp); err != nil {
		return err
	}
	for i := range resp {
		t := &resp[i]
		if t.TopicName == cfg.Topic {
			for j := range t.ErrorInPartitions {
				p := &t.ErrorInPartitions[j]
				if int(p.Partition) == cfg.Partition {
					if p.ErrorCode.HasError() {
						return p.ErrorCode
					}
				}
			}
		}
	}
	return nil
}

func timeOffset(br *broker.B, cfg *TimeConfig) error {
	var t time.Time
	switch cfg.Time {
	case "latest":
		t = broker.Latest
	case "earliest":
		t = broker.Earliest
	default:
		var err error
		t, err = time.Parse("2006-01-02T15:04:05", cfg.Time)
		if err != nil {
			return err
		}
	}
	resp, err := br.OffsetByTime(cfg.Topic, int32(cfg.Partition), t)
	if err != nil {
		return err
	}
	fmt.Println(toJSON(&resp))
	return nil
}

func consume(br *broker.B, cfg *ConsumeConfig) error {
	req := &broker.Request{
		ClientID: clientID,
		RequestMessage: &broker.FetchRequest{
			ReplicaID:   -1,
			MaxWaitTime: int32(time.Second / time.Millisecond),
			MinBytes:    int32(1024),
			FetchOffsetInTopics: []broker.FetchOffsetInTopic{
				{
					TopicName: cfg.Topic,
					FetchOffsetInPartitions: []broker.FetchOffsetInPartition{
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
	resp := broker.FetchResponse{}
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
