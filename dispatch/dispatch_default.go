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
		//log.Println("send to linkConn:", string(f.Marshal()))
		err = d.linkedConn.Write(byte(f.Opcode()), f.Marshal())
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
			log.Printf("conn [%s] response timeout\n", connID)
			return
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
		panic("mode error")
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
			panic("ctrl code error!")
		}
	}

}

func (d *defaultDispatch) Close() {
	d.linkedConn.Close()
}
