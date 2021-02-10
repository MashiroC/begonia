package heartbeat

import (
	"errors"
	"github.com/MashiroC/begonia/dispatch/frame"
)

// 包括了dispatch的一些方法，用于解决循环导包
type Heartbeat interface {
	// Send 发送一个帧
	// 发送一个帧出去，在不同的集群模式下有不同的表现
	// link模式使用
	Send(frame frame.Frame) error

	// SendTo 发送帧到指定连接
	// set模式使用
	SendTo(connID string, f frame.Frame) error

}

var (
	PingTimeout = errors.New("ping timeout")
	PongTimeout = errors.New("pong timeout")
)
