package producer

import (
	"fmt"
	"math/rand"
	"time"

	"h12.me/kafka/client"
	"h12.me/kafka/proto"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Config struct {
	Client       client.Config
	RequiredAcks int16
	Timeout      int32
}

type P struct {
	client *client.C
	config *Config
}

func New(config *Config) (*P, error) {
	client, err := client.New(&config.Client)
	if err != nil {
		return nil, err
	}
	return &P{
		client: client,
		config: config,
	}, nil
}

func (p *P) getPartition(topic string) (int32, error) {
	partitions, err := p.client.Partitions(topic)
	if err != nil {
		return 0, err
	}
	return partitions[rand.Intn(len(partitions))], nil
}

func (p *P) Produce(topic string, key, value []byte) error {
	partition, err := p.getPartition(topic)
	if err != nil {
		return err
	}
	leader, err := p.client.Leader(topic, partition)
	if err != nil {
		return err
	}
	req := &proto.Request{
		APIKey:        proto.ProduceRequestType,
		APIVersion:    0,
		CorrelationID: 1,
		ClientID:      p.config.Client.ClientID,
		RequestMessage: &proto.ProduceRequest{
			RequiredAcks: p.config.RequiredAcks,
			Timeout:      p.config.Timeout,
			MessageSetInTopics: []proto.MessageSetInTopic{
				{
					TopicName: topic,
					MessageSetInPartitions: []proto.MessageSetInPartition{
						{
							Partition: partition,
							MessageSet: []proto.OffsetMessage{
								{
									SizedMessage: proto.SizedMessage{CRCMessage: proto.CRCMessage{
										Message: proto.Message{
											Key:   key,
											Value: value,
										},
									}}},
							},
						},
					},
				},
			},
		},
	}
	resp := proto.ProduceResponse{}
	if err := leader.Do(req, &proto.Response{ResponseMessage: &resp}); err != nil {
		return err
	}
	for i := range resp {
		for j := range resp[i].OffsetInPartitions {
			if errCode := resp[i].OffsetInPartitions[j].ErrorCode; errCode != 0 {
				return fmt.Errorf("Kafka error %d", errCode)
			}
		}
	}
	return nil
}
