// Time : 2020/9/30 11:41
// Author : Kieran

// dispatch
package dispatch

import (
	"begonia2/dispatch/conn"
	"begonia2/dispatch/frame"
	"begonia2/tool/ids"
	"log"
	"reflect"
	"sync"
)

// dispatch_default.go something

type dispatchMode int

const (
	linked dispatchMode = iota + 1
	set
)

func NewByCenterCluster() Dispatcher {
	d := &defaultDispatch{}
	d.msgCh = make(chan recvMsg, 10)
	d.closeHookFunc = func(connID string, err error) {
		log.Printf("connID [%s] has some error: [%s]\n", connID, err)
	}
	return d
}

type defaultDispatch struct {
	mode dispatchMode

	linkAddr   string
	linkedConn conn.Conn
	linkID     string

	connSet  map[string]conn.Conn
	connLock sync.Mutex

	msgCh chan recvMsg

	closeHookFunc func(connID string, err error)
}

func (d *defaultDispatch) Hook(typ string, hookFunc interface{}) {
	switch typ {
	case "close":
		if f, ok := hookFunc.(func(connID string, err error)); ok {
			d.closeHookFunc = f
			return
		}
		panic("hook close func error type " + reflect.TypeOf(hookFunc).String())
	default:
		panic("hook typ error: " + typ)
	}
}

type recvMsg struct {
	connID string
	f      frame.Frame
}

// Link 建立连接，center cluster模式下，会开一条和center的tcp连接
func (d *defaultDispatch) Link(addr string) (err error) {

	d.linkAddr = addr

	c, err := conn.Dial(addr)
	if err != nil {
		return
	}

	if d.mode != 0 {
		panic("mode error")
	}

	d.mode = linked
	d.linkedConn = c

	log.Println("link", addr, "success")

	go d.work(c)

	return
}

func (d *defaultDispatch) ReLink() bool {
	err := d.Link(d.linkAddr)
	if err!=nil{
		return false
	} else {
		return true
	}
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

	if d.mode == linked {
		log.Println("send to linkConn:", string(f.Marshal()))
		d.linkedConn.Write(byte(f.Opcode()), f.Marshal())
	} else {
		panic("mode err!")
	}
	return
}

func (d *defaultDispatch) SendTo(connID string, f frame.Frame) (err error) {
	var c conn.Conn
	switch d.mode {
	case linked:
		if connID != d.linkID {
			panic("connID and linkID error")
		}

		c = d.linkedConn
	case set:
		var ok bool

		d.connLock.Lock()
		c, ok = d.connSet[connID]
		d.connLock.Unlock()

		if !ok {
			panic("connID not found")
		}
	default:
		panic("mode error")
	}

	log.Println("send to", connID, "opcode:", f.Opcode())
	err = c.Write(byte(f.Opcode()), f.Marshal())
	return
}

func (d *defaultDispatch) Recv() (connID string, f frame.Frame) {
	msg := <-d.msgCh
	connID = msg.connID
	f = msg.f
	return
}

func (d *defaultDispatch) Listen(addr string) {
	d.mode = set
	d.connSet = make(map[string]conn.Conn)

	acCh, errCh := conn.Listen(addr)

out:
	for {
		select {
		case c, ok := <-acCh:
			if !ok {
				break out
			}
			go d.work(c)
		case err, ok := <-errCh:
			if !ok {
				break out
			}
			panic(err)
		}
	}

}

func (d *defaultDispatch) work(c conn.Conn) {

	id := ids.New()
	switch d.mode {
	case linked:
		d.linkID = id
	case set:
		d.connLock.Lock()
		d.connSet[id] = c
		d.connLock.Unlock()
	default:
		panic("mode error")
	}

	for {
		opcode, data, err := c.Recv()
		if err != nil {
			//TODO:handler error
			d.closeHookFunc(id, err)
			break
		}
		log.Println("recv:", opcode, string(data))
		typ, ctrl := frame.ParseOpcode(int(opcode))
		if ctrl == frame.CtrlDefaultCode {
			var err error
			f, err := frame.UnMarshal(typ, data)
			if err != nil {
				panic(err)
			}
			d.msgCh <- recvMsg{
				connID: id,
				f:      f,
			}

		} else {
			// TODO:现在没有除了普通请求之外的ctrl code 支持
			panic("ctrl code error!")
		}
	}

}

func (d *defaultDispatch) Close() {
	d.linkedConn.Close()
}
