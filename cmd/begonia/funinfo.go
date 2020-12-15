package main

import (
	"errors"
	"github.com/MashiroC/begonia/app/coding"
	"github.com/hamba/avro"
	"go/ast"
)

func getFunInfo(decls []*ast.FuncDecl) (res []coding.FunInfo) {
	res = make([]coding.FunInfo, 0, 1)
	for _, fd := range decls {
		inSchema, inTyps := MakeSchema(fd.Name.Name, "In", fd.Type.Params)
		outSchema, outTyps := MakeSchema(fd.Name.Name, "Out", fd.Type.Results)
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
