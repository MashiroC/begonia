package logic

import (
	"github.com/MashiroC/begonia/config"
	"github.com/MashiroC/begonia/dispatch"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/tool/berr"
	"log"
	"time"
)

// logic_service.go service节点的logic层

// Service 对api层暴露的logic Service 接口
type Service interface {

	// logic 基础的logic层
	logic

	// RecvMsg 获取一个调用请求，wf是回传数据的结构体
	RecvCall() (msg *Call, wf ResultFunc)
}

// NewService 创建一个实例
func NewService(dp dispatch.Dispatcher) Service {

	c := &service{
		baseLogic: baseLogic{
			dp:       dp,
			waitChan: NewWaitChans(),
		},
	}

	// 判断是否需要在断开连接情况下重连，hook了dispatch层的close函数
	if config.C.Logic.AutoReConnection {

		c.dp.Hook("close", func(connID string, err error) {

			ok := false

			if config.C.Logic.ReConnectionRetryCount <= 0 {

				for !ok {
					log.Println("connot link to server,retry...")
					time.Sleep(time.Duration(config.C.Logic.ReConnectionIntervalSecond) * time.Second)
					ok = dp.ReLink()
				}

			} else {

				for i := 0; i < config.C.Logic.ReConnectionRetryCount && !ok; i++ {
					log.Println("connot link to server,retry", i, "limit", config.C.Logic.ReConnectionRetryCount)
					time.Sleep(time.Duration(config.C.Logic.ReConnectionIntervalSecond) * time.Second)
					ok = dp.ReLink()
				}

				if !ok {
					panic("connect closed")
				}

			}

		})

	}

	return c
}

type service struct {
	baseLogic
}

func (c *service) RecvCall() (call *Call, wf ResultFunc) {

	for {

		connID, f := c.dp.Recv()

		switch msg := f.(type) {
		case *frame.Request:

			call = &Call{
				Service: msg.Service,
				Fun:     msg.Fun,
				Param:   msg.Params,
			}

			wf = ResultFunc{
				Result: func(result *CallResult, toConnID ...string) {
					resp := frame.NewResponse(msg.ReqID, result.Result, result.Err.Error())
					if toConnID != nil {
						for _, connID := range toConnID {
							c.dp.SendTo(connID, resp)
						}
					} else {
						c.dp.Send(resp)
					}
				},
				ConnID: connID,
				ReqID:  msg.ReqID,
			}

			return

		case *frame.Response:

			err := c.waitChan.Callback(msg.ReqID, &CallResult{
				Result: msg.Result,
				Err:    berr.New("rpc", "call", msg.Err),
			})
			if err != nil {
				panic(err)
			}

		default:
			panic("err?")
		}

	}

}
