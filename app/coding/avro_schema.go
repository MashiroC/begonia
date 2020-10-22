// Time : 2020/10/20 16:04
// Author : Kieran

// coding
package coding

import (
	"reflect"
	"strconv"
	"strings"
)

// avro_schema.go something

// InSchema 根据反射类型 获得schema
func InSchema(m reflect.Method) string {

	t := m.Type

	num := t.NumIn()

	fields := make([]string, num-1)

	for i := 1; i < num; i++ {
		fieldSchema := fieldSchema("in"+strconv.FormatInt(int64(i), 10), t.In(i))
		fields[i-1] = fieldSchema
	}

	rawSchema := `
{
			"namespace":"begonia.func.` + m.Name + `",
			"type":"record",
			"name":"In",
			"fields":[
				` + strings.Join(fields, ",") + `
			]
		}`

	return rawSchema
}

func OutSchema(m reflect.Method) string {
	t := m.Type

	num := t.NumOut()

	//fmt.Println(t.Out(num-1).Kind())

	fields := make([]string, num)

	for i := 0; i < num; i++ {
		fieldSchema := fieldSchema("out"+strconv.FormatInt(int64(i+1), 10), t.Out(i))
		fields[i] = fieldSchema
	}

	rawSchema := `
{
			"namespace":"begonia.func.` + m.Name + `",
			"type":"record",
			"name":"Out",
			"fields":[
				` + strings.Join(fields, ",") + `
			]
		}`

	return rawSchema
}

func fieldSchema(name string, t reflect.Type) string {
	fType := ""
	switch t.Kind() {
	case reflect.String:
		fType = `"string"`
	case reflect.Int:
		fType = `"int"`
	case reflect.Interface:
		if t.String() == "error" {
			name = "err"
			fType = `["string","null"]`
		}
	default:
		fType += t.String()
	}
	raw := `{"name":"` + name + `","type":` + fType + `}`
	return raw
}
