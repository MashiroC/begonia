package core

import (
	"begonia2/logic"
	"begonia2/opcode/coding"
)

type call struct {
}

var Call call

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
