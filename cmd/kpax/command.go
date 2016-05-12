package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"h12.me/kpax/consumer"
	"h12.me/kpax/model"
	"h12.me/kpax/producer"
	"h12.me/kpax/proto"
)

const (
	timeFormat = "2006-01-02T15:04:05"
)

type Config struct {
	ConfigFile string  `long:"config" default:"config.json"`
	Brokers    Brokers `long:"brokers" yaml:"brokers"`

	Consume  ConsumeCommand  `command:"consume"  description:"print or count messages wthin a time range"`
	Produce  ProduceCommand  `command:"produce"  description:"produce one message to the topic"`
	Tail     TailCommand     `command:"tail" description:"print latest messages"`
	Offset   OffsetCommand   `command:"offset"   description:"print stored offsets of a topic and group"`
	Rollback RollbackCommand `command:"rollback" description:"commit offsets of a specific time for a topic"`

	Meta  MetaCommand `command:"meta"`
	Coord CoordConfig `command:"coord"`
}

type Brokers []string

func (t *Brokers) UnmarshalFlag(value string) error {
	*t = strings.Split(value, ",")
	return nil
}

type CoordConfig struct {
	Group string `long:"group"`
}

type Timestamp time.Time

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(t).Format("2006-01-02T15:04:05") + `"`), nil
}

func (t *Timestamp) UnmarshalFlag(value string) error {
	switch value {
	case "latest":
		*t = Timestamp(proto.Latest)
	case "earliest":
		*t = Timestamp(proto.Earliest)
	default:
		tm, err := time.Parse(timeFormat, value)
		if err != nil {
			return fmt.Errorf("error parsing %s: %s", value, err.Error())
		}
		*t = Timestamp(tm)
	}
	return nil
}

// parseTime parses time as RFC3339 or unix timestamp
func parseTime(v interface{}) (time.Time, error) {
	if t, ok := v.(time.Time); ok {
		return t, nil
	}
	timeText, ok := v.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("fail to parse time %v", v)
	}
	t, err := time.Parse(time.RFC3339Nano, timeText)
	if err != nil {
		unix, err := strconv.Atoi(timeText)
		if err != nil {
			return time.Time{}, fmt.Errorf("fail to parse %s: %s", timeText, err.Error())
		}
		return time.Unix(int64(unix), 0), nil
	}
	return t, nil
}

type TailCommand struct {
	Topic     string `long:"topic"`
	Partition int    `long:"partition" default:"-1"`
	Format    Format `long:"format"`
	Count     int    `long:"count" short:"n" default:"10"`
}

func (cmd *TailCommand) Exec(cl model.Cluster) error {
	if cmd.Format == "" {
		var err error
		cmd.Format, err = formatDetector{cl, cmd.Topic}.detect()
		if err != nil {
			return err
		}
	}
	var partitions []int32
	if cmd.Partition == -1 {
		var err error
		partitions, err = cl.Partitions(cmd.Topic)
		if err != nil {
			return err
		}
	} else {
		partitions = []int32{int32(cmd.Partition)}
	}
	cr := consumer.New(cl)
	var wg sync.WaitGroup
	wg.Add(len(partitions))
	for _, partition := range partitions {
		start, err := cr.FetchOffsetByTime(cmd.Topic, partition, proto.Earliest)
		if err != nil {
			return err
		}
		end, err := cr.FetchOffsetByTime(cmd.Topic, partition, proto.Latest)
		if err != nil {
			return err
		}
		if start < end-int64(cmd.Count) {
			start = end - int64(cmd.Count)
		}
		go func(partition int32, start, end int64) {
			defer wg.Done()
			offset := start
			for offset < end {
				messages, err := cr.Consume(cmd.Topic, partition, offset)
				if err != nil {
					log.Println(err)
					return
				}
				if len(messages) == 0 {
					break
				}
				for _, msg := range messages {
					line, err := cmd.Format.Sprint(msg.Value)
					if err != nil {
						log.Println(err)
						break
					}
					fmt.Println(line)
				}
				offset = messages[len(messages)-1].Offset + 1
			}
		}(partition, start, end)
	}
	wg.Wait()
	return nil
}

type ConsumeCommand struct {
	Topic     string    `long:"topic"`
	Start     Timestamp `long:"start"`
	End       Timestamp `long:"end"`
	Format    Format    `long:"format" default:"json"`
	TimeField string    `long:"time-field"`
	Count     bool      `long:"count"`
}

func (cmd *ConsumeCommand) Exec(cl model.Cluster) error {
	if cmd.Format == "" {
		var err error
		cmd.Format, err = formatDetector{cl, cmd.Topic}.detect()
		if err != nil {
			return err
		}
	}
	partitions, err := cl.Partitions(cmd.Topic)
	if err != nil {
		return err
	}
	cr := consumer.New(cl)
	var wg sync.WaitGroup
	wg.Add(len(partitions))
	var cnt int64
	timeFunc := cmd.Format.unmarshalTime(cmd.TimeField)
	for _, partition := range partitions {
		go func(partition int32) {
			defer wg.Done()
			partCnt, err := cmd.consumePartition(cr, partition, timeFunc)
			if err != nil {
				log.Println(err)
				return
			}
			atomic.AddInt64(&cnt, partCnt)
		}(partition)
	}
	wg.Wait()
	if cmd.Count {
		fmt.Println(cnt)
	}
	return nil
}

func (cmd *ConsumeCommand) consumePartition(cr *consumer.C, partition int32, timeFunc proto.GetTimeFunc) (int64, error) {
	offset, err := cr.SearchOffsetByTime(cmd.Topic, partition, time.Time(cmd.Start), timeFunc)
	if err != nil {
		return 0, err
	}
	jitterCnt := 0
	cnt := int64(0)
	for jitterCnt <= 1000 {
		messages, err := cr.Consume(cmd.Topic, partition, offset)
		if err != nil {
			return 0, err
		}
		if len(messages) == 0 {
			break
		}
		for _, msg := range messages {
			t, err := timeFunc(msg.Value)
			if err != nil {
				return 0, fmt.Errorf("fail to parse time from %s: %v", string(msg.Value), err)
			}
			if !t.Before(time.Time(cmd.Start)) && t.Before(time.Time(cmd.End)) {
				if !cmd.Count {
					line, err := cmd.Format.Sprint(msg.Value)
					if err != nil {
						log.Println(err)
						break
					}
					fmt.Println(line)
				}
				cnt++
			}
			if t.After(time.Time(cmd.End)) {
				jitterCnt++
			}
		}
		offset = messages[len(messages)-1].Offset + 1
	}
	return cnt, nil
}

type ProduceCommand struct {
	File   string `long:"file"`
	Msg    string `long:"msg"`
	Format Format `long:"format"`
	Topic  string `long:"topic"`
}

func (cmd *ProduceCommand) Exec(cl model.Cluster) error {
	pr := producer.New(cl)
	var msg []byte
	if cmd.Msg != "" {
		msg = []byte(cmd.Msg)
	}
	v := make(map[string]interface{})
	if err := json.Unmarshal(msg, &v); err != nil {
		return err
	}
	fmt.Printf("%#v\n", v)
	value, err := cmd.Format.marshal(v)
	if err != nil {
		return err
	}
	return pr.Produce(cmd.Topic, nil, value)
}

type MetaCommand struct {
	Topic     string `long:"topic"`
	Partition int    `long:"partition"`
}

func (cmd *MetaCommand) Exec(cl model.Cluster) error {
	b, err := cl.Leader(cmd.Topic, int32(cmd.Partition))
	if err != nil {
		return err
	}
	res, err := proto.Metadata(cmd.Topic).Fetch(b)
	if err != nil {
		return err
	}
	fmt.Println(toJSON(res))
	return nil
}

type Topics []string

func (ts Topics) String() string {
	return strings.Join(ts, ",")
}

func (ts *Topics) Set(s string) error {
	*ts = strings.Split(s, ",")
	return nil
}

type OffsetCommand struct {
	Group    string `long:"group"`
	Topic    string `long:"topic"`
	Earliest bool   `long:"earliest"`
	Latest   bool   `long:"latest"`
}

func (cmd *OffsetCommand) Exec(cl model.Cluster) error {
	partitions, err := cl.Partitions(cmd.Topic)
	if err != nil {
		return err
	}
	cr := consumer.New(cl)
	fmt.Printf("topic: %s, group: %s\n", cmd.Topic, cmd.Group)
	for _, partition := range partitions {
		offset := int64(-1)
		var err error
		if cmd.Latest {
			offset, err = cr.FetchOffsetByTime(cmd.Topic, partition, proto.Latest)
		} else if cmd.Earliest {
			offset, err = cr.FetchOffsetByTime(cmd.Topic, partition, proto.Earliest)
		} else {
			offset, err = cr.Offset(cmd.Topic, partition, cmd.Group)
		}
		if err != nil {
			return err
		}
		fmt.Printf("\t%d:%d\n", partition, offset)
	}
	return nil
}

type RollbackCommand struct {
	Group     string    `long:"group"`
	Topic     string    `long:"topic"`
	Start     Timestamp `long:"start"`
	Format    Format    `long:"format" default:"json"`
	TimeField string    `long:"time-field"`
	Retention int       `long:"retention"` // millisecond
}

func (cmd *RollbackCommand) Exec(cl model.Cluster) error {
	if cmd.Format == "" {
		var err error
		cmd.Format, err = formatDetector{cl, cmd.Topic}.detect()
		if err != nil {
			return err
		}
	}
	partitions, err := cl.Partitions(cmd.Topic)
	if err != nil {
		return err
	}
	cr := consumer.New(cl)
	var wg sync.WaitGroup
	wg.Add(len(partitions))
	timeFunc := cmd.Format.unmarshalTime(cmd.TimeField)
	for _, partition := range partitions {
		go func(partition int32) {
			defer wg.Done()
			offset, err := cr.SearchOffsetByTime(cmd.Topic, partition, time.Time(cmd.Start), timeFunc)
			if err != nil {
				log.Println(err)
				return
			}
			if err := cr.Commit(cmd.Topic, partition, cmd.Group, offset); err != nil {
				log.Println(err)
			}
			fmt.Printf("\t%d:%d\n", partition, offset)
		}(partition)
	}
	wg.Wait()
	return nil

}
