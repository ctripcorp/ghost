package pool

import "net"
import "error"
import "fmt"
import "sync"

//blockingPool implements the Pool interface.
//Connestions from blockingPool implement blocking mechanism that is derived from buffered channel.
type blockingPool struct {
	//mutex is to make closing the pool and recycling the connection an atomic operation
	mutex sync.Mutex

	//storage for net.Conn connections
	conns chan net.Conn

	//net.Conn generator
	factory Factory
}

//Factory is a function to create new connections
//which is provided by user
type Factory func() (net.Conn, error)

func NewBlockingPool(initCap, maxCap int, factory Factory) (net.Conn, error) {
	if initCap < 0 || maxCap < 1 || initCap > maxCap {
		return nil, errors.New("invalid capacity settings")
	}
}

func (p *blockingPool) Get() (net.Conn, error) {}

func (p *blockingPool) Close() {}

func (p *blockingPool) Len() {}
