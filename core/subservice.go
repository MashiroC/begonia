package core

import (
	"fmt"
	"log"
)

const (
	// ServiceName 默认的核心子服务服务名
	ServiceName = "CORE"
)

// C 核心子服务的单例
var C *SubService

// SubService 子服务
type SubService struct {
	services *registerServiceStore
}

// NewSubService 创建一个子服务
func NewSubService() *SubService {
	return &SubService{services: newStore()}
}

// Invoke 执行一个子服务的函数
func (s *SubService) Invoke(connID, reqID string, fun string, param []byte) (result []byte, err error) {
	switch fun {
	case "Register":
		var si ServiceInfo
		err = serviceInfoCoder.DecodeIn(param, &si)
		if err != nil {
			panic(err)
		}

		err = s.register(connID, si)
		if err != nil {
			return
		}
		result, err = successCoder.Encode(true)
	case "ServiceInfo":
		var call serviceInfoCall
		var si ServiceInfo
		err = serviceInfoCallCoder.DecodeIn(param, &call)
		if err != nil {
			return
		}

		si, err = s.serviceInfo(call.Service)
		if err != nil {
			return
		}
		result, err = serviceInfoCoder.Encode(si)
		if err != nil {
			return
		}
		fmt.Println(si)
	default:
		panic("err")
	}

	return
}

// GetToID 根据ServiceName获取连接id
func (s *SubService) GetToID(serviceName string) (connID string, ok bool) {
	service, ok := s.services.Get(serviceName)
	if !ok {
		return
	}
	connID = service.connID

	return
}

// HandleConnClose 默认的关闭连接钩子
func (s *SubService) HandleConnClose(connID string, err error) {
	log.Printf("conn [%s] closed, unlink service\n", connID)
	s.services.Unlink(connID)
}
