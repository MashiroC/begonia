package machine

import (
	"github.com/MashiroC/begonia/tool/chain"
	"github.com/shirou/gopsutil/disk"
	"strconv"
)

type DiskMonitor struct {
	chain.BaseHandler
}

func GetDiskInfo() map[string]string {
	m := make(map[string]string)
	u, _ := disk.Usage("/")
	m["disk_used"] = strconv.FormatUint(u.Used / 1024 / 1024, 10)
	m["disk_percent"] = strconv.FormatFloat(u.UsedPercent, 'f', 3, 64)
	m["disk_total"] = strconv.FormatUint(u.Total / 1024 / 1024, 10)
	return m
}

func NewDiskMonitor() *DiskMonitor {
	cm := &DiskMonitor{}

	cm.SetHandleFunc(func(req *chain.Request) {
		if req.Code & 0b100 == 0 {
			return
		}

		req.Code ^= 0b100
		m := GetDiskInfo()
		req.ResFun(m)
	})

	return cm
}
