package broker

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"h12.me/kpax/model"
)

type B struct {
	Addr     string
	Timeout  time.Duration
	QueueLen int

	br *broker
	mu sync.Mutex
}

func New(addr string) *B {
	b := &B{
		Addr:     addr,
		Timeout:  30 * time.Second,
		QueueLen: 1000,
	}
	return b
}

func NewDefault(addr string) model.Broker { return New(addr) }

func (b *B) Do(req model.Request, resp model.Response) error {
	b.mu.Lock()
	if b.br == nil {
		var err error
		b.br, err = b.newBroker()
		if err != nil {
			b.mu.Unlock()
			return err
		}
	}
	b.mu.Unlock()

	if err := b.br.do(req, resp); err != nil {
		b.Close()
		return err
	}
	return nil
}

func (b *B) Close() {
	b.mu.Lock()
	b.br.close()
	b.br = nil
	b.mu.Unlock()
}

type broker struct {
	timeout  time.Duration
	conn     net.Conn
	recvChan chan *brokerJob
	cid      int32
	mu       sync.Mutex
}

type brokerJob struct {
	req     model.Request
	resp    model.Response
	errChan chan error
}

func (b *B) newBroker() (*broker, error) {
	conn, err := net.DialTimeout("tcp", b.Addr, b.Timeout)
	if err != nil {
		return nil, err
	}
	br := &broker{
		timeout:  b.Timeout,
		conn:     conn,
		recvChan: make(chan *brokerJob, b.QueueLen),
	}
	go br.receiveLoop()
	return br, nil
}

func (b *broker) do(req model.Request, resp model.Response) error {
	job := &brokerJob{
		req:     req,
		resp:    resp,
		errChan: make(chan error),
	}
	if err := b.send(job); err != nil {
		return err
	}
	if job.requireAck() {
		return <-job.errChan
	}
	return nil
}

func (b *broker) send(job *brokerJob) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	job.req.SetID(atomic.AddInt32(&b.cid, 1))
	if err := b.conn.SetWriteDeadline(time.Now().Add(b.timeout)); err != nil {
		return err
	}
	if err := job.req.Send(b.conn); err != nil {
		return err
	}
	if job.requireAck() {
		b.recvChan <- job
	}
	return nil
}

var (
	errCorrelationIDMismatch = errors.New("correlationID mismatch")
)

func (b *broker) receiveLoop() {
	for job := range b.recvChan {
		if err := b.conn.SetReadDeadline(time.Now().Add(b.timeout)); err != nil {
			job.errChan <- err
			continue
		}
		if err := job.resp.Receive(b.conn); err != nil {
			job.errChan <- err
			continue
		}
		if job.resp.ID() != job.req.ID() {
			job.errChan <- errCorrelationIDMismatch
			continue
		}
		job.errChan <- nil
	}
	// no need to closeConn here because when recvChan closed, the receiveLoop will do it.
	b.conn.Close()
	b.conn = nil
}

func (b *broker) close() {
	close(b.recvChan)
	b.recvChan = nil
}

func (j *brokerJob) requireAck() bool { return j.resp != nil }
