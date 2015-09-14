package client

import (
	"h12.me/kafka/proto"
)

type Config struct {
	Brokers []string
}

type Client struct {
}

func New() *Client {
	return &Client{}
}

func (c *Client) Produce(req *proto.ProduceRequest) error {
	return nil
}
