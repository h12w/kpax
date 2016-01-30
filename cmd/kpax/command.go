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

	"h12.me/kpax/consumer"
	"h12.me/kpax/model"
	"h12.me/kpax/proto"
)

type Config struct {
	ConfigFile string  `long:"config" default:"config.json"`
	Brokers    Brokers `long:"brokers" json:"brokers"`

	Consume ConsumeCommand `command:"consume"`
	Commit  CommitCommand  `command:"commit"`

	Meta   MetaConfig   `command:"meta"`
	Coord  CoordConfig  `command:"coord"`
	Offset OffsetConfig `command:"offset"`
}

type Brokers []string

func (t *Brokers) UnmarshalFlag(value string) error {
	*t = strings.Split(value, ",")
	return nil
}

type CoordConfig struct {
	GroupName string `long:"group"`
}

type OffsetConfig struct {
	GroupName string `long:"group"`
	Topic     string `long:"topic"`
	Partition int    `long:"partition"`
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

type CommitCommand struct {
	GroupName string    `long:"group"`
	Topic     string    `long:"topic"`
	Start     Timestamp `long:"start"`
	Format    string    `long:"format" default:"json"`
	TimeField string    `long:"time-field"`
	Retention int       `long:"retention"` // millisecond
}

func (cmd *CommitCommand) Exec(cl model.Cluster) error {
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
			if err := cr.Commit(cmd.Topic, partition, cmd.GroupName, offset); err != nil {
				log.Println(err)
			}
		}(partition)
	}
	wg.Wait()
	return nil

}
