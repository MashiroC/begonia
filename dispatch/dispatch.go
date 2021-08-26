// Package dispatch 通讯层，应用层发出请求通过通讯层的抽象。
package dispatch

import (
	"errors"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/dispatch/router"
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

	// Start 启动
	// 会根据不同的dispatch调用不同的初始化，例如link会调Link，set会调Listen
	Start(addr string) error

	// Send 发送一个帧
	// 发送一个帧出去，在不同的集群模式下有不同的表现
	// - default:
	// 发送到服务中心
	// - other:
	// 未实现
	Send(frame frame.Frame) error

	// SendTo 发送帧到指定连接
	SendTo(connID string, f frame.Frame) error

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
	// - frame
	// - proxy
	// - ctrl
	Handle(typ string, handleFunc interface{})

	// Upgrade 将连接进行升级
	Upgrade(connID string, addr string) error
}

type baseDispatch struct {

	// hook func
	CloseHookFuncList []func(connID string, err error) // 关闭连接的hook

	LinkHookFuncList []func(connID string) // 启动连接的hook

	rt *router.Router
}

func (d *baseDispatch) Handle(typ string, in interface{}) {

	if d.rt == nil {
		d.rt = router.New()
	}
	switch typ {
	case "frame":
		if fun, ok := in.(func(connID string, f frame.Frame)); ok {
			d.rt.LgHandleFrame = fun
			return
		}
	case "ctrl":
		if fun, ok := in.(func() (code int, fun router.CtrlHandleFunc)); ok {
			code, f := fun()
			d.rt.AddCtrlHandle(code, f)
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
			d.CloseHookFuncList = append(d.CloseHookFuncList, f)
			return
		}
		panic("close func must func(connID string, err error) but " + reflect.TypeOf(hookFunc).String())
	case "link":
		if f, ok := hookFunc.(func(connID string)); ok {
			d.LinkHookFuncList = append(d.LinkHookFuncList, f)
			return
		}
		panic("start func must func(connID string) but " + reflect.TypeOf(hookFunc).String())
	default:
		panic("hook func " + name + " not exist")
	}
}

func (d *baseDispatch) DoCloseHook(connID string, err error) {
	if d.CloseHookFuncList != nil {
		for _, f := range d.CloseHookFuncList {
			f(connID, err)
		}
	}
}

func (d *baseDispatch) DoLinkHook(connID string) {
	if d.LinkHookFuncList != nil {
		for _, f := range d.LinkHookFuncList {
			f(connID)
		}
	}
}