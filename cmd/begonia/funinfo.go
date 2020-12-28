package main

import (
	"errors"
	"github.com/MashiroC/begonia/internal/coding"
	"github.com/MashiroC/begonia/tool/qarr"
	"github.com/hamba/avro"
	"go/ast"
	"strings"
)

func getFunInfo(name string, decls []*ast.FuncDecl) (res []coding.FunInfo) {
	res = make([]coding.FunInfo, 0, 1)
	reFun, ok := nameRegister[name]
	for _, fd := range decls {
		if ok {
			if !qarr.StringsIn(reFun, fd.Name.Name) {
				continue
			}
		}

		funName := fd.Name.Name
		if funName[0] >= 'a' && funName[0] <= 'z' {
			continue
		}
		inSchema, inTyps, hasContext := MakeSchema(funName, "In", fd.Type.Params)
		outSchema, outTyps, _ := MakeSchema(funName, "Out", fd.Type.Results)
		res = append(res, coding.FunInfo{
			Name:       funName,
			InSchema:   inSchema,
			OutSchema:  outSchema,
			ParamTyp:   inTyps,
			ResultTyp:  outTyps,
			HasContext: hasContext,
		})
		avro.MustParse(inSchema)
		avro.MustParse(outSchema)
	}
	return
}

func parseTarget(pkgName string, node ast.Node) (res bool) {
	defer func() {
		if re := recover(); re != nil {
			res = false
		}
	}()

	targetName := node.(*ast.TypeSpec).Name.Name

	for _, t := range targetService {
		tmp := strings.Split(t, ":")
		if len(tmp) != 2 {
			panic("-t must a:b")
		}
		if tmp[0] == targetName {
			name := pkgName + "_" + targetName
			names[name] = `"` + tmp[1] + `"`
			return true
		}
	}

	return false
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
	if se.X.(*ast.Ident).Name == "begonia" && se.Sel.Name == "NewServer" {
		// 解析
		if len(call.Args) >= 2 {
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
			name := pkgName + "_" + ident.Name
			names[name] = call.Args[0].(*ast.BasicLit).Value

			if len(call.Args) > 2 {
				reFun := make([]string, 0, 1)
				for i := 2; i < len(call.Args); i++ {
					tmp := call.Args[i].(*ast.BasicLit).Value
					tmp = tmp[1 : len(tmp)-1]
					reFun = append(reFun, tmp)
				}
				nameRegister[name] = reFun
			}

			return true
		}
	}

	return
}

func parseStruct(pkgName string, node ast.Node) bool {
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

		recv = pkgName + "_" + recv

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

	}

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
