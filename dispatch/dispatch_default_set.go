package dispatch

import (
	"fmt"
	"github.com/MashiroC/begonia/dispatch/conn"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/dispatch/heartbeat"
	"github.com/MashiroC/begonia/internal/proxy"
	"github.com/MashiroC/begonia/tool/ids"
	"log"
	"sync"
)

// dispatch_default.go something

// NewSetByDefaultCluster 在default cluster模式下创建一个dispatch
func NewSetByDefaultCluster() Dispatcher {

	d := &setDispatch{}

	d.connSet = make(map[string]conn.Conn)

	// 默认连接被关闭时只打印log
	d.CloseHookFunc = func(connID string, err error) {
		log.Printf("connID [%s] has some error: [%s]\n", connID, err)
	}

	return d
}

type setDispatch struct {
	baseDispatch

	// 代理器，如果一个节点被赋予了代理职责，会在这里检查是否要重定向
	Proxy *proxy.Handler

	// set模式相关变量
	connSet  map[string]conn.Conn // 保存连接的map
	machines map[string]map[string]string

	connLock     sync.Mutex // 锁，保证connSet线程安全
	machinesLock sync.Mutex // 保证machines的线程安全
}

// Link 建立连接，bgacenter cluster模式下，会开一条和center的tcp连接
func (d *setDispatch) Link(addr string) (err error) {
	panic("in set mode, you can't use Link()")
}

func (d *setDispatch) ReLink() bool {
	panic("in set mode, you can't use ReLink()")
}

// Send 发送一个包，在center cluster模式下直接发送到中心，中心进行调度
func (d *setDispatch) Send(f frame.Frame) (err error) {
	panic("in set mode, you can't use Send()")
}

func (d *setDispatch) HandleFrame(connID string, f frame.Frame) {

	if d.Proxy != nil {
		redirectConnID, ok := d.Proxy.Check(connID, f)
		if ok {
			d.Proxy.Action(connID, redirectConnID, f)
			return
		}
	}

	d.LgHandleFrame(connID, f)
}

func (d *setDispatch) SendTo(connID string, f frame.Frame) (err error) {
	// 及时解锁，不使用defer，避免大数据包的write协程持续占有锁
	d.connLock.Lock()
	c, ok := d.connSet[connID]
	d.connLock.Unlock()

	if !ok {
		return fmt.Errorf("dispatch send error: conn [%s] is broken or disconnection", connID)
	}

	err = c.Write(byte(f.Opcode()), f.Marshal())
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

	log.Printf("new conn addr [%s] accept, connID [%s]\n", c.Addr(), id)
	d.connLock.Lock()
	d.connSet[id] = c
	d.connLock.Unlock()
	ping := heartbeat.NewPing(7, id)
	ping.Start(d)

	for {

		opcode, data, err := c.Recv()
		if err != nil {
			c.Close()
			d.CloseHookFunc(id, err)
			d.connLock.Lock()
			delete(d.connSet, id)
			d.connLock.Unlock()
			break
		}

		// 解析opcode
		typ, ctrl := frame.ParseOpcode(int(opcode))

		switch ctrl {
		case frame.CtrlDefaultCode:
			f, err := frame.UnMarshalBasic(typ, data)
			if err != nil {
				panic(err)
			}

			go d.HandleFrame(id, f)
			//d.msgCh <- recvMsg{
			//	connID: id,
			//	f:      f,
			//}

		case frame.PingPongCtrlCode:
			f, err := frame.UnMarshalBasic(typ, data)
			if err != nil {
				panic(err)
			}
			machine := ping.HandleFrame(f)
			d.machinesLock.Lock()
			d.machines[id] = machine
			d.machinesLock.Unlock()

		default:
			panic(fmt.Sprintf("ctrl code [%s] not support", ctrl))
		}
	}

}

func (d *setDispatch) Handle(typ string, in interface{}) {
	switch typ {
	case "proxy":
		if p, ok := in.(*proxy.Handler); ok {
			d.Proxy = p
			return
		}
	default:
		d.baseDispatch.Handle(typ, in)
		return
	}
	panic("handle func not exist")
}

func (d *setDispatch) Close() {
	for _, v := range d.connSet {
		v.Close()
	}
}

func (d *setDispatch) Get(connId string) interface{} {
	d.machinesLock.Lock()
	t := d.machines[connId]
	d.machinesLock.Unlock()
	return t
}
