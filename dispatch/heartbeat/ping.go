package heartbeat

import (
	"github.com/MashiroC/begonia/config"
	"github.com/MashiroC/begonia/dispatch/conn"
	"github.com/MashiroC/begonia/dispatch/frame"
	"time"
)

// 对ping方法的一些封装
type Ping struct {
	Code         byte // ping的负载
	SendPingTime time.Duration
	RecvPongTime time.Duration
	timer        *time.Timer
}

// 开始持续对一条连接发ping
func (p *Ping) Start(c conn.Conn) {
	pingFrame := frame.NewPing(p.Code)
	ticker := time.NewTicker(p.SendPingTime)

	// 判断是否超时
	go func() {
		<-p.timer.C
		c.Close()
	}()

	for {
		<-ticker.C
		_ = c.Write(byte(pingFrame.Opcode()), pingFrame.Marshal())
		p.timer.Reset(p.RecvPongTime)
	}
}

// 获取pong的内容，转化为映射
func (p *Ping) HandleFrame(f frame.Frame) map[string]string {
	if pongFrame, ok := f.(*frame.Pong); ok {
		p.timer.Stop()
		return pongFrame.Machine
	}

	return nil
}

func NewPing(code byte) *Ping {
	return &Ping{
		Code:         code,
		SendPingTime: config.C.Dispatch.SendPingTime,
		RecvPongTime: config.C.Dispatch.GetPongTime,
		timer:        time.NewTimer(0),
	}
}
