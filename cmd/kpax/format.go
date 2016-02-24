package main

import (
	"encoding/json"
	"fmt"
	"net/url"

	"gopkg.in/vmihailenco/msgpack.v2"
	"h12.me/kpax/consumer"
	"h12.me/kpax/model"
	"h12.me/kpax/proto"
	"h12.me/uuid/hexid"
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
