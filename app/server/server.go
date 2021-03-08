// Package Server api层的service节点实现
package server

// Server 服务端的接口的一份copy
type Server interface {

	// Register 注册服务
	Register(name string, service interface{}, registerFunc ...string)

	// SetLoggerService 注册日志服务
	SetLoggerService()

	// Wait 阻塞等待
	Wait()
}

// logic_service.go something
