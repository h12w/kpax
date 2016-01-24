package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	// TODO: "github.com/jessevdk/go-flags"
	"h12.me/kafka/broker"
)

const (
	clientID = "h12.me/kafka/kafpro"
)

type Config struct {
	Broker  string        `long:"broker"`
	Meta    MetaConfig    `command:"meta"`
	Coord   CoordConfig   `command:"coord"`
	Offset  OffsetConfig  `command:"offset"`
	Commit  CommitConfig  `command:"commit"`
	Time    TimeConfig    `command:"time"`
	Consume ConsumeConfig `command:"consume"`
}

type CoordConfig struct {
	GroupName string `long:"group"`
}

type OffsetConfig struct {
	GroupName string `long:"group"`
	Topic     string `long:"topic"`
	Partition int    `long:"partition"`
}

type TimeConfig struct {
	Topic     string `long:"topic"`
	Partition int    `long:"partition"`
	Time      string `long:"time"`
}

type timestamp time.Time

func (t timestamp) String() string {
	return time.Time(t).Format("2006-01-02T15:04:05")
}

func (t *timestamp) Set(s string) error {
	tm, err := time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		return err
	}
	*t = timestamp(tm)
	return nil
}

type ConsumeConfig struct {
	Topic     string `long:"topic"`
	Partition int    `long:"partition"`
	Offset    int    `long:"offset"`
}

type MetaConfig struct {
	Topics Topics `long:"topics"`
}

type Topics []string

func (ts Topics) String() string {
	return strings.Join(ts, ",")
}

func (ts *Topics) Set(s string) error {
	*ts = strings.Split(s, ",")
	return nil
}

type CommitConfig struct {
	GroupName string `long:"group"`
	Topic     string `long:"topic"`
	Partition int    `long:"partition"`
	Offset    int    `long:"offset"`
	Retention int    `long:"retention"` // millisecond
}

func main() {

	/*
		var opts opts
		parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
		_, err := parser.Parse()
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(parser.Active.Name)
		fmt.Println(opts)
	*/

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
		flag.Var(&cfg.Meta.Topics, "topics", "topic names seperated by comma")
		flag.Parse()
		br := broker.New(broker.DefaultConfig().WithAddr(cfg.Broker))
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
	resp, err := br.GroupCoordinator(coord.GroupName)
	if err != nil {
		return err
	}
	fmt.Println(toJSON(resp))
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
			return fmt.Errorf("error parsing %s: %s", cfg.Time, err.Error())
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
	for _, t := range resp {
		for _, p := range t.FetchMessageSetInPartitions {
			ms, err := p.MessageSet.Flatten()
			if err != nil {
				return err
			}
			fmt.Println(toJSON(ms))
		}
	}
	return nil
}

func toJSON(v interface{}) string {
	buf, _ := json.MarshalIndent(v, "", "\t")
	return string(buf)
}
