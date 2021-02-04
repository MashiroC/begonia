// Package logic 逻辑层
package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/MashiroC/begonia/dispatch"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/dispatch/heartbeat"
	"github.com/MashiroC/begonia/tool/ids"
	"log"
	"reflect"
	"strings"
)

// Callback logic层回调函数的alias
type Callback = func(result *CallResult)

// Client Client接口的实现结构体
type Client struct {
	Dp        dispatch.Dispatcher // dispatch层的接口，供logic层向下继续调用
	Callbacks *CallbackStore      // 回调仓库，可以在这里注册回调，调用回调
}

// NewClient 创建一个新的 logic层 客户端
func NewClient(dp dispatch.Dispatcher) *Client {

	c := &Client{
		Dp:        dp,
		Callbacks: NewWaitChans(),
	}

	dp.Handle("frame", c.DpHandler)

	return c
}

func (c *Client) DpHandler(connID string, f frame.Frame) {

	switch resp := f.(type) {
	case *frame.Response:
		c.HandleResponse(resp)
		return
	case *frame.Ping:
		heartbeat.HandlePing(resp)
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

func (c *Client) CallSync(call *Call) *CallResult {

	ch := make(chan *CallResult)
	defer close(ch)

	c.CallAsync(call, func(result *CallResult) {
		ch <- result
	})

	return <-ch
}

func (c *Client) CallAsync(call *Call, callback Callback) {

	reqID := ids.New()
	var f frame.Frame
	f = frame.NewRequest(reqID, call.Service, call.Fun, call.Param)

	c.Callbacks.AddCallback(context.TODO(), reqID, callback)

	if err := c.Dp.Send(f); err != nil {
		err = c.Callbacks.Callback(reqID, &CallResult{
			Result: nil,
			Err:    fmt.Errorf("logic call error: %w", err),
		})
		if err != nil {
			// TODO:println => errorln
			log.Println(err)
		}
	}

}

func (c *Client) Hook(typ string, hookFunc interface{}) {

	types := strings.Split(typ, ".")

	if len(types) == 2 {

		switch types[0] {
		case "dispatch":
			c.Dp.Hook(types[1], hookFunc)
			return
		}

	} else if len(types) == 1 {
		// 现在logic还没有需要hook的
	}

	panic(fmt.Sprintf("hook func [%s] not found", typ))
}

func (c *Client) Handle(typ string, handleFunc interface{}) {
	types := strings.Split(typ, ".")

	if len(types) == 2 {

		switch types[0] {
		case "dispatch":
			c.Dp.Handle(types[1], handleFunc)
			return
		}

	} else if len(types) == 1 {
		// 现在logic还没有需要handle的
	}

	panic(fmt.Sprintf("handle func [%s] not found", typ))

}

func (c *Client) Close() {
	c.Dp.Close()
}
