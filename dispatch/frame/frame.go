// Package frame 用于在通讯层使用的frame
package frame

import (
	"errors"
	"fmt"
	"github.com/MashiroC/begonia/tool/qconv"
)

const (
	// frame中payload部分，string字段的分隔符
	breakByte = 0x00

	// CtrlDefaultCode 默认的ctrl code
	CtrlDefaultCode = 0

	// PingPongCtrlCode ping-pong的ctrl code
	PingPongCtrlCode = 7 // 0b0111

	// CtrlConnCode 将连接升级为连接池的code
	CtrlConnCode = 1

	// ProtocolVersion 默认的版本
	ProtocolVersion = 0
)

// Frame 对外暴露的接口
type Frame interface {

	// Marshal 对payload部分序列化
	Marshal() []byte

	// Opcode 组装opcode
	Opcode() int
}

// ParseOpcode 将opcode解析为typCode和ctrlCode
func ParseOpcode(opcode int) (typCode, ctrlCode int) {
	version := opcode & 15 // 15 = 0b00001111
	if version != ProtocolVersion {
		panic("协议版本不支持")
	}

	ctrlCode = opcode >> 4 & 7 // 15 = 0b0111

	typCode = opcode >> 7

	return
}

// makeOpcode 使用默认字段构建opcode
func makeOpcode(typCode int) int {
	dispatchCode := CtrlDefaultCode // 0 ~ 7

	version := ProtocolVersion // 0 ~ 15

	return ((typCode<<3)|dispatchCode)<<4 | version
}

// Unmarshal 根据typCode和序列化的数据，反序列化为frame
func Unmarshal(typCode int, data []byte) (f Frame, err error) {
	switch typCode {
	case requestTypCode:
		f, err = unMarshalRequest(data)
	case responseTypCode:
		f, err = unMarshalResponse(data)
	default:
		panic(typCode)
	}
	return
}

func makePingPongOpcode(typCode int) int {
	ctrlCode := PingPongCtrlCode // 0 ~ 7

	version := ProtocolVersion // 0 ~ 15

	return ((typCode<<3)|ctrlCode)<<4 | version
}

func UnMarshalPingPong(typCode int, data []byte) (f Frame, err error) {
	switch typCode {
	case PingTypCode:
		f, err = unMarshalPing(data)
	case PongTypCode:
		f, err = unMarshalPong(data)
	default:
		panic(typCode)
	}
	return
}

func findPosInBytes(data []byte, start int) (pos int) {
	for i := start; i < len(data); i++ {
		if data[i] == breakByte {
			return i
		}
	}

	return -1
}

func findInBytes(data []byte, pos int) (res string, endPos int, err error) {
	tmpPos := findPosInBytes(data, pos+1)

	if tmpPos == -1 {
		err = errors.New(fmt.Sprint("frame unmarshal error: ", data))
		return
	}

	res = qconv.Qb2s(data[pos+1 : tmpPos])
	endPos = tmpPos
	return
}
