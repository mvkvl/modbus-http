package queue

import (
	"errors"
	"math"
	"sync"
)

// region => constants

const (
	MinCapacity    = 1
	MaxCapacity    = 2000000
	CapacityLimit  = 0.7
	CapacityFactor = 2
)

// endregion

// region => base queue interface

type FifoQueue interface {
	Insert(item any) error
	Remove() (any, error)
	Size() int
}

// endregion

// region => channel-based queue

// region -> type

type QueueChan struct {
	capacity int
	q        chan any
}

// endregion

// region -> functions

func (q *QueueChan) Insert(item any) error {
	if err := q.expand(); err != nil {
		return err
	}
	if len(q.q) < q.capacity {
		q.q <- item
		return nil
	}
	return errors.New("QueueChan is full")
}

func (q *QueueChan) Remove() (any, error) {
	if err := q.shrink(); err != nil {
		return "", err
	}
	if len(q.q) > 0 {
		item := <-q.q
		return item, nil
	}
	return "", errors.New("QueueChan is empty")
}

func (q *QueueChan) Size() int {
	return len(q.q)
}

func (q *QueueChan) Capacity() int {
	return q.capacity
}

func (q *QueueChan) expand() error {
	if len(q.q) >= int(float32(q.capacity)*CapacityLimit) {
		if len(q.q) < MaxCapacity {
			defer close(q.q)
			newCapacity := int(math.Min(float64(q.capacity*CapacityFactor), MaxCapacity))
			newChannel := make(chan any, newCapacity)
			for v := range q.q {
				newChannel <- v
			}
			q.capacity = newCapacity
			q.q = newChannel
			return nil
		} else {
			//return errors.New("QueueChan maximum size is reached")
		}
	}
	return nil
}

func (q *QueueChan) shrink() error {
	if len(q.q) < int(float32(q.capacity)*(1-CapacityLimit)) {
		if len(q.q) > MinCapacity {
			newCapacity := int(math.Max(float64(q.capacity/CapacityFactor), MinCapacity))
			newChannel := make(chan any, newCapacity)
			for v := range q.q {
				newChannel <- v
			}
			q.capacity = newCapacity
			q.q = newChannel
			return nil
		} else {
			//return errors.New("QueueChan minimum size is reached")
		}
	}
	return nil
}

// endregion

// region -> constructor

func CreateChanQueueOfSize(capacity int) FifoQueue {
	return createChanQueue(capacity)
}

func CreateChanQueue() FifoQueue {
	return createChanQueue(MinCapacity)
}

func createChanQueue(capacity int) FifoQueue {
	return &QueueChan{
		capacity: capacity,
		q:        make(chan any, capacity),
	}
}

// endregion

// endregion

// region => array/mutex-based queue

// region -> type

type Queue struct {
	m sync.Mutex
	q []any
}

// endregion

// region -> public API

func (q *Queue) Insert(item any) error {
	q.m.Lock()
	defer q.m.Unlock()
	if len(q.q) < MaxCapacity {
		q.q = append(q.q, item)
		return nil
	}
	return errors.New("queue is full")
}

func (q *Queue) Remove() (any, error) {
	q.m.Lock()
	defer q.m.Unlock()
	if len(q.q) > 0 {
		item := q.q[0]
		q.q = q.q[1:]
		return item, nil
	}
	return nil, errors.New("queue is empty")
}

func (q *Queue) Size() int {
	return len(q.q)
}

// endregion

// region -> constructor

func CreateQueue() FifoQueue {
	return &Queue{
		q: make([]any, 0, MinCapacity),
	}
}

// endregion

// endregion
