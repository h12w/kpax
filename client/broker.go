package client

import (
	"errors"
	"net"
	"sync/atomic"

	"h12.me/kafka/proto"
)

var (
	ErrCorrelationIDMismatch = errors.New("correlationID mismatch")
)

type BrokerConfig struct {
	Conn         net.Conn
	SendChanSize int
	RecvChanSize int
}

type Broker struct {
	conn     net.Conn
	cid      int32
	sendChan chan *brokerJob
	recvChan chan *brokerJob
}

type brokerJob struct {
	req     *proto.Request
	resp    *proto.Response
	errChan chan error
}

func NewBroker(c *BrokerConfig) *Broker {
	b := &Broker{
		conn:     c.Conn,
		sendChan: make(chan *brokerJob, c.SendChanSize),
		recvChan: make(chan *brokerJob, c.RecvChanSize),
	}
	go b.sendLoop()
	go b.receiveLoop()
	return b
}

func (b *Broker) Do(req *proto.Request, resp *proto.Response) error {
	req.CorrelationID = atomic.AddInt32(&b.cid, 1)
	errChan := make(chan error)
	b.sendChan <- &brokerJob{
		req:     req,
		resp:    resp,
		errChan: errChan,
	}
	return <-errChan
}

func (b *Broker) sendLoop() {
	for job := range b.sendChan {
		if err := job.req.Send(b.conn); err != nil {
			job.errChan <- err
		}
		b.recvChan <- job
	}
}

func (b *Broker) receiveLoop() {
	for job := range b.recvChan {
		if err := job.resp.Receive(b.conn); err != nil {
			job.errChan <- err
			continue
		}
		if job.resp.CorrelationID != job.req.CorrelationID {
			job.errChan <- ErrCorrelationIDMismatch
		}
		job.errChan <- nil
	}
}
