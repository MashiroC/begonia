package coding

import (
	"reflect"
)

type ReSharpFunc func(in interface{}) interface{}

func toAvroSchemaField(t reflect.Type) string {
	return t.String()
}

// FunInfo 函数信息
type FunInfo struct {
	Name      string `avro:"name"`
	Mode      string `avro:"mode"`
	InSchema  string `avro:"inSchema"`
	OutSchema string `avro:"outSchema"`
	ParamTyp  []string
	ResultTyp []string
}

// Parse 将一个结构体的函数信息解析
func Parse(mode string, in interface{}) (fi []FunInfo, methods []reflect.Method, reSharps [][]ReSharpFunc) {
	//TODO:先简单写一下 后面再支持更多类型
	if mode != "avro" {
		panic("parse mode error")
	}

	t := reflect.TypeOf(in)

	fi = make([]FunInfo, t.NumMethod())
	methods = make([]reflect.Method, t.NumMethod())
	reSharps = make([][]ReSharpFunc, t.NumMethod())

	for i := 0; i < t.NumMethod(); i++ {

		m := t.Method(i)
		methods[i] = m

		inS := inReflectSchema(m)
		outS := outReflectSchema(m)

		reSharps[i] = parseReSharpFunc(m)

		fi[i] = FunInfo{
			Name:      m.Name,
			Mode:      mode,
			InSchema:  inS,
			OutSchema: outS,
		}
	}

	return
}
