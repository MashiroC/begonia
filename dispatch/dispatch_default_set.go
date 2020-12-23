package dispatch

import (
	"github.com/MashiroC/begonia/config"
	"github.com/MashiroC/begonia/dispatch/conn"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/tool/berr"
	"github.com/MashiroC/begonia/tool/ids"
	"log"
	"reflect"
	"sync"
	"time"
)

// dispatch_default.go something

// NewSetByDefaultCluster 在default cluster模式下创建一个dispatch
func NewSetByDefaultCluster() Dispatcher {

	d := &setDispatch{}

	d.msgCh = make(chan recvMsg, 10)
	d.connSet = make(map[string]conn.Conn)

	// 默认连接被关闭时只打印log
	d.closeHookFunc = func(connID string, err error) {
		log.Printf("connID [%s] has some error: [%s]\n", connID, err)
	}

	return d
}

type setDispatch struct {
	Machine

	// set模式相关变量
	connSet  map[string]conn.Conn // 保存连接的map
	connLock sync.Mutex           // 锁，保证connSet线程安全

	msgCh chan recvMsg // 接收消息用的管道

	// hook func
	closeHookFunc func(connID string, err error) // 关闭连接的hook
}

// Hook 在这里可以去Hook一些事件。
func (d *setDispatch) Hook(name string, hookFunc interface{}) {
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
func (d *setDispatch) Link(addr string) (err error) {
	panic(berr.New("dispatch", "link", "in set mode, you can't use Link()"))
}

func (d *setDispatch) ReLink() bool {
	panic(berr.New("dispatch", "link", "in set mode, you can't use ReLink()"))
}

// Send 发送一个包，在center cluster模式下直接发送到中心，中心进行调度
func (d *setDispatch) Send(f frame.Frame) (err error) {
	panic(berr.New("dispatch", "send", "in set mode, you can't use Send(), please use SendTo()"))
}

func (d *setDispatch) SendTo(connID string, f frame.Frame) (err error) {
	// 及时解锁，不使用defer，避免大数据包的write协程持续占有锁
	d.connLock.Lock()
	c, ok := d.connSet[connID]
	d.connLock.Unlock()

	if !ok {
		return berr.NewF("dispatch", "send", "conn [%s] is broked or disconnection", connID)
	}

	err = c.Write(byte(f.Opcode()), f.Marshal())
	return
}

func (d *setDispatch) Recv() (connID string, f frame.Frame) {
	msg := <-d.msgCh
	connID = msg.connID
	f = msg.f
	return
}

func (d *setDispatch) Listen(addr string) {

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
func (d *setDispatch) work(c conn.Conn) {

	// TODO:等hook的链做好之后，这里直接把不同的地方加到hook链上，就可以抽出来一个baseDefaultDispatch了

	id := ids.New()

	log.Printf("new conn [%s]\n", id)
	d.connLock.Lock()
	d.connSet[id] = c
	d.connLock.Unlock()

	getpongtime := config.C.Dispatch.GetPongTime
	timer := time.NewTimer(getpongtime)
	go receivePong(c, timer)
	go sendPing(d, config.C.Dispatch.SendPingTime)

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
			f, err := frame.UnMarshalPingPong(typ, data)
			if err != nil {
				panic(err)
			}
			switch p := f.(type) {
			case *frame.Ping:
				//center暂时无ping
			case *frame.Pong:
				timer.Stop()
				timer.Reset(getpongtime)
				d.StoreMachine(p.Machine)
			}
		default:
			// TODO:现在没有除了普通请求之外的ctrl code 支持
			panic(berr.NewF("dispatch", "recv", "ctrl code [%s] not support", ctrl))
		}

	}

}

func (d *setDispatch) Close() {
	for _, v := range d.connSet {
		v.Close()
	}
}

func sendPing(d *setDispatch, sendpingtime time.Duration) {
	ticker := time.NewTicker(sendpingtime)
	ping := frame.NewPing(0)
	for {
		<-ticker.C
		err := d.Send(ping)
		if err != nil {
			log.Println("sendPing err", err)
			return
		}
	}
}
func receivePong(c conn.Conn, timer *time.Timer) {
	<-timer.C
	c.Close()
}
