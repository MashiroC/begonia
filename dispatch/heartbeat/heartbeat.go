package heartbeat

import (
	"context"
	"errors"
	"github.com/MashiroC/begonia/dispatch/frame"
	"log"
	"sync"
)

// 对于单个dispatch的ping或pong
type Beat interface {
	// 开起计时器，发送ping帧或收到pong帧超时
	Start(ctx context.Context)

	// 处理收到的心跳帧，和dispatch的handle不一样
	Handle(f frame.Frame)

	// 标识是ping还是pong
	RecvType() int
}

var (
	PingTimeout = errors.New("ping timeout")
	PongTimeout = errors.New("pong timeout")
)

// 心跳用到的dispatch的两个主要函数
type closeFunc func()
type sendFunc func(connID string, f frame.Frame) error

// 连接多个对象时，用于注册多个beat对象
// 以及方便sarter调用
// 默认一条连接，其中一方只发ping，另一方只发pong
type Heart struct {
	beats map[string]Beat
	sync.Mutex
}

// 注册一个新的连接(dispatch)，并返回一个用于关闭goroutine的函数
// 在连接建立时hook
func (h *Heart) Register(typ string, connID string, close closeFunc, send sendFunc) func() {
	var beat Beat
	ctx, cancelFunc := context.WithCancel(context.Background())
	switch typ {
	case "ping":
		ping := NewPing(frame.PingPongCtrlCode, connID, close, send)
		go ping.Start(ctx)
		beat = ping

	case "pong":
		pong := NewPong(connID, close, send)
		go pong.Start(ctx)
		beat = pong

	default:
		panic("heartbeat type error: unknown heartbeat type " + "typ")
	}

	h.Lock()
	h.beats[connID] = beat
	h.Unlock()
	return cancelFunc
}

// 处理某个连接的心跳帧
// 收到帧时handle
func (h *Heart) Handle(connID string, data []byte) {
	h.Lock()
	beat, ok := h.beats[connID]
	h.Unlock()
	if beat == nil || !ok {
		log.Println("nil beat")
		return
	}

	// 反序列化出一个心跳帧
	f, err := frame.UnMarshalPingPong(beat.RecvType(), data)
	if err != nil {
		log.Println(err)
		return
	}

	beat.Handle(f)
}

func NewHeart() *Heart {
	return &Heart{
		beats: make(map[string]Beat),
		Mutex: sync.Mutex{},
	}
}

type handleFunc = func(connID string, data []byte)

// 在dispatch启动时住处handle
func Handler(h *Heart) (f func() (code int, fun handleFunc)) {
	return func() (code int, fun handleFunc) {
		return frame.PingPongCtrlCode, func(connID string, data []byte) {
			h.Handle(connID, data)
		}
	}
}

