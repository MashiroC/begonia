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

// 开始持续对一条连接发ping
func (p *Ping) Start(hb Heartbeat) {
	ticker := time.NewTicker(p.SendPingTime)
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		for {
			select {
			case <-ticker.C:
				pingFrame := frame.NewPing(p.Code)
				_ = hb.SendTo(p.ConnId, pingFrame)
				p.timer.Reset(p.RecvPongTime)
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	// 判断是否超时
	<-p.timer.C
	cancel()
	hb.Close()
}

// 获取pong的内容（机器信息），转化为映射
// 在这里暂停计时器，因为已经没必要了（不是重置）
func (p *Ping) HandleFrame(f frame.Frame) map[string]string {
	if pongFrame, ok := f.(*frame.Pong); ok {
		p.timer.Stop()
		return pongFrame.Machine
	}

	return nil
}

func NewPing(code byte, connId string) *Ping {
	return &Ping{
		Code:         code,
		SendPingTime: config.C.Dispatch.SendPingTime,
		RecvPongTime: config.C.Dispatch.GetPongTime,
		timer:        time.NewTimer(time.Hour),
		ConnId:       connId,
	}
}
