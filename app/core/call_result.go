package core

import (
	"begonia2/app"
	"begonia2/opcode/coding"
)

type result struct {
}

var Result result

func (result) ServiceInfo(b []byte) (f []app.FunInfo) {

	// TODO: 解码这个类型 构造coder
	var si ServiceInfo
	err := serviceInfoCoder.DecodeIn(b, &si)
	if err != nil {
		panic(err)
	}

	f=make([]app.FunInfo,len(si.Funs))

	for i:=0;i<len(f);i++{

		rawFun:=si.Funs[i]

		inCoder,err:=coding.NewAvro(rawFun.InSchema)
		if err!=nil{
			panic(err)
		}

		outCoder,err:=coding.NewAvro(rawFun.OutSchema)
		if err!=nil{
			panic(err)
		}

		f[i]=app.FunInfo{
			Name:     rawFun.Name,
			InCoder:  inCoder,
			OutCoder: outCoder,
		}

	}

	return
}
