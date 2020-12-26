package main

import (
	"go/ast"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type AstParam struct {
	Typ   string
	Name  string
	Child []*AstParam
}

func MakeSchema(funcName, objName string, fl *ast.FieldList) (schema string, typs []string, hasContext bool) {

	var fields []string

	if fl != nil {
		f := fl.List
		fields = make([]string, 0, len(f))
		typs = make([]string, 0, len(f))
		count := 1

		if len(f) != 0 {

			start := 0
			// 判断context
			if se, ok := f[0].Type.(*ast.SelectorExpr); ok {
				if se.Sel.Name == "Context" && se.X.(*ast.Ident).Name == "context" {
					hasContext = true
					start = 1
				}
			}

			for i := start; i < len(f); i++ {
				if len(f[i].Names) == 0 {
					fields = append(fields, fieldSchema("", "F"+strconv.FormatInt(int64(count), 10), f[i]))
					typs = append(typs, getTyp(f[i].Type))
				} else {
					for _, nativeName := range f[i].Names {
						fields = append(fields, fieldSchema(nativeName.Name, "F"+strconv.FormatInt(int64(count), 10), f[i]))
						typs = append(typs, getTyp(f[i].Type))
						count++
						//fmt.Println(typs)
					}
				}

			}
		}

	}

	// 如果函数最后一个是error， avro schema里面就把它去掉
	if fields != nil && len(fields) > 0 && len(fields[len(fields)-1]) > 38 &&
		fields[len(fields)-1][:38] == `{"name":"err","type":["string","null"]` {
		fields = fields[:len(fields)-1]
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
func fieldSchema(nativeName, fieldName string, f *ast.Field) (schema string) {

	fType := fieldKind(f.Type)
	if fType == `["string","null"]` {
		fieldName = "err"
	}

	schema = `{"name":"` + fieldName + `","type":` + fType

	if nativeName != "" {
		schema += `,"alias":"` + nativeName + `"`
	}

	schema += "}\n"

	return
}

func fieldKind(expr ast.Expr) (fType string) {

	switch in := expr.(type) {
	case *ast.StarExpr:
		panic("not support pointer")
	case *ast.ArrayType:
		// slice
		childKind := fieldKind(in.Elt)
		if childKind == "byte" || childKind == "uint8" {
			fType = `"bytes"`
		} else {
			fType = `{
				"type": "array",
				"items": ` + childKind + `
			}`
		}
	case *ast.SelectorExpr:
		// 其他包
		// TODO: 支持导入其他包的结构体
		panic("please do not use other package struct")
	case *ast.MapType:
		// map
		k := fieldKind(in.Key)
		if k != `"string"` {
			panic("map key must string but " + k)
		}
		v := fieldKind(in.Value)
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
			panic("avro not support uint")
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
			start := fset.Position(in.Obj.Decl.(*ast.TypeSpec).Pos())
			end := fset.Position(in.Obj.Decl.(*ast.TypeSpec).End())
			tmpFile, err := os.Open(start.Filename)
			if err != nil {
				panic(err)
			}
			b, err := ioutil.ReadAll(tmpFile)
			if err != nil {
				panic(err)
			}
			obj := "type " + string(b[start.Offset:end.Offset])
			f := in.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List
			key := start.Filename[:strings.LastIndex(start.Filename, string(os.PathSeparator))]
			if v, ok := objs[key]; ok {
				var flag bool
				for _, s := range v {
					if s == obj {
						flag = true
						break
					}
				}
				if !flag {
					objs[key] = append(v, obj)
				}
			} else {
				objs[key] = []string{obj}
			}

			fields := make([]string, len(f))

			for i := 0; i < len(f); i++ {
				var fieldName string
				if len(f[i].Names) == 0 {
					// 内嵌结构体
					fieldName = f[i].Type.(*ast.Ident).Name
				} else {
					fieldName = f[i].Names[0].Name
				}
				fields[i] = fieldSchema("", fieldName, f[i])
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
		//fmt.Println(expr.(*ast.SelectorExpr).Sel)
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
