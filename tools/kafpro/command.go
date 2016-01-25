package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"h12.me/kafka/broker"
	"h12.me/kafka/cluster"
	"h12.me/kafka/consumer"
)

type Config struct {
	ConfigFile string  `long:"config" default:"config.json"`
	Brokers    Brokers `long:"brokers" json:"brokers"`

	Consume ConsumeCommand `command:"consume"`

	Meta   MetaConfig   `command:"meta"`
	Coord  CoordConfig  `command:"coord"`
	Offset OffsetConfig `command:"offset"`
	Commit CommitConfig `command:"commit"`
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
		*t = Timestamp(broker.Latest)
	case "earliest":
		*t = Timestamp(broker.Earliest)
	default:
		tm, err := time.Parse("2006-01-02T15:04:05", value)
		if err != nil {
			return fmt.Errorf("error parsing %s: %s", value, err.Error())
		}
		*t = Timestamp(tm)
	}
	return nil
}

func unmarshalTime(format, field string) func([]byte) (time.Time, error) {
	switch format {
	case "json":
		return func(msg []byte) (time.Time, error) {
			m := make(map[string]interface{})
			if err := json.Unmarshal(msg, &m); err != nil {
				return time.Time{}, err
			}
			timeField, ok := m[field].(string)
			if !ok {
				return time.Time{}, fmt.Errorf("%v does not contains string field %s", m, field)
			}
			t, err := time.Parse(time.RFC3339Nano, timeField)
			if err != nil {
				return time.Time{}, err
			}
			return t, nil
		}
	}
	return nil
}

type ConsumeCommand struct {
	Topic     string    `long:"topic"`
	Start     Timestamp `long:"start"`
	End       Timestamp `long:"end"`
	TimeField string    `long:"time-field"`
	Count     bool      `long:"count"`
}

func (cmd *ConsumeCommand) Exec(cl *cluster.C) error {
	// TODO: detect format
	partitions, err := cl.Partitions(cmd.Topic)
	if err != nil {
		return err
	}
	cr, err := consumer.New(consumer.DefaultConfig(), cl)
	if err != nil {
		return err
	}
	var cnt int64
	for _, partition := range partitions {
		if err := cmd.consumePartition(cr, partition, &cnt); err != nil {
			return err
		}
	}
	fmt.Println(cnt)
	return nil
}

func (cmd *ConsumeCommand) consumePartition(cr *consumer.C, partition int32, cnt *int64) error {
	timeFunc := unmarshalTime("json", cmd.TimeField)
	offset, err := cr.SearchOffsetByTime(cmd.Topic, partition, time.Time(cmd.Start), timeFunc)
	if err != nil {
		return err
	}
	jitterCnt := 0
	for jitterCnt <= 1000 {
		messages, err := cr.Consume(cmd.Topic, partition, offset)
		if err != nil {
			return err
		}
		if len(messages) == 0 {
			break
		}
		for _, msg := range messages {
			t, err := timeFunc(msg.Value)
			if err != nil {
				return err
			}
			if !t.Before(time.Time(cmd.Start)) && t.Before(time.Time(cmd.End)) {
				if cmd.Count {
					atomic.AddInt64(cnt, 1)
				} else {
					fmt.Println(string(msg.Value))
				}
			}
			if t.After(time.Time(cmd.End)) {
				jitterCnt++
			}
		}
		offset = messages[len(messages)-1].Offset + 1
	}
	return nil
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
