package heartbeat

import "github.com/MashiroC/begonia/dispatch/frame"

// 包括了dispatch的一些方法，用于解决循环导包
type Heartbeat interface {
	// Send 发送一个帧
	// 发送一个帧出去，在不同的集群模式下有不同的表现
	// link模式使用
	Send(frame frame.Frame) error

	// SendTo 发送帧到指定连接
	// set模式使用
	SendTo(connID string, f frame.Frame) error

	// Close 释放资源
	Close()

	//存储机器信息
	Store(id string, machine map[string]string)
}

var dispatch Heartbeat
