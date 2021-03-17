package machine

import (
	"github.com/MashiroC/begonia/tool/chain"
	"github.com/shirou/gopsutil/disk"
	"strconv"
)

type DiskMonitor struct {
	chain.BaseHandler
}

// 默认路径为：/(win10下如果硬盘分区，则为所在硬盘的根目录）
// 字段名				含义
// disk_used			硬盘已使用
// disk_free			硬盘可用
// disk_total			硬盘总大小
// disk_used_percent	硬盘使用率（三位小数）
func GetDiskInfo() map[string]string {
	m := make(map[string]string)
	u, _ := disk.Usage("/")
	m["disk_used"] = strconv.FormatUint(u.Used/1024/1024, 10)
	m["disk_free"] = strconv.FormatUint(u.Free/1024/1024, 10)
	m["disk_used_percent"] = strconv.FormatFloat(u.UsedPercent, 'f', 3, 64)
	m["disk_total"] = strconv.FormatUint(u.Total/1024/1024, 10)
	return m
}

func NewDiskMonitor() *DiskMonitor {
	cm := &DiskMonitor{}

	cm.SetHandleFunc(func(req *chain.Request) {
		if req.Code&0b100 == 0 {
			return
		}

		req.Code ^= 0b100
		m := GetDiskInfo()
		req.ResFun(m)
	})

	return cm
}
