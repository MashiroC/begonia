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
	d.cancelSet = make(map[string]func())

	h := heartbeat.NewHeart()

	d.Hook("link", func(connID string) {
		closeFunc := func() {
			d.connLock.Lock()
			c := d.connSet[connID]
			d.connLock.Unlock()

			if c != nil {
				c.Close()
			}
		}
		sendFunc := func(connID string, f frame.Frame) error {
			return d.SendTo(connID, f)
		}

		d.cancelLock.Lock()
		d.cancelSet[connID] = h.Register("ping", connID, closeFunc, sendFunc)
		d.cancelLock.Unlock()
	})

	d.Hook("close", func(connID string, err error) {
		d.cancelLock.Lock()
		cancel, ok := d.cancelSet[connID]
		d.cancelLock.Unlock()
		if !ok || cancel == nil {
			log.Println("close hook error: cancel func not exit")
			return
		}

		cancel()

		d.cancelLock.Lock()
		delete(d.cancelSet, connID)
		d.cancelLock.Unlock()
	})

	d.Handle("ctrl", heartbeat.Handler(h))

	return d
}

type setDispatch struct {
	baseDispatch

	// set模式相关变量
	connSet  map[string]conn.Conn // 保存连接的map
	connLock sync.Mutex           // 锁，保证connSet线程安全

	cancelSet  map[string]func() // 保存取消goroutine的函数
	cancelLock sync.Mutex        // 保证并发安全
}

func (d *setDispatch) Start(addr string) (err error) {
	go d.listen(addr)
	return
}

// Send 发送一个包，在center cluster模式下直接发送到中心，中心进行调度
func (d *setDispatch) Send(f frame.Frame) (err error) {
	panic("in set mode, you can't use Send()")
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
	f.Release()
	return
}

func (d *setDispatch) listen(addr string) {

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

	d.DoLinkHook(id) // 变量初始化完成，这里去hook一些东西

	for {

		opcode, payload, err := c.Recv()
		if err != nil {
			c.Close()
			d.DoCloseHook(id, err)
			d.connLock.Lock()
			delete(d.connSet, id)
			d.connLock.Unlock()
			break
		}

		d.rt.Do(id, opcode, payload)
	}

}

func (d *setDispatch) Handle(typ string, in interface{}) {
	switch typ {
	case "proxy":
		if p, ok := in.(*proxy.Handler); ok {
			d.rt.Proxy = p
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

func (d *setDispatch) Upgrade(connID string, addr string) (err error) {
	d.connLock.Lock()
	defer d.connLock.Unlock()

	src := d.getConnID(addr)

	pool, ok := d.connSet[connID]
	if !ok {
		return fmt.Errorf("upgrade conn error: conn [%s] is broken or disconnection", connID)
	}

	c, ok := d.connSet[src]
	if !ok {
		return fmt.Errorf("upgrade conn error: conn [%s] is broken or disconnection", src)
	}

	err = conn.Join(pool, c)
	if err != nil {
		return
	}

	d.connSet[connID] = pool
	return
}

func (d *setDispatch) getConnID(addr string) (connID string) {
	for id, c := range d.connSet {
		if c.Addr() == addr {
			connID = id
			break
		}
	}
	return
}
