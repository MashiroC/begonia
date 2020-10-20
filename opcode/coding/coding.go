// Time : 2020/9/26 19:47
// Author : Kieran

// coding
package coding

import (
	"fmt"
	"reflect"
)

// coding.go something

type Coder interface {
	Encode(data interface{}) ([]byte, error)
	Decode([]byte) (data interface{}, err error)
	DecodeIn([]byte, interface{}) error
}

type FunInfo struct {
	Name      string
	Mode      string
	InSchema  string
	OutSchema string
}

func Parse(mode string, in interface{}) (c Coder,fi []FunInfo) {
	if mode != "avro" {
		panic("parse mode error")
	}

	t:=reflect.TypeOf(in)
	v:=reflect.ValueOf(in)

	fi=make([]FunInfo,t.NumMethod())
	for i:=0;i<t.NumMethod();i++{
		fmt.Println(t.Method(i))
	}
	fmt.Println(t,v)
	return
}
