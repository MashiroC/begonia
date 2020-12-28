package coding

import (
	"github.com/MashiroC/begonia/tool/qarr"
	"reflect"
)

type ReSharpFunc func(in interface{}) interface{}

func toAvroSchemaField(t reflect.Type) string {
	return t.String()
}

// FunInfo 函数信息
type FunInfo struct {
	Name       string
	InSchema   string
	OutSchema  string
	ParamTyp   []string
	ResultTyp  []string
	HasContext bool
}

// Parse 将一个结构体的函数信息解析
func Parse(mode string, in interface{}, registerFunc []string) (fi []FunInfo, methods []reflect.Method, reSharps [][]ReSharpFunc) {
	//TODO:先简单写一下 后面再支持更多类型
	if mode != "avro" {
		panic("parse mode error")
	}

	t := reflect.TypeOf(in)

	fi = make([]FunInfo, 0, 2)
	methods = make([]reflect.Method, 0, 2)
	reSharps = make([][]ReSharpFunc, 0, 2)

	for i := 0; i < t.NumMethod(); i++ {

		m := t.Method(i)

		if registerFunc != nil && len(registerFunc) != 0 && !qarr.StringsIn(registerFunc, m.Name) {
			continue
		}

		methods = append(methods, m)

		inS, hasContext := inReflectSchema(m)
		outS := outReflectSchema(m)

		reSharps = append(reSharps, parseReSharpFunc(m))

		fi = append(fi, FunInfo{
			Name:       m.Name,
			InSchema:   inS,
			OutSchema:  outS,
			HasContext: hasContext,
		})
	}

	return
}
