package machine

import (
	"github.com/MashiroC/begonia/tool/chain"
	"runtime"
	"strconv"
)

type CpuMonitor struct {
	chain.BaseHandler
}

func GetCpuInfo() map[string]string {
	m := make(map[string]string)
	m["cpu"] = strconv.Itoa(runtime.GOMAXPROCS(0))
	return m
}

func NewCpuMonitor() *CpuMonitor {
	cm := &CpuMonitor{}

	cm.SetHandleFunc(func(req *chain.Request) {
		if req.Code & 0b1 == 0 {
			return
		}

		req.Code ^= 0b1
		m := GetCpuInfo()
		req.ResFun(m)
	})

	return cm
}