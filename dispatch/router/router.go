package router

import (
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/internal/proxy"
	"github.com/MashiroC/begonia/tool/qconv"
	"log"
)

type CtrlHandleFunc = func(connID string, typ int, data []byte)

type Router struct {
	// handle func
	LgHandleFrame func(connID string, f frame.Frame)

	// 代理器，如果一个节点被赋予了代理职责，会在这里检查是否要重定向
	Proxy *proxy.Handler

	ctrlRoute []CtrlHandleFunc
}

func New() *Router {
	return &Router{}
}

func (r *Router) AddCtrlHandle(code int, f CtrlHandleFunc) {
	if code > 7 || code < 0 {
		panic("router add error: ctrl code must < 7 but " + qconv.I2S(code))
	}

	if r.ctrlRoute == nil {
		r.ctrlRoute = make([]CtrlHandleFunc, 8)
	}

	r.ctrlRoute[code] = f
}

func (r *Router) Do(connID string, opcode byte, payload []byte) {
	// 解析opcode
	typ, ctrl := frame.ParseOpcode(int(opcode))

	if ctrl == frame.CtrlDefaultCode {
		f, err := frame.Unmarshal(typ, payload)
		if err != nil {
			panic(err)
		}

		if r.Proxy != nil {
			redirectConnID, ok := r.Proxy.Check(connID, f)
			if ok {
				r.Proxy.Action(connID, redirectConnID, f)
				return
			}
		}

		go r.LgHandleFrame(connID, f)

	} else {

		if ctrl > 7 || ctrl < 0 {
			log.Println("recv bad frame,ctrl code [" + qconv.I2S(ctrl) + "]")
			return
		}

		f := r.ctrlRoute[ctrl]
		if f == nil {
			log.Println("recv bad frame,ctrl code [" + qconv.I2S(ctrl) + "]")
			return
		}

		f(connID, typ, payload)
	}
}
