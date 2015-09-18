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
	config *Config
	cid    int32
	conn   *connection
	mu     sync.Mutex
}

type brokerJob struct {
	req     *proto.Request
	resp    *proto.Response
	errChan chan error
}

func New(config *Config) *B {
	return &B{config: config}
}

func (b *B) Addr() string {
	return b.config.Addr
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
	if b.conn == nil || b.conn.Closed() {
		conn, err := b.newConn()
		if err != nil {
			b.mu.Unlock()
			return err
		}
		b.conn = conn
	}
	b.mu.Unlock()
	b.conn.sendChan <- job
	return nil
}

func (b *B) newConn() (*connection, error) {
	conn, err := net.Dial("tcp", b.config.Addr)
	if err != nil {
		return nil, err
	}
	c := &connection{
		conn:     conn,
		sendChan: make(chan *brokerJob),
		recvChan: make(chan *brokerJob, b.config.QueueLen),
		timeout:  b.config.Timeout,
	}
	go c.sendLoop()
	go c.receiveLoop()
	return c, nil
}

func (b *B) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.conn != nil && !b.conn.Closed() {
		b.conn.Close()
		b.conn = nil
	}
}
