package chain

type BaseHandler struct {
	nextHandler Handler
	handle      func(*Request)
}

func (bh *BaseHandler) Handle(req *Request) {
	bh.handle(req)
}

func (bh *BaseHandler) NextHandler() Handler {
	return bh.nextHandler
}

func (bh *BaseHandler) SetNext(handler Handler) {
	bh.nextHandler = handler
}

func (bh *BaseHandler) SetHandleFunc(handle func(*Request)) {
	bh.handle = handle
}

func NewBaseHandler(handle func(*Request)) *BaseHandler {
	return &BaseHandler{
		handle: handle,
	}
}
