package main

import (
	"encoding/json"
	"log"
	"os"
	"runtime"

	"github.com/jessevdk/go-flags"
	"h12.me/config"
	"h12.me/kpax/broker"
	"h12.me/kpax/cluster"
)

const (
	clientID = "h12.me/kpax/kafpro"
)

func parseJSON(cfg interface{}) error {
	var fileConfig struct {
		ConfigFile string `long:"config" default:"config.json"`
	}
	parser := flags.NewParser(&fileConfig, flags.IgnoreUnknown)
	if _, err := parser.Parse(); err != nil {
		return err
	}
	f, err := os.Open(fileConfig.ConfigFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(cfg)
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var cfg Config
	cmd, err := config.ParseCommand(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	c := cluster.New(broker.NewDefault, cfg.Brokers)
	//fmt.Println(toJSON(cfg))
	switch cmd.Name {
	case "consume":
		err = cfg.Consume.Exec(c)
	case "offset":
		err = cfg.Offset.Exec(c)
	case "rollback":
		err = cfg.Rollback.Exec(c)
	case "tail":
		err = cfg.Tail.Exec(c)
	case "produce":
		err = cfg.Produce.Exec(c)
	case "meta":
		err = cfg.Meta.Exec(c)
	default:
		log.Fatal("unkown command " + cmd.Name)
	}
	if err != nil {
		log.Fatal(err)
	}

	/*
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
	*/
}

/*
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
					if p.HasError() {
						return p.ErrorCode
					}
				}
			}
		}
	}
	return nil
}
*/

func toJSON(v interface{}) string {
	buf, _ := json.MarshalIndent(v, "", "\t")
	return string(buf)
}
