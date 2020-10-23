package core

import (
	"begonia2/app/coding"
	"begonia2/logic"
)

// call.go api层用来调用核心服务的工具函数
// 可以方便的直接获得一个logic.Call或logic.Result

type call int

// Call 单例，方便调用
var Call call

// Register 注册一个函数
func (call) Register(serviceName string, funs []coding.FunInfo) *logic.Call {

	b, err := serviceInfoCoder.Encode(ServiceInfo{
		Service: serviceName,
		Funs:    funs,
	})
	if err != nil {
		panic(err)
	}

	return &logic.Call{
		Service: ServiceName,
		Fun:     "Register",
		Param:   b,
	}
}

type serviceInfoCall struct {
	Service string `avro:"service"`
}

// ServiceInfo 获得服务信息
func (call) ServiceInfo(serviceName string) *logic.Call {
	b, err := serviceInfoCallCoder.Encode(serviceInfoCall{Service: serviceName})
	if err != nil {
		panic(err)
	}
	return &logic.Call{
		Service: ServiceName,
		Fun:     "ServiceInfo",
		Param:   b,
	}
}
