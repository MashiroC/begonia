package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// begonia 代码生成 脚手架
// 以xxx文件来生成代码
// ./begonia -s ./demo.go
// 在xxx目录下查找所有注册在begonia上的服务，然后生成代码。
// ./begonia -s ./

var (
	isGenerateService bool
	isGenerateClient  bool
	isRemove          bool
)

var (
	fset  = token.NewFileSet()
	names = make(map[string]string)
	recvs = make(map[string]Service)
	root  string
)

type Service struct {
	FuncList []*ast.FuncDecl
	File     *ast.File
}

func init() {
	shortHand := " (shorthand)"
	isGenerateServiceUsage := "generate service code from begonia"
	isGenerateClientUsage := "generate client code from begonia"
	isRemoveUsage := "remove old begonia generate code"
	flag.BoolVar(&isGenerateService, "service", false, isGenerateServiceUsage)
	flag.BoolVar(&isGenerateService, "s", false, isGenerateServiceUsage+shortHand)
	flag.BoolVar(&isGenerateClient, "client", false, isGenerateClientUsage)
	flag.BoolVar(&isGenerateClient, "c", false, isGenerateClientUsage+shortHand)

	flag.BoolVar(&isRemove, "remove", false, isRemoveUsage)
	flag.BoolVar(&isRemove, "r", false, isRemoveUsage+shortHand)
	flag.Parse()
}

func main() {

	t := time.Now()
	defer func() {
		fmt.Println("complete, total:", time.Now().Sub(t))
	}()

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
	} else {
		panic("path not dir")
	}

	for k, _ := range names {
		v := recvs[k]
		fi := getFunInfo(v.FuncList)
		if isGenerateService {
			fmt.Print("generate service ", k, " code ...")
			genServiceCode(v.File, k, fi)
			fmt.Println("\b\b\bok!")
		}

		if isGenerateClient {
			//fmt.Print("generate service", k, "call ...")
			genClientCode(v.File, k, fi)
			//fmt.Println("\b\b\bok!")
		}

	}

	gofmt(originPath)
}

func gofmt(path string) {
	c := exec.Command("go", "fmt", path+"...")
	err := c.Run()
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

		work(fPath)
	}
}

func work(path string) {
	if !strings.HasSuffix(path, ".go") ||
		strings.HasSuffix(path, ".pb.go") ||
		strings.HasSuffix(path, "_test.go") {
		return
	}

	if strings.HasSuffix(path, ".begonia.go") || strings.HasSuffix(path, ".call.go") {
		if isRemove {
			remove(path)
		}
		return
	}

	f, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		return
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

	return
}
