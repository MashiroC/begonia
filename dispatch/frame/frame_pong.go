package frame

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/MashiroC/begonia/tool/qconv"
)

const (
	// pong的typCode
	pongTypCode = 1
)

// Pong response的frame实现
//
//     4      4         8       0 || 16   [     length      ]
//	{opcode}{version}{length}{extendLength}{error}0x49{param}
//
type Pong struct {
	Machine map[string]string
	Err     string

	v      []byte // 原始payload
	opcode int    // opcode的缓存
}

// NewPong 创建一个response
func NewPong(mach map[string]string, err error) Frame {
	var errStr string
	if err == nil {
		errStr = ""
	} else {
		errStr = err.Error()
	}
	return &Pong{
		Err:     errStr,
		Machine: mach,
		opcode:  -1,
	}
}

// unMarshalPong 从payload反序列化到response
func unMarshalPong(data []byte) (resp *Pong, err error) {
	resp = &Pong{}

	buf := bytes.NewBuffer(data)

	respErrByte, err := buf.ReadBytes(breakByte)
	if err != nil {
		err = errors.New("frame unmarshal error: response error failed")
		return
	}
	resp.Err = qconv.Qb2s(respErrByte[:len(respErrByte)-1])

	ma := buf.Bytes()
	var mach map[string]string
	err = json.Unmarshal(ma, &mach)
	if err != nil {
		err = errors.New("frame unmarshal error: response result failed")
		return
	}
	resp.Machine = mach

	resp.v = data
	resp.opcode = -1

	return
}

// Marshal 序列化
func (r *Pong) Marshal() []byte {
	if r.v == nil {
		buf := make([]byte, 0, 128)

		buf = append(buf, qconv.Qs2b(r.Err)...)
		buf = append(buf, breakByte)

		mach, _ := json.Marshal(r.Machine)
		buf = append(buf, mach...)

		r.v = buf
	}

	return r.v
}

// Opcode 组装opcode
func (r *Pong) Opcode() int {
	if r.opcode == -1 {
		r.opcode = makeOpcode(pongTypCode, PingPongCtrlCode)
	}

	return r.opcode
}
