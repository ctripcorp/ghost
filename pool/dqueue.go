package pool

import (
	"container/list"
	"sync"
	"time"
)

type Deque struct {
	queue *deque
	gate  chan struct{}
}

func NewCappedDeque(maxCap int) *Deque {
	return &Deque{
		newCappedDeque(maxCap),
		make(chan struct{}, maxCap),
	}
}

// Pop pops an item from the tail of double-ended queue in timeout duraton
func (s *Deque) Pop(timeout time.Duration) (interface{}, error) {
	select {
	case <-s.gate:
		elem := s.queue.pop()
		return elem, nil
	case <-time.After(timeout):
		return nil, ErrTimeout
	}

}

// Append appends item to the tail of double-ended queue
func (s *Deque) Append(item interface{}) {
	s.queue.append(item)
	s.gate <- struct{}{}
}

// Prepend appends item to the head of double-end queue
func (s *Deque) Prepend(item interface{}) {
	s.queue.prepend(item)
	s.gate <- struct{}{}
}

type WalkFunc func(itme interface{}) interface{}

// Walk
func (s *Deque) Walk(walkFn WalkFunc) []interface{} {
	return s.queue.walk(walkFn)
}

// Deque is a head-tail linked list data structure implementation.
// It is based on a doubly linked list container, so that every
// operations time complexity is O(1).
//
// every operations over an instiated Deque are synchronized and
// safe for concurrent usage.
type deque struct {
	sync.RWMutex
	container *list.List
	capacity  int
}

// NewCappedDeque creates a Deque with the specified capacity limit.
func newCappedDeque(capacity int) *deque {
	return &deque{
		container: list.New(),
		capacity:  capacity,
	}
}

// Append inserts element at the back of the Deque in a O(1) time complexity,
// returning true if successful or false if the deque is at capacity.
func (s *deque) append(item interface{}) bool {
	s.Lock()
	defer s.Unlock()

	if s.capacity < 0 || s.container.Len() < s.capacity {
		s.container.PushBack(item)
		return true
	}

	return false
}

// prepend inserts element at the Deques front in a O(1) time complexity,
// returning true if successful or false if the deque is at capacity.
func (s *deque) prepend(item interface{}) bool {
	s.Lock()
	defer s.Unlock()

	if s.capacity < 0 || s.container.Len() < s.capacity {
		s.container.PushFront(item)
		return true
	}

	return false
}

// pop removes the last element of the deque in a O(1) time complexity
func (s *deque) pop() interface{} {
	s.Lock()
	defer s.Unlock()

	var item interface{} = nil
	var lastContainerItem *list.Element = nil

	lastContainerItem = s.container.Back()
	if lastContainerItem != nil {
		item = s.container.Remove(lastContainerItem)
	}

	return item
}

// walk traverses the deque from front, calling walkFn for each element,
// and return polymerized result
func (s *deque) walk(walkFn WalkFunc) []interface{} {
	s.Lock()
	defer s.Unlock()

	rst := make([]interface{}, 0, s.capacity)
	item := s.container.Front()

	for {
		if item != nil {
			rst = append(rst, walkFn(item.Value))
			item = item.Next()
		} else {
			break
		}
	}

	return rst
}
