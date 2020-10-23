// Package ids 寓意是id(s)，是生成唯一id的工具包
package ids

import uuid "github.com/satori/go.uuid"

// New 创建一个新id
func New() string {
	return uuid.NewV4().String()
}
