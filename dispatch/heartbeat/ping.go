package heartbeat

import (
	"context"
	"github.com/MashiroC/begonia/config"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/tool/storage"
	"time"
)

// 对ping方法的一些封装
type Ping struct {
	Code         byte          // ping的负载，用于标识pong帧需要返回的内容
	SendPingTime time.Duration // 发送ping帧的时间间隔
	RecvPongTime time.Duration // 发送ping帧后收到pong的最长时间
	timer        *time.Timer
	ConnId       string

	Close func()                                   // 关闭连接，以及hook的方法
	Send  func(connID string, f frame.Frame) error // 发送帧
}

// Start 开始持续对一条连接发ping
func (p *Ping) Start(c context.Context) {
	ticker := time.NewTicker(p.SendPingTime)

	go func(ctx context.Context) {
		for {
			select {
			case <-ticker.C:
				pingFrame := frame.NewPing(p.Code)
				_ = p.Send(p.ConnId, pingFrame)
				p.timer.Stop()
				p.timer.Reset(p.RecvPongTime)

			case <-ctx.Done():
				return
			}
		}
	}(c)

	// 判断是否超时,或者断连
	// 超时主动关闭,调用Close方法
	// 如果是context通知关闭，说明是Hook了close,直接退出goroutine即可
	select {
	case <-p.timer.C:
		p.Close()
		break
	case <-c.Done():
		break
	}
}

// Handle 获取pong的内容（机器信息），转化为映射
// 在这里暂停计时器，代表已经收到pong，发送ping的时候会reset
func (p *Ping) Handle(f frame.Frame) {
	if pongFrame, ok := f.(*frame.Pong); ok {
		p.timer.Stop()
		storage.Store(p.ConnId, pongFrame.Machine)
	}
	return
}

func (p *Ping) RecvType() int {
	return frame.PongTypCode
}

func NewPing(code byte, connID string, close func(), send func(string, frame.Frame) error) *Ping {
	return &Ping{
		Code:         code,
		SendPingTime: config.C.Dispatch.SendPingTime,
		RecvPongTime: config.C.Dispatch.GetPongTime,
		timer:        time.NewTimer(time.Hour),
		Close:        close,
		Send:         send,
		ConnId:       connID,
	}
}
