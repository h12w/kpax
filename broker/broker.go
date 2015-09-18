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
	open     bool
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

func New(config *Config) *B {
	return &B{
		config:   config,
		sendChan: make(chan *brokerJob),
		recvChan: make(chan *brokerJob, config.QueueLen),
	}
}

func (b *B) Addr() string {
	return b.config.Addr
}

func (b *B) connect() error {
	if b.open {
		return nil
	}
	var err error
	b.conn, err = net.Dial("tcp", b.config.Addr)
	if err != nil {
		return err
	}
	b.open = true
	go b.sendLoop()
	go b.receiveLoop()
	return nil
}

func (b *B) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.open {
		b.open = false
		close(b.sendChan)
		b.conn.Close()
	}
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
	defer b.mu.Unlock()
	if err := b.connect(); err != nil {
		return err
	}
	b.sendChan <- job
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
