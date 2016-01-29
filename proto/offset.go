package proto

import (
	"fmt"
	"time"

	"h12.me/kafka/model"
)

type GetTimeFunc func([]byte) (time.Time, error)

func (o *OffsetByTime) Search(cl model.Cluster, getTime GetTimeFunc) (int64, error) {
	earliest, err := (&OffsetByTime{
		Topic:     o.Topic,
		Partition: o.Partition,
		Time:      Earliest,
	}).Fetch(cl)
	if err != nil {
		return -1, err
	}
	if o.Time == Earliest {
		return earliest, nil
	}
	latest, err := (&OffsetByTime{
		Topic:     o.Topic,
		Partition: o.Partition,
		Time:      Latest,
	}).Fetch(cl)
	if err != nil {
		return -1, err
	}

	const maxMessageSize = 10000
	getter := timeGetter{
		Messages{
			Topic:       o.Topic,
			Partition:   o.Partition,
			MinBytes:    1,
			MaxBytes:    maxMessageSize,
			MaxWaitTime: 100 * time.Millisecond,
		},
		cl,
		getTime,
	}
	min, max := earliest, latest
	mid := min + (max-min)/2
	for min <= max {
		mid = min + (max-min)/2
		midTime, err := getter.get(mid)
		if err != nil {
			return -1, err
		}
		if midTime.Equal(o.Time) {
			break
		} else if midTime.Before(o.Time) {
			min = mid + 1
		} else {
			max = mid - 1
		}
	}
	return o.searchOffsetBefore(cl, earliest, mid, latest, getter)
}

func (o *OffsetByTime) searchOffsetBefore(cl model.Cluster, min, mid, max int64, getter timeGetter) (int64, error) {
	const maxJitter = 1000 // time may be interleaved in a small range
	midTime, err := getter.get(mid)
	if err != nil {
		return -1, err
	}
	if midTime.Before(o.Time) {
		mid -= maxJitter
	} else {
		mid -= 2 * maxJitter // we have passed the point, go back with double jitter
	}
	if mid < min {
		mid = min
	}
	for offset := mid; offset <= max; offset++ {
		getter.Offset = offset
		messages, err := getter.Consume(cl)
		if err != nil {
			return -1, err
		}
		if len(messages) == 0 {
			return -1, fmt.Errorf("fail to search offset: zero message count")
		}
		for _, message := range messages {
			t, err := getter.getTime(message.Value)
			if err != nil {
				return -1, err
			}
			if t.After(o.Time) || t.Equal(o.Time) {
				return offset, nil
			}
			offset = message.Offset
		}
	}
	return mid, nil
}

type timeGetter struct {
	Messages
	cl      model.Cluster
	getTime GetTimeFunc
}

func (g *timeGetter) get(offset int64) (time.Time, error) {
	g.Offset = offset
	messages, err := g.Consume(g.cl)
	if err != nil {
		return time.Time{}, err
	}
	if len(messages) == 0 {
		return time.Time{}, fmt.Errorf("no messages at topic %s, partition %d, offset %d", g.Topic, g.Partition, g.Offset)
	}
	if messages[0].Offset != offset {
		return time.Time{}, fmt.Errorf("OFFSET MISMATCH!!! %d, %d", messages[0].Offset, offset)
	}
	return g.getTime(messages[0].Value)
}
