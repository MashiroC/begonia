package frame

import (
	"errors"
	"github.com/MashiroC/begonia/tool/qconv"
)

// frame_request.go something

const (
	// request的typCode
	requestTypCode = 0
)

// Request Request的frame实现
//
// opcode4 length8 extendLength16
// req:server1 fun reqId param
//     4      4         8       0 || 16   [              length                  ]
// {opcode}{version}{length}{extendLength}{reqId}0x49{server1}0x49{fun}0x49{param}
//
type Request struct {
	ReqID   string // 请求id
	Service string // 要调用的服务
	Fun     string // 要调用的函数
	Params  []byte // 入参

	Header map[string]string

	v      []byte // 序列化后的payload，这里是一个缓存
	opcode int    // 序列化后的opcode，初始化为-1
}

// NewRequest 创建一个新的Request
func NewRequest(reqID, service, fun string, params []byte) Frame {
	return &Request{
		ReqID:   reqID,
		Service: service,
		Fun:     fun,
		Params:  params,

		opcode: -1,
	}
}

// unMarshalRequest 根据payload去反序列化出一个request
func unMarshalRequest(data []byte) (req *Request, err error) {
	req = &Request{}

	var pos int
	pos = -1

	var tmp []byte
	tmp, pos, err = findInBytes(data, pos)
	if err != nil {
		return
	}

	if len(tmp) != 0 {
		req.Header, err = unMarshalHeader(tmp)
		if err != nil {
			return
		}
	}

	req.ReqID, pos, err = findInBytesString(data, pos)
	if err != nil || len(req.ReqID) == 0 {
		return
	}
	if len(req.ReqID) == 0 {
		err = errors.New("frame unmarshal error: reqID len can not be zero")
		return
	}

	req.Service, pos, err = findInBytesString(data, pos)
	if err != nil {
		return
	}
	if len(req.Service) == 0 {
		err = errors.New("frame unmarshal error: service len can not be zero")
		return
	}

	req.Fun, pos, err = findInBytesString(data, pos)
	if err != nil {
		return
	}
	if len(req.Fun) == 0 {
		err = errors.New("frame unmarshal error: fun len can not be zero")
		return
	}

	req.Params = data[pos+1:]

	req.v = data
	req.opcode = -1
	return
}

// Marshal 序列化payload
//
//      4      4         8       0 || 16   [              length                  ]
//	{opcode}{version}{length}{extendLength}{header}0x00{reqId}0x00{server1}0x00{fun}0x00{param}
//  header: {key}0x01{value}0x01{key}0x01{value}
//
func (r *Request) Marshal() []byte {

	if r.v == nil {
		buf := make([]byte, 0, 128)

		buf = append(buf, marshalHeader(r.Header)...)
		buf = append(buf, breakByte)

		buf = append(buf, qconv.Qs2b(r.ReqID)...)
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

// Opcode 组装出一个opcode
func (r *Request) Opcode() int {
	if r.opcode == -1 {
		r.opcode = makeOpcode(requestTypCode)
	}

	return r.opcode
}

func (r *Request) Release() {

}

func unMarshalHeader(data []byte) (res map[string]string, err error) {
	res = make(map[string]string)

	flag := true
	key := ""
	pos := 0
	for {
		b := data[pos]
		if b == headerBreakByte {
			tmp := data[:pos]
			if flag {
				key = qconv.Qb2s(tmp)
			} else {
				value := qconv.Qb2s(tmp)
				res[key] = value
			}
			data = data[pos+1:]
			pos = 0
			flag = !flag
		}
		if pos == len(data)-1 {
			// end
			if flag {
				// only key
				err = errors.New("header unmarshal error, header last entry don't have value")
				return
			}
			value := qconv.Qb2s(data)
			res[key] = value
			break
		}

		pos++
	}
	return
}

func marshalHeader(header map[string]string) (data []byte) {
	data = make([]byte, 0, 128)
	if header != nil && len(header) != 0 {
		for k, v := range header {
			data = append(data, qconv.Qs2b(k)...)
			data = append(data, headerBreakByte)
			data = append(data, qconv.Qs2b(v)...)
			data = append(data, headerBreakByte)
		}
		data = data[:len(data)-1]
	}

	return data
}
