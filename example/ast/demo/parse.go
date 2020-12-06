package demo

import (
	"go/ast"
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

type AstParam struct {
	Typ   string
	Name  string
	Child []*AstParam
}

func MakeSchema(funcName, objName string, fl *ast.FieldList) (schema string, typs []string) {

	var fields []string

	if fl != nil {
		f := fl.List
		fields = make([]string, 0, len(f))
		typs = make([]string, 0, len(f))

		count := 1

		for i := 0; i < len(f); i++ {
			if len(f[i].Names) == 0 {
				fields = append(fields, fieldSchema("F"+strconv.FormatInt(int64(count), 10), f[i]))
			} else {
				for range f[i].Names {
					fields = append(fields, fieldSchema("F"+strconv.FormatInt(int64(count), 10), f[i]))
					typs = append(typs, getTyp(f[i].Type))
					count++
				}
			}

		}
	}

	schema = `
{
			"namespace":"begonia.func.` + funcName + `",
			"type":"record",
			"name":"` + objName + `",
			"fields":[
				` + strings.Join(fields, ",") + `
			]
		}`

	return
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
func fieldSchema(name string, f *ast.Field) (schema string) {
	//if len(f.Names) == 0 {
	//	name = "hello"
	//} else {
	//	name = f.Names[0].Name
	//}
	fType, isErr := fieldKind(name, f.Type)
	if isErr {
		name = "err"
	}
	schema = `{"name":"` + name + `","type":` + fType + "}\n"

	return
}

func fieldKind(name string, expr ast.Expr) (fType string, isErr bool) {

	switch in := expr.(type) {
	case *ast.StarExpr:
		panic("not support pointer")
	case *ast.ArrayType:
		// slice
		childKind, _ := fieldKind("", in.Elt)
		if childKind == "byte" || childKind == "uint8" {
			fType = `"bytes"`
		} else {
			fType = `{
				"type": "array",
				"items": ` + childKind + `
			}`
		}

	case *ast.MapType:
		// map
		k, _ := fieldKind("", in.Key)
		if k != `"string"` {
			panic("map key must string but " + k)
		}
		v, _ := fieldKind("", in.Value)
		fType = `{"type":"map","values":` + v + `}`
	case *ast.Ident:
		// 其他类型
		if in.Name == "int64" {
			fType = `"long"`
			return
		}

		if in.Name == "uint8" {
			fType = "uint8"
			return
		}

		if strings.HasPrefix(in.Name, "int") {
			fType = `"int"`
			return
		}

		// 不支持uint
		if strings.HasPrefix(in.Name, "uint") || in.Name == "interface{}" {
			isErr = true
			return
		}

		switch in.Name {
		case "float32":
			fType = `"float"`
		case "float64":
			fType = `"double"`
		case "bool":
			fType = `"boolean"`
		case "string":
			fType = `"string"`
		case "error":
			fType = `["string","null"]`
		case "byte":
			fType = `byte`
		default:
			// 结构体

			f := in.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List
			fields := make([]string, len(f))

			//log.Fatalln(in.Name,f.List)
			for i := 0; i < len(f); i++ {
				var fieldName string
				if len(f[i].Names) == 0 {
					// 内嵌结构体
					fieldName = "asd"
				} else {
					fieldName = f[i].Names[0].Name
				}
				fields[i] = fieldSchema(fieldName, f[i])
			}
			//
			fType = `{
				"type": "record",
				"name": "` + in.Name + `",
				"fields":[` + strings.Join(fields, ",") + `
				]
			}`
		}

	default:
		panic(reflect.TypeOf(expr))
	}

	return
}

func getTyp(expr ast.Expr) string {
	switch in := expr.(type) {
	case *ast.StarExpr:
		panic("not support pointer")
	case *ast.ArrayType:
		// slice
		return "[]" + getTyp(in.Elt)

	case *ast.MapType:
		// map
		return "map[" + getTyp(in.Key) + "]" + getTyp(in.Value)
	case *ast.Ident:
		// 其他类型
		return in.Name
	default:
		panic(reflect.TypeOf(in))
	}

}
