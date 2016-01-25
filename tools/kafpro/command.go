package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"h12.me/kafka/broker"
	"h12.me/kafka/cluster"
	"h12.me/kafka/consumer"
)

type Config struct {
	Brokers Brokers      `long:"brokers"`
	Meta    MetaConfig   `command:"meta"`
	Coord   CoordConfig  `command:"coord"`
	Offset  OffsetConfig `command:"offset"`
	Commit  CommitConfig `command:"commit"`

	Time TimeCommand `command:"time" description:"get offset based on time"`

	Consume ConsumeCommand `command:"consume"`
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

type TimeCommand struct {
	Topic     string    `long:"topic"`
	Start     Timestamp `long:"start"`
	TimeField string    `long:"time-field"`
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

func (cmd *TimeCommand) Exec(cl *cluster.C) error {
	// TODO: detect format
	partitions, err := cl.Partitions(cmd.Topic)
	if err != nil {
		return err
	}
	cr, err := consumer.New(consumer.DefaultConfig(), cl)
	if err != nil {
		return err
	}
	timeFunc := unmarshalTime("json", cmd.TimeField)
	for _, partition := range partitions {
		offset, err := cr.SearchOffsetByTime(cmd.Topic, partition, time.Time(cmd.Start), timeFunc)
		if err != nil {
			return err
		}
		fmt.Println(partition, offset)
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
	timeFunc := unmarshalTime("json", cmd.TimeField)
	for _, partition := range partitions {
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
				fmt.Println(string(msg.Value))
				t, err := timeFunc(msg.Value)
				if err != nil {
					return err
				}
				if t.After(time.Time(cmd.End)) {
					jitterCnt++
				}
			}
			offset = messages[len(messages)-1].Offset + 1
		}
	}
	return nil
	/*
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
	*/
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
