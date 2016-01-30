package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"gopkg.in/vmihailenco/msgpack.v2"
	"h12.me/kpax/consumer"
	"h12.me/kpax/model"
	"h12.me/kpax/proto"
)

type Config struct {
	ConfigFile string  `long:"config" default:"config.json"`
	Brokers    Brokers `long:"brokers" json:"brokers"`

	Consume  ConsumeCommand  `command:"consume" description:"print or count messages in a time range"`
	Offset   OffsetCommand   `command:"offset" description:"print stored offset"`
	Rollback RollbackCommand `command:"rollback" description:"commit older offset of a specific time"`

	Meta  MetaConfig  `command:"meta"`
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
		tm, err := time.Parse("2006-01-02T15:04:05", value)
		if err != nil {
			return fmt.Errorf("error parsing %s: %s", value, err.Error())
		}
		*t = Timestamp(tm)
	}
	return nil
}

type TimeUnmarshalFunc func([]byte) (time.Time, error)

func unmarshalTime(format, field string) func([]byte) (time.Time, error) {
	switch format {
	case "json":
		return func(msg []byte) (time.Time, error) {
			m := make(map[string]interface{})
			if err := json.Unmarshal(msg, &m); err != nil {
				return time.Time{}, err
			}
			timeField, _ := m[field].(string)
			return parseTime(timeField)
		}
	case "url":
		return func(msg []byte) (time.Time, error) {
			values, err := url.ParseQuery(string(msg))
			if err != nil {
				return time.Time{}, err
			}
			return parseTime(values.Get(field))
		}
	case "msgpack":
		return func(msg []byte) (time.Time, error) {
			m := make(map[string]interface{})
			if err := msgpack.Unmarshal(msg, &m); err != nil {
				return time.Time{}, err
			}
			timeField, _ := m[field].(string)
			return parseTime(timeField)
		}
	}
	return nil
}

// parseTime parses time as RFC3339 or unix timestamp
func parseTime(timeText string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339Nano, timeText)
	if err != nil {
		unix, err := strconv.Atoi(timeText)
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(int64(unix), 0), err
	}
	return t, err
}

type ConsumeCommand struct {
	Topic     string    `long:"topic"`
	Start     Timestamp `long:"start"`
	End       Timestamp `long:"end"`
	Format    string    `long:"format" default:"json"`
	TimeField string    `long:"time-field"`
	Count     bool      `long:"count"`
}

func (cmd *ConsumeCommand) Exec(cl model.Cluster) error {
	// TODO: detect format
	partitions, err := cl.Partitions(cmd.Topic)
	if err != nil {
		return err
	}
	cr := consumer.New(cl)
	var wg sync.WaitGroup
	wg.Add(len(partitions))
	var cnt int64
	timeFunc := unmarshalTime(cmd.Format, cmd.TimeField)
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
					fmt.Println(string(msg.Value))
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

type OffsetCommand struct {
	Group string `long:"group"`
	Topic string `long:"topic"`
}

func (cmd *OffsetCommand) Exec(cl model.Cluster) error {
	partitions, err := cl.Partitions(cmd.Topic)
	if err != nil {
		return err
	}
	cr := consumer.New(cl)
	fmt.Printf("topic: %s, group: %s\n", cmd.Topic, cmd.Group)
	for _, partition := range partitions {
		offset, err := cr.Offset(cmd.Topic, partition, cmd.Group)
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
	Format    string    `long:"format" default:"json"`
	TimeField string    `long:"time-field"`
	Retention int       `long:"retention"` // millisecond
}

func (cmd *RollbackCommand) Exec(cl model.Cluster) error {
	// TODO: detect format
	partitions, err := cl.Partitions(cmd.Topic)
	if err != nil {
		return err
	}
	cr := consumer.New(cl)
	var wg sync.WaitGroup
	wg.Add(len(partitions))
	timeFunc := unmarshalTime(cmd.Format, cmd.TimeField)
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
