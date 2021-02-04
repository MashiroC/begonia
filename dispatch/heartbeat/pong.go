package heartbeat

import (
	"context"
	"github.com/MashiroC/begonia/config"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/tool/machine"
	"time"
)

// 对pong方法的一些封装
type Pong struct {
	RecvPingTime time.Duration // 收到ping帧的最长时间间隔
	timer        *time.Timer

	ch      chan struct{} // 收到ping帧后传递信息，重置计时器
	machine *machine.Machine
}

var PongUtil *Pong

// 一定时间内没收到pong就断开连接
func startPong() {
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		for {
			select {
			case <-PongUtil.ch:
				PongUtil.timer.Stop()
				PongUtil.timer.Reset(PongUtil.RecvPingTime)
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	PongUtil.timer.Reset(PongUtil.RecvPingTime)
	<-PongUtil.timer.C
	cancel()
	dispatch.Close()
}

// 根据ping帧，返回pong帧
func HandlePing(f frame.Frame) {
	if pingFrame, ok := f.(*frame.Ping); ok {
		PongUtil.ch <- struct{}{}
		m := PongUtil.machine.GetMachineInfo(pingFrame.Code)
		pongFrame := frame.NewPong(m, nil)
		dispatch.Send(pongFrame)
	}
	return
}

func NewPong(hb Heartbeat) {
	dispatch = hb
	PongUtil = &Pong{
		RecvPingTime: config.C.Dispatch.GetPingTime,
		timer:        time.NewTimer(0),
		ch:           make(chan struct{}),
		machine:      machine.NewMachine(),
	}
	go startPong()
}
