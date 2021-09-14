package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
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
	targetService     []string
)

var (
	fset         = token.NewFileSet()
	names        = make(map[string]string)
	nameRegister = make(map[string][]string)
	recvs        = make(map[string]Service)
	objs         = make(map[string][]string)
	root         string
)

type Service struct {
	FuncList []*ast.FuncDecl
	File     *ast.File
}

func init() {
	targetServiceUsage := "generate target service (if not exist register code)"
	flag.BoolVarP(&isGenerateService, "server1", "s", false, "generate server1 code from begonia")
	flag.BoolVarP(&isGenerateClient, "client", "c", false, "generate client code from begonia")
	flag.BoolVarP(&isRemove, "remove", "r", false, "remove old begonia generate code")
	flag.StringSliceVarP(&targetService, "target", "t", []string{}, targetServiceUsage)
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
		fi := getFunInfo(k, v.FuncList)
		fmt.Println("generate server1", k, "...")
		if isGenerateService {
			fmt.Print("server1 code ...")
			genServiceCode(v.File, k, fi)
			fmt.Println("\b\b\bok!")
		}

		if isGenerateClient {
			fmt.Print("client call ...")
			genClientCode(k, fi)
			fmt.Println("\b\b\bok!")
		}

	}

	genEntity(objs)

	gofmt(originPath)
}

func gofmt(path string) {
	if !strings.HasSuffix(path, string(os.PathSeparator)) {
		path += string(os.PathSeparator)
	}

	c := exec.Command("go", "fmt", path+"...")
	err := c.Run()
	if err != nil {
		fmt.Println("go fmt error on path: [" + path + "...]")
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
		pkg = strings.Replace(path, root, "", 1)[1:]
	}

	ast.Inspect(f, func(node ast.Node) (res bool) {
		ok := parseTarget(pkg, node)
		if ok {
			return true
		}
		ok = parseObj(pkg, node)
		if ok {
			return true
		}
		ok = parseStruct(pkg, node)
		return true
	})

	return
}
