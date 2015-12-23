package broker

import (
	"time"
)

func (b *B) TopicMetadata(topics ...string) (*TopicMetadataResponse, error) {
	reqMsg := TopicMetadataRequest(topics)
	req := &Request{
		RequestMessage: &reqMsg,
	}
	respMsg := &TopicMetadataResponse{}
	if err := b.Do(req, respMsg); err != nil {
		return nil, err
	}
	return &respMsg, nil
}

func (b *B) GroupCoordinator(group string) (*GroupCoordinatorResponse, error) {
	reqMsg := GroupCoordinatorRequest(group)
	req := &Request{
		RequestMessage: &reqMsg,
	}
	respMsg := &GroupCoordinatorResponse{}
	if err := b.Do(req, respMsg); err != nil {
		return nil, err
	}
	if respMsg.ErrorCode.HasError() {
		return nil, respMsg.ErrorCode
	}
	return respMsg, nil
}

func (b *B) Produce(topic string, partition int32, messageSet []OffsetMessage) (*ProduceResponse, error) {
	cfg := &b.config.Producer
	req := &Request{
		RequestMessage: &ProduceRequest{
			RequiredAcks: cfg.RequiredAcks,
			Timeout:      int32(cfg.Timeout / time.Millisecond),
			MessageSetInTopics: []MessageSetInTopic{
				{
					TopicName: topic,
					MessageSetInPartitions: []MessageSetInPartition{
						{
							Partition:  partition,
							MessageSet: messageSet,
						},
					},
				},
			},
		},
	}

	respMsg := &ProduceResponse{}
	if err := b.Do(req, respMsg); err != nil {
		return nil, err
	}
	return respMsg, nil
}
