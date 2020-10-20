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
	Core             = &coreService{}
	signInfoCoder    coding.Coder
	serviceInfoCoder coding.Coder

	successCoder coding.Coder
)

func init() {
	var err error

	signInfoCoder, err = coding.NewAvro(signInfoRawSchema)
	if err != nil {
		panic(err)
	}

	serviceInfoCoder, err = coding.NewAvro(serviceInfoRawSchema)
	if err != nil {
		panic(err)
	}

	successCoder = &coding.SuccessCoder{}

}

type coreService struct {
}

func (c *coreService) SignCall(serviceName string, funs []coding.FunInfo) *logic.Call {

	b, err := serviceInfoCoder.Encode(SignInfo{
		Service: serviceName,
		Funs:    funs,
	})
	if err != nil {
		panic(err)
	}

	return &logic.Call{
		Service: "CORE",
		Fun:     "SignCall",
		Param:   b,
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

	//serviceInfoCoder.DecodeIn(b, &f)

	return
}
