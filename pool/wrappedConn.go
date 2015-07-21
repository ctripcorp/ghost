package pool

import "net"

//WrappedConn modify the behavior of net.Conn's Write() method and Close() method
//while other methods can be accessed transparently.
type WrappedConn struct {
	net.Conn
	pool *blockingPool
	unusable bool
}

func (c WrappedConn) Close() error {
	if c.unusable {
		c.pool.compensate()
		if c.Conn != nil {
			return c.Conn.Close()
		}
		return nil
	}
	return c.pool.put(c.Conn)
}
