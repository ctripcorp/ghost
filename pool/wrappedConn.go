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
	pool       *blockingPool
	unusable   bool
	inPool     bool
	lastAccess time.Time
	liveTime   time.Duration
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

// updateNetConn
func (c *wrappedConn) updateNetConn(newNetConn net.Conn) {
	c.Conn = newNetConn
	c.unusable = false
	c.lastAccess = time.Now()
}

// setNetConnNil set inner net conn to nil and return old net conn
func (c *wrappedConn) setNetConnNil() (oldNetConn net.Conn) {
	oldNetConn = c.Conn
	c.Conn = nil
	c.unusable = true
	c.lastAccess = time.Now()
	return
}

// closeNetConn close inner net conn and set it to nil
func (c *wrappedConn) closeNetConn() {
	if conn := c.setNetConnNil(); conn != nil {
		conn.Close()
	}
}

// getInactiveNetConn
func (c *wrappedConn) getInactiveNetConn() (conn net.Conn) {
	if time.Since(c.lastAccess) > c.liveTime && c.Conn != nil {
		conn = c.setNetConnNil()
	}
	return
}

// closeInactiveNetConn
func (c *wrappedConn) closeInactiveNetConn() {
	if conn := c.getInactiveNetConn(); conn != nil {
		conn.Close()
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
	} else {
		c.lastAccess = time.Now()
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
	} else {
		c.lastAccess = time.Now()
	}

	return
}

//wrap wraps net.Conn and start a delayClose goroutine
func (p *blockingPool) wrap(conn net.Conn, livetime time.Duration) *wrappedConn {
	return &wrappedConn{
		conn,
		p,
		true,
		true,
		time.Now(),
		livetime,
	}
}
