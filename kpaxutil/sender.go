package kpaxutil

import (
	"encoding"

	"h12.io/kpax/broker"
	"h12.io/kpax/cluster"
	"h12.io/kpax/producer"
)

type Sender interface {
	Send(topic string, value encoding.BinaryMarshaler) error
}

type simpleSender struct {
	p *producer.P
}

func NewSender(brokers []string) Sender {
	return &simpleSender{p: producer.New(cluster.New(broker.NewDefault, brokers))}
}

func (s *simpleSender) Send(topic string, value encoding.BinaryMarshaler) error {
	buf, err := value.MarshalBinary()
	if err != nil {
		return err
	}
	return s.p.Produce(topic, nil, buf)
}
