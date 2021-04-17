package heartbeat

import (
	"context"
	"github.com/MashiroC/begonia/config"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/tool/machine"
	"log"
	"time"
)

// Pong 对pong方法的一些封装
type Pong struct {
	RecvPingTime time.Duration // 收到ping帧的最长时间间隔
	timer        *time.Timer

	ch      chan struct{} // 收到ping帧后传递信息，重置计时器
	machine *machine.Machine

	connID string
	Close  func()                                   // 关闭连接，以及hook的方法
	Send   func(connID string, f frame.Frame) error // 发送帧
}

// Start 一定时间内没收到pong就断开连接
func (p *Pong) Start(c context.Context) {

	go func(ctx context.Context) {
		for {
			select {
			case <-p.ch:
				p.timer.Stop()
				p.timer.Reset(p.RecvPingTime)
			case <-ctx.Done():
				return
			}
		}
	}(c)

	// 超时主动关闭,调用Close方法
	// 如果是context通知关闭，说明是Hook了close,直接退出goroutine即可
	p.timer.Reset(p.RecvPingTime)
	select {
	case <-p.timer.C:
		p.Close()
		break
	case <-c.Done():
		break
	}
	close(p.ch)
}

// Handle 根据ping帧，返回pong帧
func (p *Pong) Handle(f frame.Frame) {
	if pingFrame, ok := f.(*frame.Ping); ok {
		p.ch <- struct{}{}
		m := p.machine.GetMachineInfo(pingFrame.Code)
		pongFrame := frame.NewPong(m, nil)
		err := p.Send(p.connID, pongFrame)
		if err != nil {
			log.Println(err)
		}
	}
	return
}

func (p *Pong) RecvType() int {
	return frame.PingTypCode
}

func NewPong(connID string, close func(), send func(connID string, f frame.Frame) error) *Pong {
	return &Pong{
		RecvPingTime: config.C.Dispatch.GetPingTime,
		timer:        time.NewTimer(config.C.Dispatch.GetPingTime),
		ch:           make(chan struct{}),
		machine:      machine.NewMachine(),
		Close:        close,
		Send:         send,
		connID:       connID,
	}
}
