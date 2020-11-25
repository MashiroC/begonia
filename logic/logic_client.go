package logic

import (
	"github.com/MashiroC/begonia/dispatch"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/tool/berr"
	"log"
	"reflect"
)

// logic_client.go 客户端相关的logic层代码

// Client 对外暴露的logic层的接口
type Client interface {

	// logic 组装了基础逻辑接口
	logic

	// Handle 阻塞处理客户端接收到的包
	// Client会从dispatch中获得消息，这里收到的消息都是rpc调用中的返回结果
	// 直接根据reqID去回调即可。
	Handle()

	// Close 关闭，释放资源
	Close()
}

// client Client接口的实现结构体
type client struct {
	baseLogic // 组装了基础逻辑结构体
}

// NewClient 创建一个新的 logic层 客户端
func NewClient(dp dispatch.Dispatcher) Client {

	c := &client{
		baseLogic: baseLogic{
			dp:       dp,
			waitChan: NewWaitChans(),
		},
	}

	go c.Handle()

	return c
}

func (c *client) Handle() {

	for {

		_, f := c.dp.Recv()
		msg, ok := f.(*frame.Response)
		if !ok {
			panic(berr.NewF("logic", "handle", "msg typ must *frame.Response but %s", reflect.TypeOf(f).String()))
		}

		reqID := msg.ReqID
		err := c.waitChan.Callback(reqID, &CallResult{
			Result: msg.Result,
			Err:    berr.New("rpc", "call", msg.Err),
		})

		if err != nil {
			// TODO:Println => Errorln
			log.Println(err)
		}
	}

}

func (c *client) Close() {
	c.dp.Close()
}
