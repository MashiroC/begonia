// Package frame 用于在通讯层使用的frame
package frame

const (
	// frame中payload部分，string字段的分隔符
	breakByte = 0x00

	// CtrlDefaultCode 默认的ctrl code
	CtrlDefaultCode = 0

	// PingPongCtrlCode ping-pong的ctrl code
	PingPongCtrlCode = 7 // 0b0111

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
func makeOpcode(typCode int, ctrlCode int) int {
	// dispatchCode := CtrlDefaultCode // 0 ~ 7

	version := ProtocolVersion // 0 ~ 15

	return ((typCode<<3)|ctrlCode)<<4 | version
}

// UnMarshalBasic 根据typCode和序列化的数据，反序列化为frame
func UnMarshalBasic(typCode int, data []byte) (f Frame, err error) {
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

func UnMarshalPingPong(typCode int, data []byte) (f Frame, err error) {
	switch typCode {
	case pingTypCode:
		f, err = unMarshalPing(data)
	case pongTypCode:
		f, err = unMarshalPong(data)
	default:
		panic(typCode)
	}
	return
}
