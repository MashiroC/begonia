package core

import (
	"github.com/MashiroC/begonia/app/coding"
	"github.com/MashiroC/begonia/logic/containers"
)

// call.go api层用来调用核心服务的工具函数
// 可以方便的直接获得一个logic.Call或logic.Result

type call int

// Call 单例，方便调用
var Call call

// Register 注册一个函数
func (call) Register(serviceName string, funs []coding.FunInfo) *containers.Call {

	b, err := serviceInfoCoder.Encode(ServiceInfo{
		Service: serviceName,
		Funs:    funs,
	})
	if err != nil {
		panic(err)
	}

	return &containers.Call{
		Service: ServiceName,
		Fun:     "Register",
		Param:   b,
	}
}

type serviceInfoCall struct {
	Service string `avro:"service"`
}

// ServiceInfo 获得服务信息
func (call) ServiceInfo(serviceName string) *containers.Call {
	b, err := serviceInfoCallCoder.Encode(serviceInfoCall{Service: serviceName})
	if err != nil {
		panic(err)
	}
	return &containers.Call{
		Service: ServiceName,
		Fun:     "ServiceInfo",
		Param:   b,
	}
}
