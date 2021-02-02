package machine

import (
	"github.com/MashiroC/begonia/tool/chain"
	"github.com/shirou/gopsutil/mem"
	"strconv"
)

type MemMonitor struct {
	chain.BaseHandler
}

func GetMemInfo() map[string]string {
	m := make(map[string]string)
	memory, _ := mem.VirtualMemory()
	m["mem_total"] = strconv.FormatUint(memory.Total / 1024 / 1024, 10)
	m["mem_free"] = strconv.FormatUint(memory.Free / 1024 / 1024, 10)
	m["mem_used_percent"] = strconv.FormatFloat(memory.UsedPercent, 'f', 3, 64)
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
