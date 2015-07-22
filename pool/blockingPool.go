package pool

import "net"
import "errors"
import "sync"
import "time"
import "fmt"

//blockingPool implements the Pool interface.
//Connestions from blockingPool offer a kind of blocking mechanism that is derived from buffered channel.
type blockingPool struct {
	//mutex is to make closing the pool and recycling the connection an atomic operation
	mutex sync.Mutex

	//timeout to Get, default to 3
	timeout time.Duration

	//storage for net.Conn connections
	conns chan net.Conn

	//net.Conn generator
	factory Factory
}

//Factory is a function to create new connections
//which is provided by the user
type Factory func() (net.Conn, error)

//Create a new blocking pool. As no new connections would be made when the pool is busy, 
//the number of connections of the pool is kept no more than initCap and maxCap does not 
//make senss but the api is reserved. The timeout to block Get() is set to 3 by default 
//concerning that it is better to be related with Get() method.
func NewBlockingPool(initCap, maxCap int, factory Factory) (Pool, error) {
	if initCap < 0 || maxCap < 1 || initCap > maxCap {
		return nil, errors.New("invalid capacity settings")
	}

	newPool := &blockingPool{
		timeout: 3,
		conns: make(chan net.Conn, maxCap),
		factory: factory,
	}

	for i := 0; i < initCap; i++ {
		conn, err := factory()
		if err != nil {
			newPool.Close()
			return nil, fmt.Errorf("error counted when calling factory: %s", err)
		}
		newPool.conns <- conn
	}
	return newPool, nil
}

//Get blocks for an available connection.
func (p *blockingPool) Get() (net.Conn, error) {
	//in case that pool is closed and pool.conns is set to nil
	conns := p.conns
	if conns == nil {
		return nil, ErrClosed
	}

	select {
	case conn := <-conns:
		return p.wrap(conn), nil/*not wrapped yet*/
	case <-time.After(time.Second*p.timeout):
		return nil, errors.New("timeout")
	}
}

//put puts the connection back to the pool. If the pool is closed, put simply close 
//any connections received and return immediately. A nil net.Conn is illegal and will be rejected.
func (p *blockingPool) put(conn net.Conn) error {
	if conn == nil {
		return errors.New("connection is nil.")
	}

	//in case that pool is closed and pool.conns is set to nil
	conns := p.conns
	if conns == nil {
		return conn.Close()
	}

	//It is impossible to block as number of connections is never more than length of channel
	conns <-conn
	return nil
}

//Create a new connection and put it into the channel.
func (p *blockingPool) compensate() error {
	conn, err := p.factory()
	if err != nil {
		//The author hopes this error never happends.
		//p.Close()
		return fmt.Errorf("error counted when calling factory: %s", err)
	}

	//in case that pool is closed and pool.conns is set to nil
	conns := p.conns
	if conns == nil {
		return nil
	}
	p.conns <-conn
	return nil
}

//Close set connection channel to nil and close all the relative connections.
//Yet not implemented.
func (p *blockingPool) Close() {}

//Len return the number of current active(in use or available) connections.
//Yet not implemented.
func (p *blockingPool) Len() int {
	return 0;
}
