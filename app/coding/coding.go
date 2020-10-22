// Time : 2020/9/26 19:47
// Author : Kieran

// coding
package coding

import (
	"reflect"
)

// coding.go something

type Coder interface {
	Encode(data interface{}) ([]byte, error)
	Decode([]byte) (data interface{}, err error)
	DecodeIn([]byte, interface{}) error
}

type FunInfo struct {
	Name      string `avro:"name"`
	Mode      string `avro:"mode"`
	InSchema  string `avro:"inSchema"`
	OutSchema string `avro:"outSchema"`
}

func Parse(mode string, in interface{}) (fi []FunInfo, methods []reflect.Method) {
	//TODO:先简单写一下 后面再支持更多类型
	if mode != "avro" {
		panic("parse mode error")
	}

	t := reflect.TypeOf(in)

	fi = make([]FunInfo, t.NumMethod())
	methods = make([]reflect.Method, t.NumMethod())

	for i := 0; i < t.NumMethod(); i++ {

		m := t.Method(i)
		methods[i] = m

		inS := InSchema(m)
		outS := OutSchema(m)

		fi[i] = FunInfo{
			Name:      m.Name,
			Mode:      mode,
			InSchema:  inS,
			OutSchema: outS,
		}
	}

	return
}
