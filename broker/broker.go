package broker

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"h12.me/kpax/model"
)

var (
	ErrCorrelationIDMismatch = errors.New("correlationID mismatch")
)

type B struct {
	Addr     string
	Timeout  time.Duration
	QueueLen int

	connection
	cid      int32
	recvChan chan *brokerJob
	mu       sync.Mutex
}

type connection struct {
	conn net.Conn
	mu   sync.Mutex
	*B
}

type brokerJob struct {
	req     model.Request
	resp    model.Response
	errChan chan error
}

func New(addr string) *B {
	b := &B{
		Addr:     addr,
		Timeout:  30 * time.Second,
		QueueLen: 1000,
	}
	b.connection.B = b
	return b
}

func NewDefault(addr string) model.Broker { return New(addr) }

func (b *B) Do(req model.Request, resp model.Response) error {
	job := &brokerJob{
		req:     req,
		resp:    resp,
		errChan: make(chan error),
	}
	if err := b.send(job); err != nil {
		return err
	}
	return <-job.errChan
}

func (b *B) send(job *brokerJob) error {
	conn, err := b.getConn()
	if err != nil {
		return err
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	job.req.SetID(atomic.AddInt32(&b.cid, 1))
	if err := job.req.Send(conn); err != nil {
		b.closeConn()
		return err
	}
	if b.recvChan == nil {
		b.recvChan = make(chan *brokerJob, b.QueueLen)
		go b.receiveLoop()
	}
	if job.requireAck() {
		b.recvChan <- job
	}
	return nil
}
func (j *brokerJob) requireAck() bool { return j.resp != nil }

func (c *B) receiveLoop() {
	for job := range c.recvChan {
		conn, err := c.getConn()
		if err != nil {
			job.errChan <- err
			continue
		}
		if err := job.resp.Receive(conn); err != nil {
			c.closeConn()
			job.errChan <- err
			continue
		}
		if job.resp.ID() != job.req.ID() {
			c.closeConn()
			job.errChan <- ErrCorrelationIDMismatch
			continue
		}
		job.errChan <- nil
	}
	c.closeConn()
}

func (b *B) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.recvChan != nil {
		close(b.recvChan)
		b.recvChan = nil
	}
}

func (c *connection) getConn() (net.Conn, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil {
		conn, err := net.DialTimeout("tcp", c.Addr, c.Timeout)
		if err != nil {
			return nil, err
		}
		c.conn = conn
	}
	return c.conn, c.conn.SetDeadline(time.Now().Add(c.Timeout))
}

func (c *connection) closeConn() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}
