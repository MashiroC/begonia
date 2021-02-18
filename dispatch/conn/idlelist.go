package conn

import (
	"time"
)

type idleList struct {
	len         int
	front, back *poolConn
}

type poolConn struct {
	c          *defaultConn
	t          time.Time
	next, prev *poolConn
}

func (l *idleList) pushFront(pc *poolConn) {
	pc.next = l.front
	pc.prev = nil
	if l.len == 0 {
		l.back = pc
	} else {
		l.front.prev = pc
	}
	l.front = pc
	l.len++
	return
}

func (l *idleList) popFront() {
	pc := l.front
	l.len--
	if l.len == 0 {
		l.front, l.back = nil, nil
	} else {
		pc.next.prev = nil
		l.front = pc.next
	}
	pc.next, pc.prev = nil, nil
}

func (l *idleList) popBack() {
	pc := l.back
	l.len--
	if l.len == 0 {
		l.front, l.back = nil, nil
	} else {
		pc.prev.next = nil
		l.back = pc.prev
	}
	pc.next, pc.prev = nil, nil
}
