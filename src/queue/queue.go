package queue

import (
	"errors"
	"sync"
)

// region => constants

const (
	MinCapacity = 1
	MaxCapacity = 2000000
)

// endregion

// region => base queue interface

type FifoQueue interface {
	Insert(item any) error
	Remove() (any, error)
	Size() int
}

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
