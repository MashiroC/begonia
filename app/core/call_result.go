package core

import "begonia2/app"

type result struct {
}

var Result result

func (result) ServiceInfo(b []byte) (f []app.FunInfo) {

	// TODO: 解码这个类型 构造coder

	//serviceInfoCoder.DecodeIn(b, &f)

	return
}
