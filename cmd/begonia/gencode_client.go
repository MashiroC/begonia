package main

import (
	"encoding/json"
	"github.com/MashiroC/begonia/app/coding"
	"github.com/MashiroC/begonia/tool/qconv"
	"go/ast"
	"html/template"
	"os"
	"strings"
)

func genClientCode(f *ast.File, fullName string, fi []coding.FunInfo) {
	//var aliasMap map[string]string
	tmpl, err := template.New("client").Funcs(template.FuncMap{
		"raw":    raw,
		"concat": concat,
		"add":    add,
		"makeRes": func(s []string) (res string) {
			for i := 0; i < len(s); i++ {
				if s[i] == "error" && i == len(s)-1 {
					res += "err"
				} else {
					res += "res" + qconv.I2S(i+1)
				}
				res += ", "
			}
			if len(res) != 0 {
				res = res[:len(res)-2]
			}
			return
		},
		"hasRes": func(s []string) bool {
			return s != nil && len(s) != 0
		},
		"join": func(param []string, sep string) string {
			return strings.Join(param, sep)
		},
		"getNamesList": func(schema string) (res string) {
			var m map[string]interface{}
			json.Unmarshal([]byte(schema), &m)

			for _, v := range m["fields"].([]interface{}) {
				aliasIn, ok := v.(map[string]interface{})["alias"]
				if ok {
					res += aliasIn.(string) + ", "
				} else {
					res += v.(map[string]interface{})["name"].(string) + ", "
				}
			}

			if len(res) >= 2 {
				res = res[:len(res)-2]
			}
			return
		},
		"getFields": func(fun coding.FunInfo, typ string) (res string) {
			var m map[string]interface{}
			var list []string

			if typ == "in" {
				json.Unmarshal([]byte(fun.InSchema), &m)
				list = fun.ParamTyp
			} else if typ == "out" {
				json.Unmarshal([]byte(fun.OutSchema), &m)
				list = fun.ResultTyp
			} else {
				panic("error")
			}

			for i, v := range m["fields"].([]interface{}) {
				aliasIn, ok := v.(map[string]interface{})["alias"]
				if ok {
					res += aliasIn.(string) + " " + list[i] + ", "
				} else {
					res += v.(map[string]interface{})["name"].(string) + " " + list[i] + ", "
				}
			}

			if len(res) >= 2 {
				res = res[:len(res)-2]
			}
			return
		},
		"resultMode": func(resultTyp []string) (i int) {
			if len(resultTyp) > 1 {
				if len(resultTyp) == 2 && resultTyp[len(resultTyp)-1] == "error" {
					return 1
				}
				return 2
			}

			if len(resultTyp) == 1 && resultTyp[0] == "error" {
				return 0
			}

			return len(resultTyp)
		},
		"getAlias": func(schema string) (alias []string) {
			alias = make([]string, 0, 1)
			var m map[string]interface{}
			json.Unmarshal([]byte(schema), &m)

			for _, v := range m["fields"].([]interface{}) {
				aliasIn, ok := v.(map[string]interface{})["alias"]
				if ok {
					alias = append(alias, aliasIn.(string))
				} else {
					alias = append(alias, v.(map[string]interface{})["name"].(string))
				}
			}

			if len(alias) > 0 && alias[len(alias)-1] == "err" {
				alias = alias[:len(alias)-1]
			}

			return
		},
	}).Parse(clientTmplStr)
	if err != nil {
		panic(err)
	}

	tmp := strings.Split(fullName, ".")
	serviceName := tmp[len(tmp)-1]
	dirPath := strings.Join(tmp[:len(tmp)-3], string(os.PathSeparator))

	path := root + string(os.PathSeparator) + dirPath + string(os.PathSeparator) + serviceName + ".call.go"
	//fmt.Println(path)
	//fmt.Println(serviceName)

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	registerName := names[fullName]
	registerName = registerName[1 : len(registerName)-1]
	err = tmpl.Execute(file, map[string]interface{}{
		"source":       strings.Join(tmp[:len(tmp)-2], string(os.PathSeparator)) + ".go",
		"ServiceName":  serviceName,
		"RegisterName": registerName,
		"fi":           fi,
		"package":      f.Name,
		"alias":        "",
	})
	if err != nil {
		panic(err)
	}

}
