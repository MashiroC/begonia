// Time : 2020/9/19 15:55
// Author : Kieran

// center
package center

import "begonia2/logic"

// center.go something

// Center 服务中心的接口，对外统一用接口
type Center interface {
	Run(addr string)
}

type rCenter struct {
	lg logic.MixNode
}

func (c *rCenter) Run(addr string) {
	c.lg.
}
