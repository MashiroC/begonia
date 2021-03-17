package machine

import (
	"github.com/MashiroC/begonia/tool/chain"
	"github.com/shirou/gopsutil/cpu"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type CpuMonitor struct {
	chain.BaseHandler
	sync.RWMutex // 保证percent的并发安全
	percent float64 // cpu使用率
}

func (c *CpuMonitor) MonitorUsePercent() {
	for {
		percents, _ := cpu.Percent(time.Second, false)
		c.Lock()
		c.percent = percents[0]
		c.Unlock()
	}
}

// 字段名	  	含义
// cpu_proc	  	cpu核数
// cpu_percent	cpu使用率
func GetCpuInfo() map[string]string {
	m := make(map[string]string)
	m["cpu_proc"] = strconv.Itoa(runtime.GOMAXPROCS(0))
	return m
}

func NewCpuMonitor() *CpuMonitor {
	cm := &CpuMonitor{}
	go cm.MonitorUsePercent()
	cm.SetHandleFunc(func(req *chain.Request) {
		if req.Code&0b1 == 0 {
			return
		}

		req.Code ^= 0b1
		m := GetCpuInfo()
		cm.RLock()
		m["cpu_percent"] = strconv.FormatFloat(cm.percent, 'f', 3, 64)
		cm.RUnlock()
		req.ResFun(m)
	})

	return cm
}
