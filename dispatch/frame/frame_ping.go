package frame

import (
	"bytes"
	"github.com/MashiroC/begonia/tool/qconv"
	"strconv"
	"time"
)

const (
	// ping的typCode
	pingTypCode = 0
)

// Ping Request的frame实现
//
// opcode4 length8 extendLength16
// req:service fun reqId param
//     4      4         8       0 || 16   [              length                  ]
// {opcode}{version}{length}{extendLength}{reqId}0x49{service}0x49{fun}0x49{param}
//
type Ping struct {
	PingPongTime time.Duration

	v      []byte // 序列化后的payload，这里是一个缓存
	opcode int    // 序列化后的opcode，初始化为-1
}

// NewPing 创建一个新的Request
func NewPing(pingPongTime time.Duration) Frame {
	return &Ping{
		PingPongTime: pingPongTime,
		opcode:       -1,
	}
}

// unMarshalRequest 根据payload去反序列化出一个request
func unMarshalPing(data []byte) (ping *Ping, err error) {

	ping = &Ping{}
	buf := bytes.NewBuffer(data)
	t := buf.Bytes()
	ppt, _ := strconv.ParseInt(qconv.Qb2s(t), 10, 64)
	ping.PingPongTime = time.Duration(ppt)
	ping.v = data
	ping.opcode = -1

	return
}

// Marshal 序列化payload
//
//      4      4         8       0 || 16   [              length                  ]
//	{opcode}{version}{length}{extendLength}{reqId}0x49{service}0x49{fun}0x49{param}
//
func (ping *Ping) Marshal() []byte {

	if ping.v == nil {
		buf := make([]byte, 0, 128)

		t := strconv.Itoa(int(ping.PingPongTime))
		buf = append(buf, qconv.Qs2b(t)...)

		ping.v = buf
	}

	return ping.v
}

// Opcode 组装出一个opcode
func (ping *Ping) Opcode() int {
	if ping.opcode == -1 {
		ping.opcode = makeOpcode(pingTypCode, PingPongCtrlCode)
	}

	return ping.opcode
}
