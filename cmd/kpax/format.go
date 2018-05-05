package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"gopkg.in/vmihailenco/msgpack.v2"
	"h12.io/kpax/consumer"
	"h12.io/kpax/model"
	"h12.io/kpax/proto"
	"h12.io/uuid/hexid"
)

type Format string

const (
	UnknownFormat Format = ""
	MsgPackFormat Format = "msgpack"
	JSONFormat    Format = "json"
	URLFormat     Format = "url"
)

type formatDetector struct {
	cl    model.Cluster
	topic string
}

func (d formatDetector) detect() (Format, error) {
	cl, topic := d.cl, d.topic
	partitions, err := cl.Partitions(topic)
	if err != nil {
		return UnknownFormat, err
	}
	cr := consumer.New(cl)
	for _, partition := range partitions {
		offset, err := cr.FetchOffsetByTime(topic, partition, proto.Latest)
		if err != nil {
			continue
		}
		messages, err := cr.Consume(topic, partition, offset-1)
		if err != nil {
			continue
		}
		for _, msg := range messages {
			return detectFormat(msg.Value)
		}
	}
	return UnknownFormat, nil
}
func detectFormat(value []byte) (Format, error) {
	m := make(map[string]interface{})
	if err := msgpack.Unmarshal(value, &m); err == nil {
		return MsgPackFormat, nil
	} else if err := json.Unmarshal(value, &m); err == nil {
		return JSONFormat, nil
	} else if _, err := url.ParseQuery(string(value)); err == nil {
		return URLFormat, nil
	}
	return UnknownFormat, nil
}

func (format Format) Sprint(value []byte) (string, error) {
	switch format {
	case MsgPackFormat:
		m := make(map[string]interface{})
		if err := msgpack.Unmarshal(value, &m); err != nil {
			break
		}
		buf, err := json.Marshal(hexid.Restore(m))
		if err != nil {
			return "", fmt.Errorf("fail to marshal %#v: %s", m, err.Error())
		}
		return string(buf), nil
	}
	return string(value), nil
}

func (format Format) marshal(v interface{}) ([]byte, error) {
	switch format {
	case MsgPackFormat:
		return msgpack.Marshal(v)
	case JSONFormat:
		return json.Marshal(v)
	}
	return []byte(fmt.Sprint(v)), nil
}

type TimeUnmarshalFunc func([]byte) (time.Time, error)

func (format Format) unmarshalTime(field string) func([]byte) (time.Time, error) {
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
			m = hexid.Restore(m).(map[string]interface{})
			return parseTime(m[field])
		}
	}
	return nil
}
