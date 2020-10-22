// Time : 2020/9/19 15:10
// Author : Kieran

// dispatch 通讯层，应用层发出请求通过通讯层的抽象。
package dispatch

import (
	"begonia2/dispatch/frame"
)

/*
 通讯层有两个实现。

 default  -  默认，begonia的tcp连接方式
 grpc     -  使用grpc的实现
*/
type Dispatcher interface {

	// Link 连接到某个服务或中心
	Link(addr string)

	// Send 发送请求
	Send(frame frame.Frame) error

	SendTo(connID string, f frame.Frame) error

	// Recv 接收到一个请求 对于客户端是接收到了响应 对于服务端是接收到了一个请求
	Recv() (connID string,f frame.Frame)

	Listen(addr string)

	Close()

	Hook(typ string,hookFunc interface{})
}
