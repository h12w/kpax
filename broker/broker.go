package broker

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"h12.me/kafka/proto"
)

var (
	ErrCorrelationIDMismatch = errors.New("correlationID mismatch")
	ErrBrokerClosed          = errors.New("broker is closed")
)

type Config struct {
	Addr     string
	QueueLen int
	Timeout  time.Duration
}

type B struct {
	config   *Config
	conn     net.Conn
	closed   bool
	cid      int32
	sendChan chan *brokerJob
	recvChan chan *brokerJob
	mu       sync.Mutex
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
		sendChan: make(chan *brokerJob),
		recvChan: make(chan *brokerJob, config.QueueLen),
	}
	go b.sendLoop()
	go b.receiveLoop()
	return b, nil
}

func (b *B) Close() {
	b.mu.Lock()
	if !b.closed {
		b.closed = true
		close(b.sendChan)
		b.conn.Close()
	}
	b.mu.Unlock()
}

func (b *B) Do(req *proto.Request, resp proto.ResponseMessage) error {
	req.CorrelationID = atomic.AddInt32(&b.cid, 1)
	errChan := make(chan error)
	if err := b.sendJob(&brokerJob{
		req:     req,
		resp:    &proto.Response{ResponseMessage: resp},
		errChan: errChan,
	}); err != nil {
		return err
	}
	return <-errChan
}

func (b *B) sendJob(job *brokerJob) error {
	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		return ErrBrokerClosed
	}
	b.sendChan <- job
	b.mu.Unlock()
	return nil
}

func (b *B) sendLoop() {
	for job := range b.sendChan {
		b.conn.SetWriteDeadline(time.Now().Add(b.config.Timeout))
		if err := job.req.Send(b.conn); err != nil {
			b.Close()
			job.errChan <- err
			close(b.recvChan)
		}
		b.recvChan <- job
	}
	close(b.recvChan)
}

func (b *B) receiveLoop() {
	for job := range b.recvChan {
		b.conn.SetReadDeadline(time.Now().Add(b.config.Timeout))
		if err := job.resp.Receive(b.conn); err != nil {
			b.Close()
			job.errChan <- err
		}
		if job.resp.CorrelationID != job.req.CorrelationID {
			job.errChan <- ErrCorrelationIDMismatch
		}
		job.errChan <- nil
	}
}
