package main

import (
	"encoding/json"
	"github.com/MashiroC/begonia/internal/coding"
	"html/template"
	"os"
	"strings"
)

func genClientCode(fullName string, fi []coding.FunInfo) {
	//var aliasMap map[string]string
	tmpl := getClientTmpl()

	serviceName, dirPath := parseFullName(fullName)

	path := strings.Join([]string{root, dirPath, "call"}, string(os.PathSeparator))

	makeCall(path)

	file, err := os.Create(path + string(os.PathSeparator) + serviceName + ".begonia.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	registerName := names[fullName]
	registerName = registerName[1 : len(registerName)-1]
	err = tmpl.Execute(file, map[string]interface{}{
		"source":       fullName[:strings.LastIndex(fullName, "_")],
		"ServiceName":  serviceName,
		"RegisterName": registerName,
		"fi":           fi,
		"alias":        "",
	})
	if err != nil {
		panic(err)
	}

}

func makeCall(path string) {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(err)
	}

	_, err = os.Stat(path + string(os.PathSeparator) + "call.go")
	if err != nil {
		callFile, err := os.Create(path + string(os.PathSeparator) + "call.begonia.go")
		if err != nil {
			panic(err)
		}
		callTmpl, err := template.New("call").Parse(baseClientTmplStr)
		err = callTmpl.Execute(callFile, nil)
		if err != nil {
			panic(err)
		}
	}
}

func getClientTmpl() *template.Template {
	tmpl, err := template.New("client").Funcs(template.FuncMap{
		"raw":    raw,
		"concat": concat,
		"add":    add,
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

			if typ == "out" {
				res += "err error, "
			}

			if len(res) >= 2 {
				res = res[:len(res)-2]
			}
			return
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
	return tmpl
}
