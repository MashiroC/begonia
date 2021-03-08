// 日志中心服务
package centerlog

import (
	"errors"
	"strconv"
)

//go:generate begonia -s -c -r ../

type CenterLogService struct {
	store *logServiceStore
}

func NewCenterLogService() *CenterLogService {
	return &CenterLogService{
		store: newStore(),
	}
}

// 注册日志服务
func (c *CenterLogService) RegisterLogService(name string) error {
	// 初始化一个日志服务
	l := newLogService()
	l.SetLogServicePath(name)
	l.Name = name
	return c.store.Add(name, l)
}

// 信息写入
func (c *CenterLogService) Write(name string, msg []byte) (string, error) {
	l, ok := c.store.Get(name)
	if !ok {
		return "", errors.New("No find server")
	}
	write, err := l.Write(msg)
	return strconv.Itoa(write), err
}

// 查看日志服务
func (c *CenterLogService)GetLog(name string)([]byte,error){
	l, ok := c.store.Get(name)
	if !ok {
		return nil, errors.New("No find server")
	}
	return l.GetLogServiceInfo()
}