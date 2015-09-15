package broker

import (
	"errors"
	"net"
	"sync/atomic"

	"h12.me/kafka/proto"
)

var (
	ErrCorrelationIDMismatch = errors.New("correlationID mismatch")
)

type Config struct {
	Addr         string
	SendQueueLen int
	RecvQueueLen int
}

type B struct {
	config   *Config
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

func New(config *Config) (*B, error) {
	conn, err := net.Dial("tcp", config.Addr)
	if err != nil {
		return nil, err
	}
	b := &B{
		config:   config,
		conn:     conn,
		sendChan: make(chan *brokerJob, config.SendQueueLen),
		recvChan: make(chan *brokerJob, config.RecvQueueLen),
	}
	go b.sendLoop()
	go b.receiveLoop()
	return b, nil
}

func (b *B) Close() {
	b.conn.Close()
}

func (b *B) Do(req *proto.Request, resp *proto.Response) error {
	req.CorrelationID = atomic.AddInt32(&b.cid, 1)
	errChan := make(chan error)
	b.sendChan <- &brokerJob{
		req:     req,
		resp:    resp,
		errChan: errChan,
	}
	return <-errChan
}

func (b *B) sendLoop() {
	for job := range b.sendChan {
		if err := job.req.Send(b.conn); err != nil {
			job.errChan <- err
		}
		b.recvChan <- job
	}
}

func (b *B) receiveLoop() {
	for job := range b.recvChan {
		if err := job.resp.Receive(b.conn); err != nil {
			job.errChan <- err
		}
		if job.resp.CorrelationID != job.req.CorrelationID {
			job.errChan <- ErrCorrelationIDMismatch
		}
		job.errChan <- nil
	}
}
