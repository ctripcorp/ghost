package pool

import (
	"github.com/oleiade/lane"
	"time"
)

type wrappedDeque struct {
	Queue *lane.Deque
	gate  chan struct{}
}

func newWrappedDeque(maxCap int) *wrappedDeque {
	q := lane.NewCappedDeque(maxCap)
	e := make(chan struct{}, maxCap)
	return &wrappedDeque{q, e}
}

// popTail pops an item from the tail of double-ended queue in timeout duraton
func (q *wrappedDeque) popTail(timeout time.Duration) (interface{}, error) {
	select {
	case <-q.gate:
		elem := q.Queue.Pop()
		return elem, nil
	case <-time.After(timeout):
		return nil, ErrTimeout
	}

}

// appendTail appends item to the tail of double-ended queue
func (q *wrappedDeque) appendTail(item interface{}) {
	q.Queue.Append(item)
	q.gate <- struct{}{}
}

// appendHead appends item to the head of double-end queue
func (q *wrappedDeque) appendHead(item interface{}) {
	q.Queue.Prepend(item)
	q.gate <- struct{}{}
}
