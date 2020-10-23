// Package service api层的service节点实现
package service

import (
	"begonia2/app/coding"
)

var success coding.SuccessCoder

// logic_service.go something

// Service 服务端的接口
type Service interface {
	Register(name string, service interface{})
	Wait()
}

// astService ast树代码生成的ast service api
type astService struct {
}