package dispatch

import (
	"github.com/MashiroC/begonia/dispatch/conn"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/tool/berr"
	"github.com/MashiroC/begonia/tool/ids"
	"log"
	"reflect"
	"sync"
)

// dispatch_default.go something

// dispatchMode 该dispatch的模式，是单连接还是多连接
type dispatchMode int

const (
	linked dispatchMode = iota + 1 // 单连接
	set                            // 多链接
)

// NewByDefaultCluster 在default cluster模式下创建一个dispatch
func NewByDefaultCluster() Dispatcher {

	d := &defaultDispatch{}

	d.msgCh = make(chan recvMsg, 10)

	// 默认连接被关闭时只打印log
	d.closeHookFunc = func(connID string, err error) {
		log.Printf("connID [%s] has some error: [%s]\n", connID, err)
	}

	return d
}

type defaultDispatch struct {

	// mode 该dispatch的模式
	mode dispatchMode

	// link模式相关变量
	linkAddr   string    // 单连接的地址
	linkedConn conn.Conn // 连接
	linkID     string    // 连接的id

	// set模式相关变量
	connSet  map[string]conn.Conn // 保存连接的map
	connLock sync.Mutex           // 锁，保证connSet线程安全

	msgCh chan recvMsg // 接收消息用的管道

	// hook func
	closeHookFunc func(connID string, err error) // 关闭连接的hook
}

// Hook 在这里可以去Hook一些事件。
func (d *defaultDispatch) Hook(name string, hookFunc interface{}) {
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

type recvMsg struct {
	connID string
	f      frame.Frame
}

// Link 建立连接，center cluster模式下，会开一条和center的tcp连接
func (d *defaultDispatch) Link(addr string) (err error) {

	d.linkAddr = addr

	c, err := conn.Dial(addr)
	if err != nil {
		return berr.Warp("dispatch", "link", err)
	}

	if d.mode != 0 {
		return berr.New("dispatch", "link", "mode must not zero")
	}

	d.mode = linked
	d.linkedConn = c

	go d.work(c)

	return
}

func (d *defaultDispatch) ReLink() bool {
	err := d.Link(d.linkAddr)
	return err == nil
}

// Send 发送一个包，在center cluster模式下直接发送到中心，中心进行调度
func (d *defaultDispatch) Send(f frame.Frame) (err error) {

	// TODO:请求实现幂等 断连时排序等待连接重连 这里暂时先直接传过去
	if d.mode == linked {
		err = d.linkedConn.Write(byte(f.Opcode()), f.Marshal())
	} else {
		err = berr.New("dispatch", "send", "in set mode, you can't use Send(), please use SendTo()")
	}

	return
}

func (d *defaultDispatch) SendTo(connID string, f frame.Frame) (err error) {
	var c conn.Conn
	switch d.mode {
	case linked:
		if connID != d.linkID {
			err = berr.New("dispatch", "send", "in linked mode, you can't use SendTo() to another conn, please use Send() or passing manager center connID")
			return
		}

		c = d.linkedConn
	case set:
		var ok bool

		d.connLock.Lock()
		c, ok = d.connSet[connID]
		d.connLock.Unlock()

		if !ok {
			return berr.NewF("dispatch", "send", "conn [%s] is broked or disconnection", connID)
		}
	default:
		panic(berr.NewF("dispatch", "mode", "mode [%s] not support", d.mode))
	}

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
			//TODO: println更换errorln
			log.Println("dispatch listen error:", err.Error())
		}
	}

}

// work 获得一个新的连接之后持续监听连接，然后把消息发送到msgCh里
func (d *defaultDispatch) work(c conn.Conn) {

	id := ids.New()

	switch d.mode {
	case linked:
		d.linkID = id
		log.Printf("link [%s] success\n", id)
	case set:
		log.Printf("new conn [%s]\n", id)
		d.connLock.Lock()
		d.connSet[id] = c
		d.connLock.Unlock()
	default:
		panic(berr.NewF("dispatch", "mode", "mode [%s] not support", d.mode))
	}

	for {

		opcode, data, err := c.Recv()
		if err != nil {
			c.Close()
			d.closeHookFunc(id, err)
			d.connLock.Lock()
			delete(d.connSet, id)
			d.connLock.Unlock()
			break
		}

		// 解析opcode
		typ, ctrl := frame.ParseOpcode(int(opcode))

		if ctrl == frame.CtrlDefaultCode {

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
			panic(berr.NewF("dispatch", "recv", "ctrl code [%s] not support", ctrl))
		}
	}

}

func (d *defaultDispatch) Close() {
	d.linkedConn.Close()
}
