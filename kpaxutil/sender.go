package kpaxutil

import (
	"encoding"

	"h12.me/kpax/broker"
	"h12.me/kpax/cluster"
	"h12.me/kpax/producer"
)

type Sender struct {
	p *producer.P
}

func NewSender(brokers []string) *Sender {
	return &Sender{p: producer.New(cluster.New(broker.New, brokers))}
}

func (s *Sender) Send(topic string, value encoding.BinaryMarshaler) error {
	buf, err := value.MarshalBinary()
	if err != nil {
		return err
	}
	return s.p.Produce(topic, nil, buf)
}
