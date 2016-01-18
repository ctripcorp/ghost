package pool

import (
	"errors"
	"net"
	"sync"
	"time"
)

//blockingPool implements the Pool interface.
//Connestions from blockingPool offer a kind of blocking mechanism that is derived from buffered channel.
type blockingPool struct {
	//mutex is to make closing the pool and recycling the connection an atomic operation
	mutex sync.Mutex

	//check unalive connections period, default to 5s
	checkPeriod time.Duration

	//timeout to Get, default to 3s
	timeout time.Duration

	//storage for net.Conn connections
	conns *Deque

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
func NewBlockingPool(initCap, maxCap int, livetime time.Duration, factory Factory) (Pool, error) {
	if initCap < 0 || maxCap < 1 || initCap > maxCap {
		return nil, errors.New("invalid capacity settings")
	}

	newPool := &blockingPool{
		checkPeriod: 5 * time.Second,
		timeout:     3 * time.Second,
		conns:       NewCappedDeque(maxCap),
		factory:     factory,
	}

	for i := 0; i < initCap; i++ {
		newPool.conns.Append(newPool.wrap(nil, livetime))
	}

	go newPool.checkUnactiveConn()
	return newPool, nil
}

func (p *blockingPool) checkUnactiveConn() {
	for {
		time.Sleep(p.checkPeriod)
		result := p.conns.Walk(func(item interface{}) interface{} {
			c := item.(*wrappedConn)
			return c.getInactiveNetConn()
		})
		for _, v := range result {
			if v != nil {
				conn := v.(net.Conn)
				conn.Close()
			}
		}
	}
}

//Get blocks for an available connection.
func (p *blockingPool) Get() (net.Conn, error) {
	//in case that pool is closed or pool.conns is set to nil
	conns := p.conns
	if conns == nil {
		return nil, ErrClosed
	}
	// new connection is popped from tail of Deque
	item, err := p.conns.Pop(p.timeout)
	if err != nil {
		return nil, err
	}
	conn := item.(*wrappedConn)
	conn.inPool = false
	conn.closeInactiveNetConn()

	if conn.Conn == nil {
		if c, err := p.factory(); err != nil {
			//conn.Conn is possibly nil coz factory() may fail, in which case conn is immediately
			//put back to the pool
			conn.Close()
			return nil, err
		} else {
			conn.updateNetConn(c)
		}
	}

	return conn, nil
}

//put puts the connection back to the pool. If the pool is closed, put simply close
//any connections received and return immediately. A nil net.Conn is illegal and will be rejected.
func (p *blockingPool) put(conn *wrappedConn) error {
	//in case that pool is closed and pool.conns is set to nil
	conns := p.conns
	if conns == nil {
		conn.closeNetConn()
		return ErrClosed
	}

	if conn.unusable {
		//if conn is marked unusable, close inner net.Conn, and append it to head of Deque
		conn.closeNetConn()
		p.conns.Prepend(conn)

	} else {
		//else append conn to tail of Deque
		p.conns.Append(conn)
	}
	conn.inPool = true
	return nil
}

//TODO
//Close set connection channel to nil and close all the relative connections.
//Yet not implemented.
func (p *blockingPool) Close() {}

//TODO
//Len return the number of current active(in use or available) connections.
func (p *blockingPool) Len() int {
	return 0
}
