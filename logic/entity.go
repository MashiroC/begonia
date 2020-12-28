package logic

import "github.com/MashiroC/begonia/dispatch/frame"

// entity.go logic层通用的一些结构体

type Calls interface {
	Frame(reqID string) frame.Frame
}

// Call rpc调用的请求
// 一般由api层组装
type Call struct {
	Service string // 调用的服务名
	Fun     string // 调用的函数名
	Param   []byte // 远程函数的入参(已序列化)
}

func (c *Call) Frame(reqID string) frame.Frame {
	return frame.NewRequest(reqID, c.Service, c.Fun, c.Param)
}

// CallResult rpc调用的响应
// 一般在logic层封装传递给api层
// 如果rpc调用过程中发生错误， 包括远程函数主动返回的错误，框架层面的错误时，
// [Err]变量不为空。正常情况下为空字符串
type CallResult struct {
	Result []byte // rpc调用的结果，远程函数的出参(已序列化)
	Err    error  // 错误
}

func (c *CallResult) Frame(reqID string) frame.Frame {
	return frame.NewResponse(reqID, c.Result, c.Err)
}

// Redirect 结果重定向
// 在api层需要将收到的请求转发到另一条连接时，传递该变量，logic层会转发请求
var Redirect = &CallResult{
	Result: nil,
	Err:    nil,
}

// ResultFunc 回传结果的结构体
// 用于api层接收消息后，需要返回结果时调用其中的Result函数
type ResultFunc struct {

	// Result 返回结果的函数
	// 如果传toConnID，result会被转递给toConnID代表的连接
	// 如果不传，result会被传递给请求的连接
	Result func(result Calls)

	// ConnID 请求的连接id
	ConnID string

	// ReqID 请求id
	ReqID string
}
