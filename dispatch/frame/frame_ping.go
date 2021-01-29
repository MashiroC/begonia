package frame

import (
	"bytes"
)

const (
	// ping的typCode
	pingTypCode = 0
)

// Ping Request的frame实现
//
// opcode4 length8 extendLength16
// req:service fun reqId param
//     4      4         8       0 || 16   [length]
// {opcode}{version}{length}{extendLength}{flag}
//
type Ping struct {
	Code byte // 需要获取的机器信息

	v      []byte // 序列化后的payload，这里是一个缓存
	opcode int    // 序列化后的opcode，初始化为-1
}

// NewPing 创建一个新的Request
func NewPing(code byte) Frame {
	return &Ping{
		Code:   code,
		opcode: -1,
	}
}

// unMarshalRequest 根据payload去反序列化出一个request
func unMarshalPing(data []byte) (ping *Ping, err error) {

	ping = &Ping{}
	buf := bytes.NewBuffer(data)

	ping.Code = buf.Bytes()[0]

	ping.v = data
	ping.opcode = -1
	return
}

// Marshal 序列化payload
//
//     4      4         8       0 || 16   [       	length    	  ]
// {opcode}{version}{length}{extendLength}{PingPongTime}0x00{Param}
//
func (ping *Ping) Marshal() []byte {

	if ping.v == nil {
		buf := make([]byte, 0, 128)

		buf = append(buf, ping.Code)
		ping.v = buf
	}

	return ping.v
}

// Opcode 组装出一个opcode
func (ping *Ping) Opcode() int {
	if ping.opcode == -1 {
		ping.opcode = makeOpcode(PingPongCtrlCode)
	}

	return ping.opcode
}
