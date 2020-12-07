package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/MashiroC/begonia/app/coding"
	"github.com/MashiroC/begonia/example/ast/demo"
	"github.com/MashiroC/begonia/tool/qconv"
	"github.com/hamba/avro"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// begonia 代码生成 脚手架
// 以xxx文件来生成代码
// ./begonia -g ./demo.go
// 在根目录下查找所有注册在begonia上的服务，然后生成代码。
// ./begonia -g .

var (
	isGenerate bool
	fset       = token.NewFileSet()
	names      = make(map[string]struct{})
	recvs      = make(map[string]Service)
	root       string
)

type Service struct {
	FuncList []*ast.FuncDecl
	File     *ast.File
}

func init() {
	shortHand := " (shorthand)"
	isGenerateUsage := "generate code from begonia"
	flag.BoolVar(&isGenerate, "generate", false, isGenerateUsage+shortHand)
	flag.BoolVar(&isGenerate, "g", false, isGenerateUsage)
	flag.Parse()
}

func main() {
	//filePath := os.Args
	//if *command {
	//
	//}
	t := time.Now()
	defer func() {
		fmt.Println(time.Now().Sub(t))
	}()

	if !isGenerate {
		return
	}

	originPath := os.Args[len(os.Args)-1]
	path, err := filepath.Abs(originPath)
	root = path
	if err != nil {
		panic(err)
	}

	fi, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	if fi.IsDir() {
		dfs(path)
	}

	for k, _ := range names {
		fmt.Print("generate service: ", k, " ...")
		v := recvs[k]
		fi := getFunInfo(v.FuncList)
		genCode(v.File, k, fi)
		fmt.Println("\b\b\bok!")
	}

	c := exec.Command("go", "fmt", originPath+"...")
	err = c.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func dfs(path string) {
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range fs {
		fPath := path + string(os.PathSeparator) + file.Name()
		fi, err := os.Stat(fPath)
		if err != nil {
			panic(err)
		}
		if fi.IsDir() {
			dfs(fPath)
			continue
		}

		if !strings.HasSuffix(fPath, ".go") ||
			strings.HasSuffix(fPath, ".pb.go") ||
			strings.HasSuffix(fPath, "_test.go") ||
			strings.HasSuffix(fPath, ".begonia.go") {
			continue
		}
		f, err := parser.ParseFile(fset, fPath, nil, parser.AllErrors)

		if err != nil {
			continue
		}

		var pkg string
		if path != root {
			pkg = strings.ReplaceAll(strings.Replace(path, root, "", 1)[1:], string(os.PathSeparator), ".")
		}
		ast.Inspect(f, func(node ast.Node) (res bool) {
			ok := parseObj(pkg, node)
			if ok {
				return true
			}
			ok = parseStruct(pkg, node)
			return true
		})
		//ast.Print(fset, f)
		//panic("asd")

	}
}

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
	}).Parse(demo.TmplStr)
	if err != nil {
		panic(err)
	}
	tmp := strings.Split(fullServiceName, ".")
	serviceName := tmp[len(tmp)-1]

	path := root + string(os.PathSeparator) + strings.ReplaceAll(fullServiceName, ".", string(os.PathSeparator)) + ".begonia.go"

	file, err := os.Create(path)
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

func getFunInfo(decls []*ast.FuncDecl) (res []coding.FunInfo) {
	res = make([]coding.FunInfo, 0, 1)
	for _, fd := range decls {
		inSchema, inTyps := demo.MakeSchema(fd.Name.Name, "In", fd.Type.Params)
		outSchema, outTyps := demo.MakeSchema(fd.Name.Name, "Out", fd.Type.Results)
		res = append(res, coding.FunInfo{
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
	return
}

func parseObj(pkgName string, node ast.Node) (res bool) {
	defer func() {
		if re := recover(); re != nil {
			res = false
		}
	}()

	// 别问我为什么这么写 语法树就是这样的
	call := node.(*ast.CallExpr)
	se := call.Fun.(*ast.SelectorExpr).X.(*ast.Ident).Obj.Decl.(*ast.AssignStmt).Rhs[0].(*ast.CallExpr).Fun.(*ast.SelectorExpr)
	if se.X.(*ast.Ident).Name == "begonia" && se.Sel.Name == "NewService" {
		// 解析
		if len(call.Args) == 2 {
			var ue *ast.UnaryExpr
			if tmp, ok := call.Args[1].(*ast.Ident); ok {
				ue = tmp.Obj.Decl.(*ast.AssignStmt).Rhs[0].(*ast.UnaryExpr)
			} else {
				ue = call.Args[1].(*ast.UnaryExpr)
			}
			var ident *ast.Ident
			if cl, ok := ue.X.(*ast.CompositeLit); ok {
				ident = cl.Type.(*ast.Ident)
			} else {
				ident = ue.X.(*ast.Ident).Obj.Decl.(*ast.AssignStmt).Rhs[0].(*ast.CallExpr).Fun.(*ast.Ident)
			}
			name := pkgName + "." + ident.Name
			names[name] = struct{}{}

			return true
		}
	}

	return
}

func parseStruct(pkgName string, node ast.Node) bool {
	//serviceName string, fi []coding.FunInfo

	//res = make(map[string][]coding.FunInfo)
	f, ok := node.(*ast.File)
	if !ok {
		return false
	}

	for i := 0; i < len(f.Decls); i++ {
		fd, ok := f.Decls[i].(*ast.FuncDecl)
		if !ok || fd.Recv == nil {
			continue
		}

		recv, err := getRecv(fd.Recv.List[0].Type)
		if err != nil {
			continue
		}

		recv = pkgName + "." + recv

		if re, ok := recvs[recv]; ok {
			recvs[recv] = Service{
				FuncList: append(re.FuncList, fd),
				File:     f,
			}
		} else {
			recvs[recv] = Service{
				FuncList: []*ast.FuncDecl{fd},
				File:     f,
			}
		}

		//	if _, ok := res[recv]; !ok {
		//		res[recv] = make([]coding.FunInfo, 0, 1)
		//	}
		//
		//	inSchema, inTyps := demo.MakeSchema(fd.Name.Name, "In", fd.Type.Params)
		//	outSchema, outTyps := demo.MakeSchema(fd.Name.Name, "Out", fd.Type.Results)
		//	res[recv] = append(res[recv], coding.FunInfo{
		//		Name:      fd.Name.Name,
		//		Mode:      "avro",
		//		InSchema:  inSchema,
		//		OutSchema: outSchema,
		//		ParamTyp:  inTyps,
		//		ResultTyp: outTyps,
		//	})
		//	avro.MustParse(inSchema)
		//	avro.MustParse(outSchema)
	}
	//codegen:=demo.Codegen(inSchema,outSchema)

	return false
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
