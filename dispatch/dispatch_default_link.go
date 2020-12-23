package dispatch

import (
	"github.com/MashiroC/begonia/config"
	"github.com/MashiroC/begonia/dispatch/conn"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/tool/berr"
	"github.com/MashiroC/begonia/tool/ids"
	"github.com/MashiroC/begonia/tool/machine"
	"log"
	"reflect"
	"time"
)

// dispatch_default.go something

// NewByDefaultCluster 在default cluster模式下创建一个dispatch
func NewLinkedByDefaultCluster() Dispatcher {

	d := &linkDispatch{}

	d.msgCh = make(chan recvMsg, 10)

	// 默认连接被关闭时只打印log
	d.closeHookFunc = func(connID string, err error) {
		log.Printf("connID [%s] has some error: [%s]\n", connID, err)
	}

	return d
}

type linkDispatch struct {

	// link模式相关变量
	linkAddr   string    // 单连接的地址
	linkedConn conn.Conn // 连接
	linkID     string    // 连接的id

	msgCh chan recvMsg // 接收消息用的管道

	// hook func
	closeHookFunc func(connID string, err error) // 关闭连接的hook
}

// Hook 在这里可以去Hook一些事件。
func (d *linkDispatch) Hook(name string, hookFunc interface{}) {
	switch name {
	case "close":
		if f, ok := hookFunc.(func(connID string, err error)); ok {
			d.closeHookFunc = f
			return
		}
		panic(berr.New("dispatch", "hook", "close func must func(connID string, err error) but "+reflect.TypeOf(hookFunc).String()))
	default:
		panic(berr.New("dispatch", "hook", "hook func "+name+"not found"))
	}
}

// Link 建立连接，center cluster模式下，会开一条和center的tcp连接
func (d *linkDispatch) Link(addr string) (err error) {

	d.linkAddr = addr

	c, err := conn.Dial(addr)
	if err != nil {
		return berr.Warp("dispatch", "link", err)
	}

	d.linkedConn = c

	go d.work(c)

	return
}

func (d *linkDispatch) ReLink() bool {
	err := d.Link(d.linkAddr)
	return err == nil
}

// Send 发送一个包，在center cluster模式下直接发送到中心，中心进行调度
func (d *linkDispatch) Send(f frame.Frame) (err error) {
	// TODO:请求实现幂等 断连时排序等待连接重连 这里暂时先直接传过去
	err = d.linkedConn.Write(byte(f.Opcode()), f.Marshal())
	return
}

func (d *linkDispatch) SendTo(connID string, f frame.Frame) (err error) {
	if connID != d.linkID {
		err = berr.New("dispatch", "send", "in linked mode, you can't use SendTo() to another conn, please use Send() or passing manager center connID")
		return
	}

	err = d.linkedConn.Write(byte(f.Opcode()), f.Marshal())
	return
}

func (d *linkDispatch) Recv() (connID string, f frame.Frame) {
	msg := <-d.msgCh
	connID = msg.connID
	f = msg.f
	return
}

func (d *linkDispatch) Listen(addr string) {
	panic(berr.New("dispatch", "listen", "link mode can't use Listen()"))
}

// work 获得一个新的连接之后持续监听连接，然后把消息发送到msgCh里
func (d *linkDispatch) work(c conn.Conn) {

	id := ids.New()

	d.linkID = id
	log.Printf("link [%s] success\n", id)

	getpingtime := config.C.Dispatch.GetPingTime
	timer := time.NewTimer(getpingtime)
	go isTimeOut(timer, d)

	for {

		opcode, data, err := c.Recv()
		if err != nil {
			c.Close()
			d.closeHookFunc(id, err)
			break
		}

		// 解析opcode
		typ, ctrl := frame.ParseOpcode(int(opcode))

		switch ctrl {
		case frame.BasicCtrlCode:
			f, err := frame.UnMarshalBasic(typ, data)
			if err != nil {
				panic(err)
			}

			d.msgCh <- recvMsg{
				connID: id,
				f:      f,
			}
		case frame.PingPongCtrlCode:
			timer.Stop()
			timer.Reset(getpingtime)
			info, err := machine.M.GetMachineInfo()
			pong := frame.NewPong(info, err)

			err = d.Send(pong)

			if err != nil {
				log.Println("sendPong err", err)
			}
		default:
			panic(berr.NewF("dispatch", "recv", "ctrl code [%s] not support", ctrl))
		}
	}

}

func (d *linkDispatch) Close() {
	d.linkedConn.Close()
}

func isTimeOut(timer *time.Timer, d *linkDispatch) {
	<-timer.C
	d.Close()
}
