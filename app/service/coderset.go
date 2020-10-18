// Time : 2020/9/30 21:21
// Author : Kieran

// service
package service

import (
	"begonia2/opcode/coding"
	"reflect"
)

// coderset.go something

//TODO coderset
func newCoderSet() *coderSet {
	return &coderSet{}
}

// TODO: coderSet的实现
type coderSet struct {
}

func (s *coderSet) get(service, funName string) (fun reflectFun) {
	return
}

func (s *coderSet) store(service, funName string, fun reflectFun) {

}

type reflectFun struct {
	in  coding.Coder
	out coding.Coder
	obj interface{}
	method reflect.Method
}
