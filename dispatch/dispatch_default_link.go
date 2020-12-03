package dispatch

import (
	"github.com/MashiroC/begonia/config"
	"github.com/MashiroC/begonia/dispatch/conn"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/tool/berr"
	"github.com/MashiroC/begonia/tool/ids"
	"log"
	"time"
)

// dispatch_default.go something

// NewByDefaultCluster 在default cluster模式下创建一个dispatch
func NewLinkedByDefaultCluster() Dispatcher {

	d := &linkDispatch{}

	// 判断是否需要在断开连接情况下重连，hook了dispatch层的close函数
	if !config.C.Dispatch.AutoReConnection {

		// 不配置自动重连时 默认连接被关闭时只打印log
		// TODO:hook变链式执行
		d.CloseHookFunc = func(connID string, err error) {
			log.Printf("connID [%s] has some error: [%s]\n", connID, err)
		}
	} else {

		d.CloseHookFunc = func(connID string, err error) {

			// 用一个协程跑 避免阻塞
			go func() {
				ok := false

				if config.C.Dispatch.ReConnectionRetryCount <= 0 {

					for !ok {
						log.Println("connot link to server,retry...")
						time.Sleep(time.Duration(config.C.Dispatch.ReConnectionIntervalSecond) * time.Second)
						ok = d.ReLink()
					}

				} else {

					for i := 0; i < config.C.Dispatch.ReConnectionRetryCount && !ok; i++ {
						log.Println("connot link to server,retry", i, "limit", config.C.Dispatch.ReConnectionRetryCount)
						time.Sleep(time.Duration(config.C.Dispatch.ReConnectionIntervalSecond) * time.Second)
						ok = d.ReLink()
					}

					if !ok {
						panic("connect closed")
					}

				}
			}()

		}

	}

	return d
}

type linkDispatch struct {
	baseDispatch

	// link模式相关变量
	linkAddr   string    // 单连接的地址
	linkedConn conn.Conn // 连接
	linkID     string    // 连接的id
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

func (d *linkDispatch) Listen(addr string) {
	panic(berr.New("dispatch", "listen", "link mode can't use Listen()"))
}

// work 获得一个新的连接之后持续监听连接，然后把消息发送到msgCh里
func (d *linkDispatch) work(c conn.Conn) {

	id := ids.New()

	d.linkID = id
	log.Printf("link [%s] success\n", id)

	for {

		opcode, data, err := c.Recv()
		if err != nil {
			c.Close()
			d.CloseHookFunc(id, err)
			break
		}

		// 解析opcode
		typ, ctrl := frame.ParseOpcode(int(opcode))

		if ctrl == frame.CtrlDefaultCode {

			f, err := frame.UnMarshal(typ, data)
			if err != nil {
				panic(err)
			}

			//d.msgCh <- recvMsg{
			//	connID: id,
			//	f:      f,
			//}
			go d.LgHandleFrame(id,f)

		} else {
			// TODO:现在没有除了普通请求之外的ctrl code 支持
			panic(berr.NewF("dispatch", "recv", "ctrl code [%s] not support", ctrl))
		}
	}

}

func (d *linkDispatch) Close() {
	d.linkedConn.Close()
}
