package pool

import (
	"net"
	"time"
)

//WrappedConn modify the behavior of net.Conn's Write() method and Close() method
//while other methods can be accessed transparently.
type WrappedConn struct {
	net.Conn
	pool *blockingPool
	unusable bool
	start time.Time
}

//TODO
func (c *WrappedConn) Close() error {
	return c.pool.put(c)
}

//Write checkout the error returned from the origin Write() method.
//If the error is not nil, the connection is marked as unusable.
func (c *WrappedConn) Write(b []byte) (n int, err error) {
	//c.Conn is certainly not nil
	n, err = c.Conn.Write(b)
	if err != nil {
		c.unusable = true
	}
	return
}

//Read works the same as Write.
func (c *WrappedConn) Read(b []byte) (n int, err error) {
	//c.Conn is certainly not nil
	n, err = c.Conn.Read(b)
	if err != nil {
		c.unusable = true
	}
	return
}

func (p *blockingPool) wrap(conn net.Conn) *WrappedConn {
	return &WrappedConn{
		conn,
		p,
		false,
		time.Now(),
	}
}
