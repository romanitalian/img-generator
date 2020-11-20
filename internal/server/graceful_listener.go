package server

import (
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

type GracefulListener struct {
	ln          net.Listener
	maxWaitTime time.Duration
	done        chan struct{}
	connsCount  uint64
	shutdown    uint64
}

func NewGracefulListener(ln net.Listener, maxWaitTime time.Duration) net.Listener {
	return &GracefulListener{
		ln:          ln,
		maxWaitTime: maxWaitTime,
		done:        make(chan struct{}),
	}
}

func (ln *GracefulListener) Accept() (net.Conn, error) {
	c, err := ln.ln.Accept()
	if err != nil {
		return nil, err
	}
	atomic.AddUint64(&ln.connsCount, 1)
	return &gracefulConn{
		Conn: c,
		ln:   ln,
	}, nil
}

func (ln *GracefulListener) Addr() net.Addr {
	return ln.ln.Addr()
}

func (ln *GracefulListener) Close() error {
	err := ln.ln.Close()
	if err != nil {
		return err
	}

	atomic.AddUint64(&ln.shutdown, 1)
	if atomic.LoadUint64(&ln.connsCount) == 0 {
		close(ln.done)
		return nil
	}
	select {
	case <-ln.done:
		return nil
	case <-time.After(ln.maxWaitTime):
		return fmt.Errorf("cannot graceful shutdown in waiting time %s", ln.maxWaitTime)
	}
}

func (ln *GracefulListener) closeConn() {
	cnt := atomic.AddUint64(&ln.connsCount, ^uint64(0))
	if atomic.LoadUint64(&ln.shutdown) != 0 && cnt == 0 {
		close(ln.done)
	}
}

type gracefulConn struct {
	net.Conn
	ln *GracefulListener
}

func (c *gracefulConn) Close() error {
	err := c.Conn.Close()
	if err != nil {
		return err
	}
	cnt := atomic.AddUint64(&c.ln.connsCount, ^uint64(0))
	if atomic.LoadUint64(&c.ln.shutdown) != 0 && cnt == 0 {
		close(c.ln.done)
	}

	return nil
}
