package core

import (
	"github.com/MashiroC/begonia/app"
	"github.com/MashiroC/begonia/app/coding"
)

type result struct {
}

// Result 单例，用来获得远程函数调用的结果
var Result result

// ServiceInfo 服务信息
func (result) ServiceInfo(b []byte) (f []app.FunInfo) {

	// TODO: 解码这个类型 构造coder
	var si ServiceInfo
	err := serviceInfoCoder.DecodeIn(b, &si)
	if err != nil {
		panic(err)
	}

	f = make([]app.FunInfo, len(si.Funs))

	for i := 0; i < len(f); i++ {

		rawFun := si.Funs[i]

		inCoder, err := coding.NewAvro(rawFun.InSchema)
		if err != nil {
			panic(err)
		}

		outCoder, err := coding.NewAvro(rawFun.OutSchema)
		if err != nil {
			panic(err)
		}

		f[i] = app.FunInfo{
			Name:     rawFun.Name,
			InCoder:  inCoder,
			OutCoder: outCoder,
		}

	}

	return
}
