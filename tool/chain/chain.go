package chain

// 向一条责任链发起请求
type Request struct {
	Code   byte              // 请求标识码
	ResFun func(interface{}) // 用户写结果的函数
}

type Handler interface {
	Handle(req *Request)
	NextHandler() Handler
	SetNext(handler Handler)
}

type Chain struct {
	firstHandler Handler
	done         func(code byte) bool // 用于判断是否停止继续传递的函数
}

// 处理一个请求
func (c *Chain) Handle(req *Request) {
	h := c.firstHandler
	for h != nil {
		h.Handle(req)
		if c.done(req.Code) {
			break
		}
		h = h.NextHandler()
	}
}

// 注册一个实例
func (c *Chain) Sign(handler Handler) {
	handler.SetNext(c.firstHandler)
	c.firstHandler = handler
}

func NewChain() *Chain {
	return &Chain{
		firstHandler: nil,
		done: func(code byte) bool {
			return code == 0
		},
	}
}
