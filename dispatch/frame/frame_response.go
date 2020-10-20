// Time : 2020/8/6 12:31
// Author : MashiroC

// frame
package frame

import (
	"begonia2/tool/qconv"
	"bytes"
	"errors"
)

// frame_response.go something

const (
	ResponseTypCode = 1
)

type Response struct {
	ReqId  string
	Err    string
	Result []byte

	opcode int
	v      []byte
}

func (r *Response) Marshal() []byte {
	/* opcode4 length8 extendLength16
	req:service fun reqId param
	    4      4         8       0 || 16   [              length                  ]
	{opcode}{version}{length}{extendLength}{reqId}0x49{service}0x49{fun}0x49{param}

	resp:reqId,error,data

	{opcode}{length}{extendLength}{reqId}{error}{data}
	*/
	if r.v == nil {
		buf := make([]byte, 0, 128)

		buf = append(buf, qconv.Qs2b(r.ReqId)...)
		buf = append(buf, breakByte)

		buf = append(buf, qconv.Qs2b(r.Err)...)
		buf = append(buf, breakByte)

		buf = append(buf, r.Result...)

		r.v = buf
	}

	return r.v
}

func (r *Response) Opcode() int {
	if r.opcode == -1 {
		r.opcode = makeOpcode(ResponseTypCode)
	}

	return r.opcode
}

func NewResponse(reqId string, result []byte, err string) Frame {
	return &Response{
		ReqId:  reqId,
		Err:    err,
		Result: result,
		opcode: -1,
	}
}

func unMarshalResponse(data []byte) (resp *Response, err error) {
	/* opcode4 length8 extendLength16
	req:service fun reqId param
	    4      4         8       0 || 16   [              length                  ]
	{opcode}{version}{length}{extendLength}{reqId}0x49{service}0x49{fun}0x49{param}

	resp:reqId,error,data

	{opcode}{length}{extendLength}{reqId}{error}{data}
	*/
	resp = &Response{}

	buf := bytes.NewBuffer(data)

	reqIdByte, err := buf.ReadBytes(breakByte)
	if err != nil || len(reqIdByte) <= 1 {
		err = errors.New("unmarshal response reqId failed")
		return
	}
	resp.ReqId = qconv.Qb2s(reqIdByte[:len(reqIdByte)-1])

	respErrByte, err := buf.ReadBytes(breakByte)
	if err != nil {
		err = errors.New("unmarshal response error failed")
		return
	}
	resp.Err = qconv.Qb2s(respErrByte[:len(respErrByte)-1])

	resp.Result = buf.Bytes()

	resp.v = data
	resp.opcode = -1

	return
}
