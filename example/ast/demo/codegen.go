package demo

import (
	"github.com/MashiroC/begonia/app/coding"
)

type CodeGenService interface {
	Do(fun string, param []byte) (result []byte, err error)
	FuncList() []coding.FunInfo
}

//var (
//	_DemoFuncList []coding.FunInfo
//
//	_DemoEchoDoInSchema  = `{}`
//	_DemoEchoDoOutSchema = `{}`
//	_DemoEchoDoInCoder   coding.Coder
//	_DemoEchoDoOutCoder  coding.Coder
//
//	_DemoAddDoInSchema  = `{}`
//	_DemoAddDoOutSchema = `{}`
//	_DemoAddDoInCoder   coding.Coder
//	_DemoAddDoOutCoder  coding.Coder
//)
//
//type _DemoEchoDoIn struct {
//	F1 string
//	F2 string
//}
//
//type _DemoEchoDoOut struct {
//	F1 string
//}
//
//type _DemoAddDoIn struct {
//	F1 int
//}
//
//type _DemoAddDoOut struct {
//	F1 int
//}
//
//func init() {
//	//var err error
//	//_DemoEchoDoInCoder, err = coding.NewAvro(_DemoEchoDoInSchema)
//	//if err != nil {
//	//	panic(err)
//	//}
//	//
//	//_DemoEchoDoOutCoder, err = coding.NewAvro(_DemoEchoDoOutSchema)
//	//if err != nil {
//	//	panic(err)
//	//}
//	//_DemoEchoDoInCoder, err = coding.NewAvro(_DemoEchoDoInSchema)
//	//if err != nil {
//	//	panic(err)
//	//}
//	//
//	//_DemoFuncList = []coding.FunInfo{
//	//	{Name: "Echo", Mode: "avro", InSchema: _DemoEchoDoInSchema, OutSchema: _DemoEchoDoOutSchema},
//	//	{Name: "Add", Mode: "avro", InSchema: _DemoAddDoInSchema, OutSchema: _DemoAddDoOutSchema},
//	//}
//}
//
//func (d *Demo) Do(fun string, param []byte) (result []byte, err error) {
//	switch fun {
//	case "Echo":
//		var in _DemoEchoDoIn
//		err := _DemoEchoDoInCoder.DecodeIn(param, &in)
//		if err != nil {
//			panic(err)
//		}
//		res1:= d.Echo(
//			in.F1,
//			in.F2,
//			)
//		var out _DemoEchoDoOut
//		out.F1 = res1
//		res, err := _DemoEchoDoOutCoder.Encode(out)
//		if err != nil {
//			panic(err)
//		}
//		return res, nil
//	case "Add":
//		var in _DemoAddDoIn
//		err := _DemoAddDoInCoder.DecodeIn(param, &in)
//		if err != nil {
//			panic(err)
//		}
//		res1, err := d.Add(in.F1)
//		var out _DemoAddDoOut
//		out.F1 = res1
//		res, err := _DemoAddDoOutCoder.Encode(out)
//		if err != nil {
//			panic(err)
//		}
//		return res, nil
//	default:
//		err = errors.New("rpc call error: fun not exist")
//	}
//	return
//}
//
//func (d *Demo) FuncList() []coding.FunInfo {
//	return _DemoFuncList
//}
