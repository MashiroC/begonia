package logic

import (
	"github.com/MashiroC/begonia/dispatch"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/tool/berr"
	"log"
	"reflect"
)

// logic_client.go 客户端相关的logic层代码

// Client Client接口的实现结构体
type Client struct {
	baseLogic // 组装了基础逻辑结构体
}

// NewClient 创建一个新的 logic层 客户端
func NewClient(dp dispatch.Dispatcher) *Client {

	c := &Client{
		baseLogic: baseLogic{
			dp:       dp,
			waitChan: NewWaitChans(),
		},
	}

	return c
}

func (c *Client) DpHandler(connID string,f frame.Frame) {

	if resp, ok := f.(*frame.Response); ok {
		c.HandleResponse(resp)
		return
	}

	panic(berr.NewF("logic", "handle", "msg typ must *frame.Response but %s", reflect.TypeOf(f).String()))

}

func (c *Client) HandleResponse(resp *frame.Response) {
	reqID := resp.ReqID
	err := c.waitChan.Callback(reqID, &CallResult{
		Result: resp.Result,
		Err:    berr.New("rpc", "call", resp.Err),
	})

	if err != nil {
		// TODO:Println => Errorln
		log.Println(err)
	}
}

func (c *Client) Close() {
	c.dp.Close()
}
