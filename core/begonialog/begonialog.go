// 日志中心服务
package centerlog

import (
	"context"
	"errors"
	"github.com/olivere/elastic/v7"
	"log"
)

var clientEs *elastic.Client
var ctx = context.Background()
var ErrNotFound = errors.New("server msg error")

// 日志：[时间] [title] [fields] \n callers
// 储存信息
type Msg struct {
	ServerName string            `json:"server_name"` // 服务名
	Level      int               `json:"level"`       // 日志等级
	Fields     map[string]string `json:"fields"`   // 日志信息
	Time       int64             `json:"time"`        // 时间
	Callers    []string          `json:"callers"`     // 路径
}

type CenterLogService struct {
}

// 日志信息写入
// @params
// name: 服务名
// fields: log的字段(包括msg)
func (c *CenterLogService) Save( msg []Msg) error {
	// 调用Es写入
	go c.putEs(msg)
	return nil
}

func (c *CenterLogService) putEs( msg []Msg) {
	for _, v := range msg {
		if err := v.PutMsg(); err != nil {
			log.Println(err)
			return
		}
	}
	return
}

func NewCenterLogService() *CenterLogService {
	return &CenterLogService{
	}
}

//TODO:
// 1. 远程调用
// 2. defer
// 3. 定时 | 条数 写入
// 4. 一定等级马上写入
// 5. 代码使用
// 6. 完善es最好相当与一个tool PUT DELETE GET POST（Index -> 库， Type -> 表）
// 7. 链路追踪，终究还是得有个query和list
// 8. Web

//go:generate begonia -s -c -r ../