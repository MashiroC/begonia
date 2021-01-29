package proxy

import (
	"github.com/MashiroC/begonia/dispatch/frame"
)

type HandlerAction = func(connID, redirectConnID string, f frame.Frame)

type CheckFunc = func(connID string, f frame.Frame) (redirectConnID string, ok bool)

func NewHandler() *Handler {
	return &Handler{}
}

type Handler struct {
	Check CheckFunc

	handlerChains []HandlerAction
}

func (c *Handler) AddAction(action HandlerAction) {
	if c.handlerChains == nil {
		c.handlerChains = make([]HandlerAction, 0, 2)
	}

	c.handlerChains = append(c.handlerChains, action)
}

func (c *Handler) Action(connID, redirectConnID string, f frame.Frame) {
	if c.handlerChains == nil {
		return
	}

	for i := 0; i < len(c.handlerChains); i++ {
		c.handlerChains[i](connID, redirectConnID, f)
	}
}