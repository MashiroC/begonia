// Time : 2020/9/30 21:21
// Author : Kieran

// service
package service

import (
	"begonia2/opcode/coding"
	"errors"
	"reflect"
	"sync"
)

// coderset.go something

//TODO coderset
func newCoderSet() *coderSet {
	return &coderSet{
		m: make(map[string]map[string]reflectFun),
	}
}

// TODO: coderSet的实现
type coderSet struct {
	m map[string]map[string]reflectFun
	l sync.RWMutex
}

func (s *coderSet) get(service, funName string) (fun reflectFun,err error) {
	s.l.RLock()
	defer s.l.RUnlock()

	funs,ok:=s.m[service]
	if !ok{
		err = errors.New("service not found")
		return
	}

	fun,ok=funs[funName]
	if !ok{
		err = errors.New("fun not found")
	}
	return
}

func (s *coderSet) store(service, funName string, fun reflectFun) {
	s.l.Lock()
	defer s.l.Unlock()

	funs, ok := s.m[service]

	if !ok {
		s.m[service] = map[string]reflectFun{
			funName: fun,
		}
	} else {
		funs[funName] = fun
	}

}

type reflectFun struct {
	in     coding.Coder
	out    coding.Coder
	obj    interface{}
	method reflect.Method
}
