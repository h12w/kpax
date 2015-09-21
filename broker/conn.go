package broker

import (
	"net"
	"sync"
	"time"

	"h12.me/kafka/proto"
)

type connection struct {
	conn     net.Conn
	timeout  time.Duration
	sendChan chan *brokerJob
	recvChan chan *brokerJob
	mu       sync.Mutex
}

func (c *connection) Closed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.sendChan == nil
}

func (c *connection) Close() {
	c.mu.Lock()
	close(c.sendChan)
	c.sendChan = nil
	c.mu.Unlock()
}

func (c *connection) sendLoop() {
	for job := range c.sendChan {
		c.conn.SetWriteDeadline(time.Now().Add(c.timeout))
		if err := job.req.Send(c.conn); err != nil {
			job.errChan <- err
			if err == proto.ErrConn {
				c.Close()
				close(c.recvChan)
			}
		}
		c.recvChan <- job
	}
	close(c.recvChan)
}

func (c *connection) receiveLoop() {
	for job := range c.recvChan {
		c.conn.SetReadDeadline(time.Now().Add(c.timeout))
		if err := job.resp.Receive(c.conn); err != nil {
			job.errChan <- err
			if err == proto.ErrConn {
				c.Close()
			}
		}
		if job.resp.CorrelationID != job.req.CorrelationID {
			job.errChan <- ErrCorrelationIDMismatch
		}
		job.errChan <- nil
	}
	c.conn.Close() // safe to close the connection here
}
