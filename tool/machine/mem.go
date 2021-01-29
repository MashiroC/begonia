package machine

import (
	"github.com/MashiroC/begonia/tool/chain"
	"runtime"
	"strconv"
)

type MemMonitor struct {
	chain.BaseHandler
}

func GetMemInfo() map[string]string {
	m := make(map[string]string)
	m["mem"] = strconv.Itoa(runtime.MemProfileRate)
	return m
}

func NewMemMonitor() *MemMonitor {
	mm := &MemMonitor{}

	mm.SetHandleFunc(func(req *chain.Request) {
		if req.Code & 0b10 == 0 {
			return
		}

		req.Code ^= 0b10
		m := GetMemInfo()
		req.ResFun(m)
	})

	return mm
}
