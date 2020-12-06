package main

import (
	"errors"
	"fmt"
	"github.com/MashiroC/begonia/app/coding"
	"github.com/MashiroC/begonia/example/ast/demo"
	"github.com/MashiroC/begonia/tool/qconv"
	"github.com/hamba/avro"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"os"
	"path/filepath"
)

func main() {
	fset := token.NewFileSet()
	fmt.Println(os.Getwd())
	path, _ := filepath.Abs("./example/ast/demo/dddd.go")
	f, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		panic(err)
	}

	serviceName, fi := parse(f)
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
	}).Parse(demo.TmplStr)
	if err != nil {
		panic(err)
	}

	file, err := os.Create("./example/ast/demo/demo_begonia.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = tmpl.Execute(file, map[string]interface{}{
		"source":      "example/ast/demo/dddd.go",
		"ServiceName": serviceName,
		"fi":          fi,
		"package":     f.Name,
	})
	if err != nil {
		panic(err)
	}
}

func parse(f *ast.File) (serviceName string, fi []coding.FunInfo) {
	fi = make([]coding.FunInfo, 0, 1)

	for i := 0; i < len(f.Decls); i++ {
		fd, ok := f.Decls[i].(*ast.FuncDecl)
		if !ok || fd.Recv == nil {
			continue
		}

		recv, err := getRecv(fd.Recv.List[0].Type)
		if err != nil {
			continue
		}
		serviceName = recv

		inSchema, inTyps := demo.MakeSchema(fd.Name.Name, "In", fd.Type.Params)
		outSchema, outTyps := demo.MakeSchema(fd.Name.Name, "Out", fd.Type.Results)
		fi = append(fi, coding.FunInfo{
			Name:      fd.Name.Name,
			Mode:      "avro",
			InSchema:  inSchema,
			OutSchema: outSchema,
			ParamTyp:  inTyps,
			ResultTyp: outTyps,
		})
		avro.MustParse(inSchema)
		avro.MustParse(outSchema)
	}
	//codegen:=demo.Codegen(inSchema,outSchema)

	return
}

func getRecv(expr ast.Expr) (string, error) {
	name, err := unPointer(expr)
	if err == nil {
		return name, nil
	}
	return "", err
}

func unPointer(expr ast.Expr) (name string, err error) {
	var ident *ast.Ident
	if star, ok := expr.(*ast.StarExpr); ok {
		// 取指针
		ident = star.X.(*ast.Ident)
		name = ident.Name
		return
	}
	err = errors.New("ast parse error: it not pointer")
	return
}
