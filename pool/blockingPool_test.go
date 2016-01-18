package pool

import (
	"net"
	"testing"
	"time"
)

func makeConn() (net.Conn, error) {
	conn, err := net.Dial("tcp", "10.2.6.233:22122")
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func Test_Get(t *testing.T) {
	pool, _ := NewBlockingPool(5, 5, 10*time.Second, makeConn)
	for i := 0; i < 100; i++ {
		if conn, err := pool.Get(); err != nil {
			t.Error(err)
		} else {
			conn.Close()
		}
	}
}

func Test_Write(t *testing.T) {
	pool, _ := NewBlockingPool(5, 5, 10*time.Second, makeConn)
	conn, err := pool.Get()
	if err != nil {
		t.Error(err)
	}

	blob := make([]byte, 1024)
	if _, err = conn.Write(blob); err != nil {
		t.Error(err)
	}

	conn.Close()
	conn, err = pool.Get()
	if err != nil {
		t.Error(err)
	}
	conn.Write(blob)
	conn.Close()

}

func Test_Read(t *testing.T) {
	pool, _ := NewBlockingPool(5, 5, 10*time.Second, makeConn)
	conn, err := pool.Get()
	if err != nil {
		t.Error(err)
	}
	blob := make([]byte, 1024)
	conn.Close()
	if _, err = conn.Read(blob); err.Error() != "read conn fail: conn is in pool" {
		t.Error(err)
	}
}
