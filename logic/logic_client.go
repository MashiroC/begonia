package logic

import (
	"errors"
	"fmt"
	"github.com/MashiroC/begonia/dispatch"
	"github.com/MashiroC/begonia/dispatch/frame"
	"log"
	"reflect"
)

// logic_client.go 客户端相关的logic层代码

// Client Client接口的实现结构体
type Client struct {
	BaseLogic // 组装了基础逻辑结构体
}

// NewClient 创建一个新的 logic层 客户端
func NewClient(dp dispatch.Dispatcher) *Client {

	c := &Client{
		BaseLogic: BaseLogic{
			Dp:        dp,
			Callbacks: NewWaitChans(),
		},
	}

	return c
}

func (c *Client) DpHandler(connID string, f frame.Frame) {

	if resp, ok := f.(*frame.Response); ok {
		c.HandleResponse(resp)
		return
	}

	panic(fmt.Sprintf("logic handle error: msg typ must *frame.Response but %s", reflect.TypeOf(f)))
}

func (c *Client) HandleResponse(resp *frame.Response) {
	reqID := resp.ReqID

	var err error
	if resp.Err != "" {
		err = errors.New(resp.Err)
	}
	err = c.Callbacks.Callback(reqID, &CallResult{
		Result: resp.Result,
		Err:    err,
	})

	if err != nil {
		// TODO:Println => Errorln
		log.Println(err)
	}
}

func (c *Client) Close() {
	c.Dp.Close()
}
