// Time : 2020/9/30 11:41
// Author : Kieran

// dispatch
package dispatch

import (
	"begonia2/dispatch/conn"
	"begonia2/dispatch/frame"
	"log"
)

// dispatch_default.go something

func NewCenterCluster() Dispatcher {
	return &defaultDispatch{}
}

type defaultDispatch struct {
	c conn.Conn
}

// Link 建立连接，center cluster模式下，会开一条和center的tcp连接
func (d *defaultDispatch) Link(addr string) {
	c, err := conn.Dial(addr)
	if err != nil {
		// TODO:handle err
		panic(err)
	}

	d.c = c
}

// Send 发送一个包，在center cluster模式下直接发送到中心，中心进行调度
func (d *defaultDispatch) Send(f frame.Frame) (err error) {
	/* opcode4 length8 extendLength16
	req:service fun reqId param
	    4      4         8       0 || 16   [              length                  ]
	{opcode}{version}{length}{extendLength}{reqId}0x49{service}0x49{fun}0x49{param}

	resp:reqId,error,data

	{opcode}{length}{extendLength}{reqId}{error}{data}
	*/
	// TODO:请求实现幂等 断连时排序等待连接重连 这里暂时先直接传过去

	d.c.Write(byte(f.Opcode()), f.Marshal())
	return
}

func (d *defaultDispatch) Recv() (f frame.Frame) {

	for f == nil {
		opcode, data, err := d.c.Recv()
		if err != nil {
			//TODO:handler error
			log.Println(err)
			continue
		}

		typ, ctrl := frame.ParseOpcode(int(opcode))
		//TODO:handler error

		if ctrl == frame.CtrlDefaultCode {
			var err error
			f, err = frame.UnMarshal(typ, data)
			if err != nil {
				panic(err)
			}
		} else {
			// TODO:现在没有除了普通请求之外的ctrl code 支持
			panic("ctrl code error!")
		}
	}

	return
}

func (d *defaultDispatch) Close() {
	d.c.Close()
}
