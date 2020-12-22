// Package dispatch 通讯层，应用层发出请求通过通讯层的抽象。
package dispatch

import (
	"github.com/MashiroC/begonia/dispatch/frame"
	"sync"
)

/*
 通讯层有三种类型。
 default cluster (实现中)
 p2p cluster(计划中)
 manager cluster (计划中)
*/

// Dispatcher 通讯层的对外暴露的接口
type Dispatcher interface {

	// Link 连接到某个服务或中心
	// 会直接连接到指定的地址，[error]是用来返回连接时候的错误值的。
	// 连接断开不会在这里返回错误，而是提供一个hook，通过hook "close" 来捕获断开连接
	Link(config map[string]interface{}) error

	// ReLink 重新连接
	// 需要先调用 Link 之后才能调用ReLink，相当于是重新调用了一次Link，返回这次重连是否成功
	ReLink() bool

	// Send 发送一个帧
	// 发送一个帧出去，在不同的集群模式下有不同的表现
	// - default:
	// 发送到服务中心
	// - other:
	// 未实现
	Send(frame frame.Frame) error

	// SendTo 发送帧到指定连接
	SendTo(connID string, f frame.Frame) error

	// Recv 接收一个请求
	Recv() (connID string, f frame.Frame)

	// Listen 对一个地址开始监听
	Listen(addr string)

	// Close 释放资源
	Close()

	// Hook 对某些地方进行hook
	// 目前可以hook的：
	// - close
	Hook(typ string, hookFunc interface{})
}

type Machine struct {
	sync.Mutex
	Info map[string]string
}

func (m *Machine) Get(key string) (value string, has bool) {
	m.Lock()
	defer m.Unlock()
	value, has = m.Info[key]
	return
}

func (m *Machine) StoreMachine(info map[string]string) {
	m.Lock()
	defer m.Unlock()
	m.Info = info
}
