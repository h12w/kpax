package broker

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"h12.me/kpax/model"
)

type AsyncBroker struct {
	Timeout  time.Duration
	QueueLen int
	addr     string

	mu sync.Mutex
	br *broker
}

func NewAsyncBroker(addr string) *AsyncBroker {
	b := &AsyncBroker{
		addr:     addr,
		Timeout:  30 * time.Second,
		QueueLen: 1000,
	}
	return b
}

func NewDefault(addr string) model.Broker { return NewAsyncBroker(addr) }

func (b *AsyncBroker) Do(req model.Request, resp model.Response) error {
	var (
		br  *broker
		err error
	)
	b.mu.Lock()
	if b.br == nil {
		b.br, err = b.newBroker()
	}
	br = b.br
	b.mu.Unlock()
	if err != nil {
		return err
	}

	if err := <-br.do(req, resp); err != nil {
		b.mu.Lock()
		if b.br == br {
			b.br = nil
		}
		b.mu.Unlock()
		br.close()
		return err
	}
	return nil
}

func (b *AsyncBroker) Close() {
	b.mu.Lock()
	if b.br != nil {
		b.br.close()
		b.br = nil
	}
	b.mu.Unlock()
}

type broker struct {
	timeout time.Duration
	conn    net.Conn
	cid     int32

	mu       sync.Mutex
	recvChan chan *brokerJob
}

type brokerJob struct {
	req     model.Request
	resp    model.Response
	errChan chan error
}

func (b *AsyncBroker) newBroker() (*broker, error) {
	conn, err := net.DialTimeout("tcp", b.addr, b.Timeout)
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

func (b *broker) do(req model.Request, resp model.Response) chan error {
	job := &brokerJob{
		req:     req,
		resp:    resp,
		errChan: make(chan error),
	}
	if err := b.send(job); err != nil {
		go func() {
			job.errChan <- err
		}()
	}
	return job.errChan
}

func (b *broker) send(job *brokerJob) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.recvChan == nil {
		return errChannelAlreadyClosed
	}
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

func (b *broker) close() {
	b.mu.Lock()
	if b.recvChan != nil {
		close(b.recvChan)
		b.recvChan = nil
	}
	b.mu.Unlock()
}

var (
	errCorrelationIDMismatch = errors.New("correlationID mismatch")
	errChannelAlreadyClosed  = errors.New("channel already closed")
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

func (j *brokerJob) requireAck() bool { return j.resp != nil }
