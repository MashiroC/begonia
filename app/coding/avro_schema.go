package coding

import (
	"reflect"
	"strconv"
	"strings"
)

type parseMode int

const (
	_invalid parseMode = iota
	modeNormal
	modeSlice
)

// avro_schema.go something

// inReflectSchema 根据反射类型 获得schema
func inReflectSchema(m reflect.Method) string {
	namespace := "begonia.func." + m.Name
	name := "In"

	t := m.Type
	num := t.NumIn()

	typ := make([]reflect.Type, 0, num-1)
	for i := 1; i < num; i++ {
		typ = append(typ, t.In(i))
	}

	res := makeSchema(namespace, name, typ)
	return res
}

// outReflectSchema 根据反射 获得出参schema
func outReflectSchema(m reflect.Method) string {
	namespace := "begonia.func." + m.Name
	name := "Out"

	t := m.Type
	num := t.NumOut()

	typ := make([]reflect.Type, 0, num)
	for i := 0; i < num; i++ {
		typ = append(typ, t.Out(i))
	}

	return makeSchema(namespace, name, typ)
}

func makeSchema(namespace, name string, typ []reflect.Type) string {
	fields := make([]string, len(typ))

	for i := 0; i < len(typ); i++ {
		fields[i] = fieldSchema("F"+strconv.FormatInt(int64(i+1), 10), typ[i])
	}

	rawSchema := `
{
			"namespace":"` + namespace + `",
			"type":"record",
			"name":"` + name + `",
			"fields":[
				` + strings.Join(fields, ",") + `
			]
		}`

	return rawSchema
}

// fieldSchema 解析对应的类型到avro schema，如果遇到指针，会取指针后再对类型做解析
// 目前支持的类型：
// int, int8 ~ int32 -> int
// int64             -> long
// float32           -> float
// float64           -> double
// bool              -> boolean
// string            -> string
// error             -> ["string","null"]
// slice             -> array
// struct            -> record
// []uint8           -> bytes
// map[string]kind   -> map
//
// ps:
// - 目前avro类型不支持无符号整数，uint uint8~uint64全部不支持，唯一的例外是[]uint8会被解析成bytes
// - map的value，中只要上述类型支持的，都可以支持，不支持interface{}
// - 结构体支持嵌套、内嵌
// - 不支持array，请使用slice
func fieldSchema(name string, t reflect.Type) (schema string) {
	fType, isErr := fieldKind(modeNormal, t)
	if isErr {
		name = "err"
	}
	schema = `{"name":"` + name + `","type":` + fType + "}\n"

	return
}

func fieldKind(mode parseMode, t reflect.Type) (fType string, isErr bool) {
	switch t.Kind() {
	case reflect.String:
		fType = `"string"`
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int:
		fType = `"int"`
	case reflect.Int64:
		fType = `"long"`
	case reflect.Float32:
		fType = `"float"`
	case reflect.Float64:
		fType = `"double"`
	case reflect.Bool:
		fType = `"boolean"`
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			fType = `"bytes"`
		} else {
			childKind, _ := fieldKind(modeSlice, t.Elem())

			fType = `{
				"type": "array",
				"items": ` + childKind + `
			}`

		}
	case reflect.Interface:
		if mode == modeSlice {
			panic("slice not supported interface")
		}
		if t.String() == "error" {
			fType = `["string","null"]`
			isErr = true
		} else {
			panic("avro parse not supported")
		}
	case reflect.Ptr:
		if mode == modeSlice {
			panic("github.com/MashiroC/begonia not supported ptr")
		}
		fType, isErr = fieldKind(modeNormal, t.Elem())
	case reflect.Struct:
		if mode == modeSlice {
			panic("slice not supported struct")
		}
		n := t.NumField()
		fields := make([]string, n)

		for i := 0; i < n; i++ {
			field := t.Field(i)
			fields[i] = fieldSchema(field.Name, field.Type)
		}

		fType = `{
				"type": "record",
				"name": "` + t.Name() + `",
				"fields":[` + strings.Join(fields, ",") + `
				]
			}`
	case reflect.Map:
		if mode == modeSlice {
			panic("slice not supported map")
		}

		child, _ := fieldKind(modeNormal, t.Elem())
		fType = `{"type":"map","values":` + child + `}`
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		panic("avro not supported uint")
	default:
		fType += t.String()
	}
	return
}
