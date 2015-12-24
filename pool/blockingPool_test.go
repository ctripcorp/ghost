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
	var err error
	conns := make([]net.Conn, 5)
	for i := 0; i < 5; i++ {
		conns[i], err = pool.Get()
		if err != nil {
			t.Error(err)
		}
	}
	for i := 0; i < 5; i++ {
		conns[i].Close()
	}
}

func Test_Write(t *testing.T) {
	pool, _ := NewBlockingPool(5, 5, 10*time.Second, makeConn)
	conn, err := pool.Get()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	conn.Write([]byte("hello world"))
	conn.Close()
	if _, err = conn.Write([]byte("hello world")); err.Error() != "write conn fail: conn is in pool" {
		t.Error(err)
	}
}

func Test_Read(t *testing.T) {
	pool, _ := NewBlockingPool(5, 5, 10*time.Second, makeConn)
	conn, err := pool.Get()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	blob := make([]byte, 1024)
	conn.Close()
	if _, err = conn.Read(blob); err.Error() != "read conn fail: conn is in pool" {
		t.Error(err)
	}
}
