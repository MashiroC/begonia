package heartbeat

import (
	"context"
	"github.com/MashiroC/begonia/config"
	"github.com/MashiroC/begonia/dispatch/frame"
	"time"
)

// 对ping方法的一些封装
type Ping struct {
	Code         byte          // ping的负载，用于标识pong帧需要返回的内容
	SendPingTime time.Duration // 发送ping帧的时间间隔
	RecvPongTime time.Duration // 发送ping帧后收到pong的最长时间
	timer        *time.Timer
	ConnId       string
}

var PingUtil *Ping

// 开始持续对一条连接发ping
func (p *Ping) Start() {
	ticker := time.NewTicker(p.SendPingTime)
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		for {
			select {
			case <-ticker.C:
				pingFrame := frame.NewPing(p.Code)
				_ = dispatch.SendTo(p.ConnId, pingFrame)
				p.timer.Reset(p.RecvPongTime)
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	// 判断是否超时
	<-p.timer.C
	cancel()
	dispatch.Close()
}

// 获取pong的内容（机器信息），转化为映射
// 在这里暂停计时器，代表已经收到pong，发送ping的时候会reset
func HandlePong(f frame.Frame) {
	if pongFrame, ok := f.(*frame.Pong); ok {
		PingUtil.timer.Stop()
		dispatch.Store(PingUtil.ConnId, pongFrame.Machine)
	}
	return
}

func NewPing(code byte, connId string, hb Heartbeat) *Ping {
	dispatch = hb
	PingUtil = &Ping{
		Code:         code,
		SendPingTime: config.C.Dispatch.SendPingTime,
		RecvPongTime: config.C.Dispatch.GetPongTime,
		timer:        time.NewTimer(time.Hour),
		ConnId:       connId,
	}
	return PingUtil
}
