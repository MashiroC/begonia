package service

import (
	"errors"
	"github.com/MashiroC/begonia/app/coding"
	"reflect"
	"sync"
)

// service_store.go 存放远程函数的数据结构

// newServiceStore 创建一个新的实例
func newServiceStore() *serviceStore {
	return &serviceStore{
		m: make(map[string]map[string]reflectFun),
	}
}

// serviceStore 实现的数据结构
type serviceStore struct {
	m map[string]map[string]reflectFun // 实际存储的map 两层map service - fun - reflectFun
	l sync.RWMutex                     // 线程安全的锁
}

// get 获得一个远程函数的信息
func (s *serviceStore) get(service, funName string) (fun reflectFun, err error) {

	s.l.RLock()
	defer s.l.RUnlock()

	funs, ok := s.m[service]
	if !ok {
		err = errors.New("service not found")
		return
	}

	fun, ok = funs[funName]
	if !ok {
		err = errors.New("fun not found")
	}
	return
}

// store 存储一个远程函数的信息
func (s *serviceStore) store(service, funName string, fun reflectFun) {
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

// reflectFun 存储的使用反射实现的远程函数的信息
type reflectFun struct {
	in      coding.Coder
	out     coding.Coder
	reSharp []coding.ReSharpFunc
	obj     interface{}
	method  reflect.Method
}
