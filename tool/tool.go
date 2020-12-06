package tool

import (
	"errors"
	"sync"
	"time"
)

func Timeout(f func()) error {
	var pos bool
	var isTimeout bool
	var l, l2 sync.Mutex
	l.Lock()

	time.AfterFunc(3*time.Second, func() {
		l2.Lock()
		if !pos {
			pos = true
			isTimeout = true
			l.Unlock()
		}
		l2.Unlock()
	})

	go func() {
		f()
		l2.Lock()
		if !pos {
			pos = true
			l.Unlock()
		}
		l2.Unlock()
	}()
	l.Lock()
	l.Unlock()
	if isTimeout {
		return errors.New("timeout")
	}
	return nil
}
