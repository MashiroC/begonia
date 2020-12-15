package main

import (
	"github.com/MashiroC/begonia/app/coding"
	"github.com/MashiroC/begonia/tool/qconv"
	"go/ast"
	"html/template"
	"os"
	"strings"
)

func genCode(f *ast.File, fullServiceName string, fi []coding.FunInfo) {
	tmpl, err := template.New("test").Funcs(template.FuncMap{
		"raw": func(str string) template.HTML {
			return template.HTML(str)
		},
		"concat": func(str ...string) string {
			res := ""
			for i := 0; i < len(str); i++ {
				res += str[i]
			}
			return res
		},
		"add": func(a int) int {
			return a + 1
		},
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
	}).Parse(tmplStr)
	if err != nil {
		panic(err)
	}

	tmp := strings.Split(fullServiceName, ".")
	serviceName := tmp[len(tmp)-1]
	dirPath := strings.Join(tmp[:len(tmp)-3], string(os.PathSeparator))

	path := root + string(os.PathSeparator) + dirPath + string(os.PathSeparator) + serviceName + ".begonia.go"

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = tmpl.Execute(file, map[string]interface{}{
		"source":      strings.Join(tmp[:len(tmp)-2], string(os.PathSeparator)) + ".go",
		"ServiceName": serviceName,
		"fi":          fi,
		"package":     f.Name,
	})
	if err != nil {
		panic(err)
	}
}
