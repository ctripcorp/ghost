package pool

import "net"

//WrappedConn modify the behavior of net.Conn's Write() method and Close() method
//while other methods can be accessed transparently.
type WrappedConn struct {
	net.Conn
	pool *blockingPool
	unusable bool
}

//Close put the connection back to the pool.
//If the connection is marked unusable, Close close the connection and call 
//blockingPool.compensate which create a new connection and put it instead.
func (c *WrappedConn) Close() error {
	if c.unusable {
		c.pool.compensate()
		if c.Conn != nil {
			return c.Conn.Close()
		}
		return nil
	}
	return c.pool.put(c.Conn)
}

//Write checkout the error returned from the origin Write() method.
//If the error is not nil, the connection is marked as unusable.
func (c *WrappedConn) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	if err != nil {
		c.unusable = true
	}
	return
}

//Read works the same as Write.
func (c *WrappedConn) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	if err != nil {
		c.unusable = true
	}
	return
}

func (p *blockingPool) wrap(conn net.Conn) net.Conn {
	return &WrappedConn{
		conn,
		p,
		false,
	}
}
