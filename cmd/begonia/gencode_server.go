package main

import (
	"github.com/MashiroC/begonia/app/coding"
	"github.com/MashiroC/begonia/tool/qconv"
	"go/ast"
	"html/template"
	"os"
	"strings"
)

func genServiceCode(f *ast.File, fullServiceName string, fi []coding.FunInfo) {
	tmpl := getServiceTmpl()

	serviceName, dirPath := parseFullName(fullServiceName)

	path := strings.Join([]string{root, dirPath, serviceName + ".begonia.go"}, string(os.PathSeparator))

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = tmpl.Execute(file, map[string]interface{}{
		"source":      fullServiceName[:strings.LastIndex(fullServiceName, "_")],
		"ServiceName": serviceName,
		"fi":          fi,
		"package":     f.Name,
	})
	if err != nil {
		panic(err)
	}
}

func getServiceTmpl() *template.Template {
	tmpl, err := template.New("server").Funcs(template.FuncMap{
		"raw":    raw,
		"concat": concat,
		"add":    add,
		"makeRes": func(s []string) (res string) {
			for i := 0; i < len(s); i++ {
				if s[i] == "error" && i == len(s)-1 {
					res += "inErr"
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
		"isLastError": func(s []string) bool {
			return s[len(s)-1] == "error"
		},
		"hasRes": func(s []string) bool {
			return s != nil && len(s) != 0
		},
	}).Parse(serviceTmplStr)
	if err != nil {
		panic(err)
	}
	return tmpl
}
