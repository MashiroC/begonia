package frame

import (
	"bytes"
	"errors"
	"github.com/MashiroC/begonia/tool/qconv"
)

// frame_response.go something

const (
	// response的typCode
	responseTypCode = 1
)

// Response response的frame实现
//
//     4      4         8       0 || 16   [     length      ]
//	{opcode}{version}{length}{extendLength}{error}0x49{param}
//
type Response struct {
	ReqID  string // 请求id
	Err    string // 调用中的错误
	Result []byte // 调用结果

	v      []byte // 原始payload
	opcode int    // opcode的缓存
}

// NewResponse 创建一个response
func NewResponse(reqID string, result []byte, err error) Frame {
	var errStr string
	if err == nil {
		errStr = ""
	} else {
		errStr = err.Error()
	}
	return &Response{
		ReqID:  reqID,
		Err:    errStr,
		Result: result,
		opcode: -1,
	}
}

// unMarshalResponse 从payload反序列化到response
func unMarshalResponse(data []byte) (resp *Response, err error) {
	resp = &Response{}

	buf := bytes.NewBuffer(data)

	reqIDBytes, err := buf.ReadBytes(breakByte)
	if err != nil || len(reqIDBytes) <= 1 {
		err = errors.New("frame unmarshal error: response reqId failed")
		return
	}
	resp.ReqID = qconv.Qb2s(reqIDBytes[:len(reqIDBytes)-1])

	respErrByte, err := buf.ReadBytes(breakByte)
	if err != nil {
		err = errors.New("frame unmarshal error: response error failed")
		return
	}
	resp.Err = qconv.Qb2s(respErrByte[:len(respErrByte)-1])

	resp.Result = buf.Bytes()

	resp.v = data
	resp.opcode = -1

	return
}

// Marshal 序列化
func (r *Response) Marshal() []byte {
	if r.v == nil {
		buf := make([]byte, 0, 128)

		buf = append(buf, qconv.Qs2b(r.ReqID)...)
		buf = append(buf, breakByte)

		buf = append(buf, qconv.Qs2b(r.Err)...)
		buf = append(buf, breakByte)

		buf = append(buf, r.Result...)

		r.v = buf
	}

	return r.v
}

// Opcode 组装opcode
func (r *Response) Opcode() int {
	if r.opcode == -1 {
		r.opcode = makeOpcode(responseTypCode, BasicCtrlCode)
	}

	return r.opcode
}
