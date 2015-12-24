package pool

import (
	"fmt"
	"net"
	"time"
)

//wrappedConn modify the behavior of net.Conn's Write() method and Close() method
//while other methods can be accessed transparently.
type wrappedConn struct {
	net.Conn
	pool             *blockingPool
	unusable         bool
	inPool           bool
	liveTime         time.Duration
	setDelayClose    chan struct{}
	cancelDelayClose chan struct{}
}

// setEnterPool will start delay close function and mark inPool = true
func (c *wrappedConn) setEnterPool() {
	if c.Conn != nil {
		c.setDelayClose <- struct{}{}
	}
	c.inPool = true
}

// setOutPool will cancel delay close function and mark inPool = false
func (c *wrappedConn) setOutPool() {
	if c.Conn != nil {
		c.cancelDelayClose <- struct{}{}
	}
	c.inPool = false
}

//TODO
func (c *wrappedConn) Close() error {
	if c.inPool {
		return fmt.Errorf("close conn fail: conn is in pool")
	}

	if err := c.pool.put(c); err != nil {
		return err
	}

	return nil
}

// closeConnect close inner net.Conn and mark the wrapped connecton unusable
func (c *wrappedConn) closeConnect() {
	if c.Conn != nil {
		c.Conn.Close()
		c.Conn = nil
		c.unusable = true
	}
}

//Write checkout the error returned from the origin Write() method.
//If the error is not nil, the connection is marked as unusable.
func (c *wrappedConn) Write(b []byte) (n int, err error) {
	if c.inPool {
		err = fmt.Errorf("write conn fail: conn is in pool")
		return
	}

	//c.Conn is certainly not nil
	n, err = c.Conn.Write(b)
	if err != nil {
		c.unusable = true
	}
	return
}

//Read works the same as Write.
func (c *wrappedConn) Read(b []byte) (n int, err error) {
	if c.inPool {
		err = fmt.Errorf("read conn fail: conn is in pool")
		return
	}

	//c.Conn is certainly not nil
	n, err = c.Conn.Read(b)
	if err != nil {
		c.unusable = true
	}

	return
}

//delayClose listens on setDelayClose and cancelDelayClose channel. When recveiving from
//setDelayClose channel, it start a ticker and wait ticker finish or receiving from
//cancelDelayClose channel
func (c *wrappedConn) delayClose() {
	var ticker *time.Ticker
	for {
		select {
		case <-c.setDelayClose:
			ticker = time.NewTicker(c.liveTime)
			select {
			case <-ticker.C:
				c.closeConnect()
			case <-c.cancelDelayClose:
				if ticker != nil {
					ticker.Stop()
				}
			}
		case <-c.cancelDelayClose:
		}
	}
}

//wrap wraps net.Conn and start a delayClose goroutine
func (p *blockingPool) wrap(conn net.Conn, livetime time.Duration) *wrappedConn {
	c := &wrappedConn{
		conn,
		p,
		true,
		true,
		livetime,
		make(chan struct{}),
		make(chan struct{}),
	}

	go c.delayClose()

	return c

}
