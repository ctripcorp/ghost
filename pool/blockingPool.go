package pool

import "net"
import "errors"
import "sync"
import "time"

//blockingPool implements the Pool interface.
//Connestions from blockingPool offer a kind of blocking mechanism that is derived from buffered channel.
type blockingPool struct {
	//mutex is to make closing the pool and recycling the connection an atomic operation
	mutex sync.Mutex

	//timeout to Get, default to 3
	timeout time.Duration

	//storage for net.Conn connections
	conns chan *WrappedConn

	//net.Conn generator
	factory Factory
}

//Factory is a function to create new connections
//which is provided by the user
type Factory func() (net.Conn, error)

//Create a new blocking pool. As no new connections would be made when the pool is busy, 
//the number of connections of the pool is kept no more than initCap and maxCap does not 
//make sense but the api is reserved. The timeout to block Get() is set to 3 by default 
//concerning that it is better to be related with Get() method.
func NewBlockingPool(initCap, maxCap int, factory Factory) (Pool, error) {
	if initCap < 0 || maxCap < 1 || initCap > maxCap {
		return nil, errors.New("invalid capacity settings")
	}

	newPool := &blockingPool{
		timeout: 3,
		conns: make(chan *WrappedConn, maxCap),
		factory: factory,
	}

	for i := 0; i < initCap; i++ {
		conn, _ := factory()
		//It doesn't matter if conn it nil, nil is checked whenever Get() is called, 
		//so error counted from factory() is ignored.
		newPool.conns <- newPool.wrap(conn)
	}
	return newPool, nil
}

//Get blocks for an available connection.
func (p *blockingPool) Get() (net.Conn, error) {
	//in case that pool is closed or pool.conns is set to nil
	conns := p.conns
	if conns == nil {
		return nil, ErrClosed
	}

	select {
	case conn := <-conns:
		if conn.Conn == nil {
			var err error
			conn.Conn, err = p.factory()
			if err != nil {
				p.put(conn)
				return nil, err
			}
		}
		conn.unusable = false
		return conn, nil
	case <-time.After(time.Second*p.timeout):
		return nil, ErrTimeout
	}
}

//put puts the connection back to the pool. If the pool is closed, put simply close 
//any connections received and return immediately. A nil net.Conn is illegal and will be rejected.
func (p *blockingPool) put(conn *WrappedConn) error {
	//in case that pool is closed and pool.conns is set to nil
	conns := p.conns
	if conns == nil {
		//conn.Conn is possibly nil coz factory() may fail, in which case conn is immediately 
		//put back to the pool
		if conn.Conn != nil {
			conn.Conn.Close()
			conn.Conn = nil
		}
		return ErrClosed
	}

	//if conn is marked unusable, underlying net.Conn is set to nil
	if conn.unusable {
		if conn.Conn != nil {
			conn.Conn.Close()
			conn.Conn = nil
		}
	}

	//It is impossible to block as number of connections is never more than length of channel
	conns <-conn
	return nil
}

//TODO
//Close set connection channel to nil and close all the relative connections.
//Yet not implemented.
func (p *blockingPool) Close() {}

//TODO
//Len return the number of current active(in use or available) connections.
func (p *blockingPool) Len() int {
	return 0;
}

