package frame

import (
	"errors"
)

const (
	breakByte       = 0x00
	CtrlDefaultCode = 0
	ProtocolVersion = 1
)

var (
	unMarshalErr = errors.New("unmarshal frame error")
)

type Datas = map[string]interface{}

type Frame interface {
	Marshal() []byte
	Opcode() int
}

func ParseOpcode(opcode int) (typCode, ctrlCode int) {
	version := opcode & 0b00001111
	if version != ProtocolVersion {
		panic("协议版本不支持")
	}

	ctrlCode = opcode >> 4 & 0b0111


	typCode = opcode >> 7

	return
}

func makeOpcode(typCode int) int {
	dispatchCode := CtrlDefaultCode // 0 ~ 7

	version := ProtocolVersion // 0 ~ 15

	return ((typCode<<3)|dispatchCode)<<4 | version
}

func UnMarshal(typCode int, data []byte) (f Frame, err error) {
	switch typCode {
	case RequestTypCode:
		f, err = unMarshalRequest(data)
	case ResponseOpCode:
		f, err = unMarshalResponse(data)
	}
	return
}
