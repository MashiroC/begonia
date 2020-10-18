// Time : 2020/8/6 11:59
// Author : MashiroC

// frame
package frame

import (
	"begonia2/tool/qconv"
	"bytes"
	"errors"
)

// frame_request.go something

const (
	RequestTypCode = 1
)

type Request struct {
	ReqId   string
	Service string
	Fun     string
	Params  []byte

	v      []byte
	opcode int
}

func (r *Request) Marshal() []byte {
	/* opcode4 length8 extendLength16
	req:service fun reqId param
	    4      4         8       0 || 16   [              length                  ]
	{opcode}{version}{length}{extendLength}{reqId}0x49{service}0x49{fun}0x49{param}
	*/
	if r.v == nil {
		buf := make([]byte, 0, 128)

		buf = append(buf, qconv.Qs2b(r.ReqId)...)
		buf = append(buf, breakByte)

		buf = append(buf, qconv.Qs2b(r.Service)...)
		buf = append(buf, breakByte)

		buf = append(buf, qconv.Qs2b(r.Fun)...)
		buf = append(buf, breakByte)

		buf = append(buf, r.Params...)

		r.v = buf
	}

	return r.v
}

func NewRequest(reqId, service, fun string, params []byte) Frame {
	return &Request{
		ReqId:   reqId,
		Service: service,
		Fun:     fun,
		Params:  params,
		opcode:  -1,
	}
}

func (r *Request) Opcode() int {
	if r.opcode == -1 {
		r.opcode=makeOpcode(RequestTypCode)
	}

	return r.opcode
}

func unMarshalRequest(data []byte) (req *Request, err error) {
	/* opcode4 length8 extendLength16
	req:service fun reqId param
	    4      4         8       0 || 16   [              length                  ]
	{opcode}{version}{length}{extendLength}{reqId}0x49{service}0x49{fun}0x49{param}
	*/
	req = &Request{}

	buf := bytes.NewBuffer(data)
	reqIdByte, err := buf.ReadBytes(breakByte)
	if err != nil || len(reqIdByte) <= 1 {
		err = errors.New("unmarshal request reqId failed")
		return
	}
	req.ReqId = qconv.Qb2s(reqIdByte[:len(reqIdByte)-1])

	serviceByte, err := buf.ReadBytes(breakByte)
	if err != nil || len(serviceByte) <= 1 {
		err = errors.New("unmarshal request service failed")
		return
	}
	req.Service = qconv.Qb2s(serviceByte[:len(serviceByte)-1])

	funByte, err := buf.ReadBytes(breakByte)
	if err != nil || len(funByte) <= 1 {
		err = errors.New("unmarshal request fun failed")
		return
	}
	req.Fun = qconv.Qb2s(funByte[:len(funByte)-1])

	req.Params = buf.Bytes()

	return
}
