// Time : 2020/9/28 19:48
// Author : Kieran

// containers
package containers

import (
	"begonia2/app/coding"
	"sync"
)

// coderSet.go something

func NewCoderSet() *CoderSet {
	return &CoderSet{m: make(map[string]coding.Coder)}
}

type CoderSet struct {
	l sync.RWMutex
	m map[string]coding.Coder
}

func (c *CoderSet) Set(k string, v coding.Coder) {
	c.l.Lock()
	defer c.l.Unlock()

	c.m[k] = v
}

func (c *CoderSet) Get(k string) coding.Coder {
	c.l.RLock()
	defer c.l.RUnlock()

	return c.m[k]
}
