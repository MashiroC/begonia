package core

import "github.com/MashiroC/begonia/internal/coding"

// ServiceInfo 服务信息
// 在注册服务时，为入参
// 获取服务时 为出参
type ServiceInfo struct {
	Service string           `avro:"server"`
	Funs    []coding.FunInfo `avro:"funs"`
}
