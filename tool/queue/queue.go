package queue

import (
	"container/list"
	"sync"
)

type Queue struct {
	l    *list.List
	Lock sync.Mutex
}

func New() *Queue {
	q := &Queue{l: list.New()}
	return q
}

func (q *Queue) Push(val interface{}) {
	q.l.PushFront(val)
}

func (q *Queue) PopBack() (val interface{}) {
	e := q.l.Back()
	if e != nil {
		val = q.l.Remove(e)
	}
	return val
}

func (q *Queue) PopFront() (val interface{}) {
	e := q.l.Front()
	if e != nil {
		val = q.l.Remove(e)
	}
	return val
}

func (q *Queue) Back() (val interface{}) {
	e := q.l.Back()
	return e.Value
}

func (q *Queue) Front() (val interface{}) {
	e := q.l.Front()
	return e.Value
}

func (q *Queue) IsEmpty() bool {
	e := q.l.Back()
	return e == nil
}

func (q *Queue) Len() int {
	l := q.l.Len()
	return l
}
