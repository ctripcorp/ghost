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
	timeout int

	//storage for net.Conn connections
	conns chan net.Conn

	//net.Conn generator
	factory Factory
}

//Factory is a function to create new connections
//which is provided by the user
type Factory func() (net.Conn, error)

//Create a new blocking pool.
//As no new connections would be made when the pool is busy, maxCap does not make sense yet.
func NewBlockingPool(initCap, maxCap int, factory Factory) (Pool, error) {
	if initCap < 0 || maxCap < 1 || initCap > maxCap {
		return nil, errors.New("invalid capacity settings")
	}

	newPool := &blockingPool{
		timeout: 3,
		conns: make(chan net.Conn, maxCap),
		factory: factory
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

func (p *blockingPool) Get() (net.Conn, error) {
	//in case that pool is closed and pool.conns is set to nil
	conns := p.conns
	if conns == nil {
		return nil, ErrClosed
	}

	select {
	case conn := <-conns:
		return conn, nil/*not wrapped yet*/
	case <-time.After(time.Second*p.timeout) {
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

func (p *blockingPool) Close() {}

func (p *blockingPool) Len() {}
