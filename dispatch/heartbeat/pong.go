package heartbeat

import (
	"github.com/MashiroC/begonia/config"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/tool/machine"
	"time"
)

// 对pong方法的一些封装
type Pong struct {
	RecvPingTime time.Duration
	timer        *time.Timer
	// 收到ping帧后取消计时器
	ch           chan struct{}
	machine      *machine.Machine
}

// 一定时间内没收到pong就断开连接
func (p *Pong) Start(c Heartbeat) {
	go func() {
		<-p.timer.C
		c.Close()
	}()

	for {
		<-p.ch
		p.timer.Stop()
		p.timer.Reset(p.RecvPingTime)
	}
}

// 根据ping帧，返回pong帧
func (p *Pong) HandleFrame(f frame.Frame) frame.Frame {
	if pingFrame, ok := f.(*frame.Ping); ok {
		p.ch <- struct{}{}
		m := p.machine.GetMachineInfo(pingFrame.Code)
		pongFrame := frame.NewPong(m, nil)
		return pongFrame
	}
	return nil
}

func NewPong() *Pong {
	return &Pong{
		RecvPingTime: config.C.Dispatch.GetPingTime,
		timer:        time.NewTimer(0),
		ch:           make(chan struct{}),
		machine:      machine.NewMachine(),
	}
}
