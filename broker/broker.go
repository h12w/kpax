package broker

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrCorrelationIDMismatch = errors.New("correlationID mismatch")
	ErrBrokerClosed          = errors.New("broker is closed")
)

type Config struct {
	QueueLen   int
	Connection ConnectionConfig
	Producer   ProducerConfig
}

type ConnectionConfig struct {
	Addr    string
	Timeout time.Duration
}

type ProducerConfig struct {
	RequiredAcks int16
	Timeout      time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		QueueLen: 1000,
		Connection: ConnectionConfig{
			Timeout: 30 * time.Second,
		},
		Producer: ProducerConfig{
			RequiredAcks: AckLocal,
			Timeout:      10 * time.Second,
		},
	}
}

func (c *Config) WithAddr(addr string) *Config {
	c.Connection.Addr = addr
	return c
}

type B struct {
	config *Config
	cid    int32
	conn   *connection
	mu     sync.Mutex
}

type brokerJob struct {
	req     *Request
	resp    *Response
	errChan chan error
}

func New(config *Config) *B {
	return &B{config: config}
}

func (b *B) Addr() string {
	return b.config.Connection.Addr
}

func (b *B) Do(req *Request, resp ResponseMessage) error {
	req.CorrelationID = atomic.AddInt32(&b.cid, 1)
	errChan := make(chan error)
	if err := b.sendJob(&brokerJob{
		req:     req,
		resp:    &Response{ResponseMessage: resp},
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
	cfg := &b.config.Connection
	conn, err := net.Dial("tcp", cfg.Addr)
	if err != nil {
		return nil, err
	}
	c := &connection{
		conn:     conn,
		sendChan: make(chan *brokerJob),
		recvChan: make(chan *brokerJob, b.config.QueueLen),
		timeout:  cfg.Timeout,
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
