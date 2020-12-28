// Package Server api层的service节点实现
package server

import (
	"github.com/MashiroC/begonia/internal/coding"
)

var success coding.SuccessCoder

// Server 服务端的接口的一份copy
type Server interface {
	Register(name string, service interface{}, registerFunc ...string)
	Wait()
}

// logic_service.go something
