package coding

import "reflect"

func toAvroSchemaField(t reflect.Type) string {
	return t.String()
}

// FunInfo 函数信息
type FunInfo struct {
	Name      string `avro:"name"`
	Mode      string `avro:"mode"`
	InSchema  string `avro:"inSchema"`
	OutSchema string `avro:"outSchema"`
}

// Parse 将一个结构体的函数信息解析
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