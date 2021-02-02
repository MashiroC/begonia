// Package dispatch 通讯层，应用层发出请求通过通讯层的抽象。
package dispatch

import (
	"errors"
	"github.com/MashiroC/begonia/dispatch/frame"
	"reflect"
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
	Link(addr string) error

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

	// Listen 对一个地址开始监听
	Listen(addr string)

	// Close 释放资源
	Close()

	// Hook 对某些地方进行hook
	// 目前可以hook的：
	// - close
	Hook(typ string, hookFunc interface{})

	// Handle 对某些地方添加一个handle func来处理一些情况。
	// example:
	// dp.Handle("request",func(f *frame.Response) { fmt.Println(f) })
	// 目前可以handle的：
	// - client.handleResponse (response)
	// - client.handleRequest  (request)
	Handle(typ string, handleFunc interface{})

	// 用于获取一些信息，不同模式下获取的信息不同
	// 目前支持的：
	// - set：获取机器信息
	Get(id string) interface{}
}

type baseDispatch struct {

	// handle func
	LgHandleFrame func(connID string, f frame.Frame)

	// hook func
	CloseHookFunc func(connID string, err error) // 关闭连接的hook
}

func (d *baseDispatch) HandleFrame(connID string, f frame.Frame) {
	d.LgHandleFrame(connID, f)
}

func (d *baseDispatch) Handle(typ string, in interface{}) {
	switch typ {
	case "frame":
		if fun, ok := in.(func(connID string, f frame.Frame)); ok {
			d.LgHandleFrame = fun
			return
		}
	default:
		panic(errors.New("dispatch handle error: you handle func not exist"))
	}
	panic(errors.New("dispatch handle error: handle func not match"))
}

// Hook 在这里可以去Hook一些事件。
func (d *baseDispatch) Hook(name string, hookFunc interface{}) {
	switch name {
	case "close":
		if f, ok := hookFunc.(func(connID string, err error)); ok {
			d.CloseHookFunc = f
			return
		}
		panic("close func must func(connID string, err error) but " + reflect.TypeOf(hookFunc).String())
	default:
		panic("hook func " + name + " not exist")
	}
}
