// Time : 2020/9/28 20:24
// Author : Kieran

// client
package app

import (
	"begonia2/logic"
	"begonia2/opcode/coding"
)

// coreservice.go something

var (
	CoreServiceName = "CoreBegonia"
)

var (
	Core                = &coreService{}
	signInfoCoder       coding.Coder
	signInfoResultCoder coding.Coder
	signCoder           coding.Coder
)

func init() {
	var err error

	signInfoCoder, err = coding.NewAvro(signInfoRawSchema)
	if err != nil {
		panic(err)
	}

	signInfoResultCoder, err = coding.NewAvro(signInfoResultRawSchema)
	if err != nil {
		panic(err)
	}

}

type coreService struct {
}

func (c *coreService) Sign(serviceName string, funs []coding.FunInfo) *logic.Call {
	//var fs []rFun
	//b, err := signCoder.Encode(fs)
	// TODO: 核心服务还没实现
	return &logic.Call{
		Service: "CORE",
		Fun:     "Sign",
		Param:   []byte{1,2,3},
	}
}

type SignInfoReq struct {
	Service string `avro:"service"`
}

func (c *coreService) SignInfo(serviceName string) *logic.Call {
	b, err := signInfoCoder.Encode(SignInfoReq{Service: serviceName})
	if err != nil {
		panic(err)
	}
	return &logic.Call{
		Service: CoreServiceName,
		Fun:     "SignInfo",
		Param:   b,
	}
}

func (c *coreService) SignInfoResult(b []byte) (f []FunInfo) {

	// TODO: 解码这个类型 构造coder

	//signInfoResultCoder.DecodeIn(b, &f)

	return
}
