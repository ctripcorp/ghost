package pool

import "net"
import "errors"

//ErrClosed occurs when Len() or Get() is called after Close() is
var (
	ErrClosed = errors.New("pool is closed")
	ErrTimeout = errors.New("timeout")
)

//Pool interface describes a connection pool.
type Pool interface {
	//Get returns an available(new or reused) connection from the pool.
	//It blocks when no connection is available.
	//Closing a connection puts the connection back to the pool.
	//Duplicated close of a connection is endurable as the second close is ignored.
	Get() (net.Conn, error)

	//Close closes the pool and all its connections.
	//Call Get() when the pool is closed is counted as an error.
	Close()

	//Len returns the current number of connections in the pool,
	//including those that are in use or free.
	Len() int
}
