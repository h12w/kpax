package broker

/*
import (
	"net"
	"sync"
	"sync/atomic"
	"time"

	"gopkg.in/fatih/pool.v2"
	"h12.io/kpax/model"
)

const defaultTimeout = 30 * time.Second

type PoolBroker struct {
	Timeout time.Duration
	MaxCap  int
	addr    string
	cid     int32

	p  pool.Pool
	mu sync.Mutex
}

func NewPoolBroker(addr string) *PoolBroker {
	return &PoolBroker{
		addr:    addr,
		MaxCap:  10,
		Timeout: defaultTimeout,
	}
}

func (b *PoolBroker) Do(req model.Request, resp model.Response) error {
	var (
		err error
		p   pool.Pool
	)
	b.mu.Lock()
	if b.p == nil {
		b.p, err = pool.NewChannelPool(0, b.MaxCap, func() (net.Conn, error) {
			return net.DialTimeout("tcp", b.addr, b.Timeout)
		})
	}
	p = b.p
	b.mu.Unlock()
	if err != nil {
		return err
	}

	conn, err := p.Get()
	if err != nil {
		return err
	}
	defer conn.Close()
	if err := conn.SetWriteDeadline(time.Now().Add(b.Timeout)); err != nil {
		if poolConn, ok := conn.(pool.PoolConn); ok {
			poolConn.MarkUnusable()
		}
		return err
	}
	req.SetID(atomic.AddInt32(&b.cid, 1))
	if err := req.Send(conn); err != nil {
		if poolConn, ok := conn.(pool.PoolConn); ok {
			poolConn.MarkUnusable()
		}
		return err
	}
	if resp == nil {
		return nil
	}
	if err := conn.SetReadDeadline(time.Now().Add(b.Timeout)); err != nil {
		if poolConn, ok := conn.(pool.PoolConn); ok {
			poolConn.MarkUnusable()
		}
		return err
	}
	if err := resp.Receive(conn); err != nil {
		if poolConn, ok := conn.(pool.PoolConn); ok {
			poolConn.MarkUnusable()
		}
		return err
	}
	if resp.ID() != req.ID() {
		if poolConn, ok := conn.(pool.PoolConn); ok {
			poolConn.MarkUnusable()
		}
		return errCorrelationIDMismatch
	}
	return nil
}

func (b *PoolBroker) Close() {
	b.mu.Lock()
	if b.p != nil {
		b.p.Close()
		b.p = nil
	}
	b.mu.Unlock()
}
*/
