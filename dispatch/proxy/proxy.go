package proxy

import (
	"github.com/MashiroC/begonia/core"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/logic/containers"
)

type HandlerAction = func(connID, redirectConnID string, f frame.Frame)

type Handler interface {
	Check(connID string, f frame.Frame) (redirectConnID string, ok bool)
	Action(connID, redirectConnID string, f frame.Frame)
	AddAction(action HandlerAction)
}

func NewCenterProxyHandler() Handler {
	return &CenterProxyHandler{}
}

type baseHandler struct {
	handlerChains []HandlerAction
}

func (c *baseHandler) AddAction(action HandlerAction) {
	if c.handlerChains == nil {
		c.handlerChains = make([]HandlerAction, 0, 2)
	}

	c.handlerChains = append(c.handlerChains, action)
}

func (c *baseHandler) Action(connID, redirectConnID string, f frame.Frame) {
	if c.handlerChains == nil {
		return
	}

	for i := 0; i < len(c.handlerChains); i++ {
		c.handlerChains[i](connID, redirectConnID, f)
	}
}

type CenterProxyHandler struct {
	baseHandler

	waitChan containers.WaitChans
}

func (c *CenterProxyHandler) Check(connID string, f frame.Frame) (redirectConnID string, ok bool) {

	// Response不走proxy器
	if _, okk := f.(*frame.Response); okk {
		return
	}

	req := f.(*frame.Request)

	if req.Service != core.ServiceName {

		redirectConnID, ok = core.C.GetToID(req.Service)
		if !ok {
			panic("unknown bu ok error")
		}
	}
	return
}
