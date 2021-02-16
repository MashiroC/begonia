package conn

import (
	"errors"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/tool/qconv"
	"github.com/MashiroC/begonia/tool/queue"
	"sync"
	"time"
)

type info struct {
	opcode byte
	data   []byte
	err    error
}

type pool struct {
	maxPoolSize  int // 连接池拥有的最大连接数
	corePoolSize int // 连接池至少持有的连接数
	poolSize     int // 连接池当前拥有的连接数

	ttl       time.Duration // 当连接池现有连接数超过核心连接数时，多余连接能存活的最大时间
	canDial   bool
	wait      bool
	waitCh    chan struct{}
	data      *queue.Queue
	connSet   idleList
	connLock  sync.Mutex
	localAddr string
}

func (p *pool) Addr() string {
	pc := p.connSet.back.c
	return pc.Addr()
}

func (p *pool) Write(opcode byte, data []byte) (err error) {
	pc, err := p.get()
	if err != nil {
		return
	}
	err = pc.c.Write(opcode, data)
	p.put(pc)

	return
}

func (p *pool) Recv() (opcode byte, data []byte, err error) {
	for p.data.IsEmpty() {
	}

	ele := p.data.PopBack()
	info := ele.(info)
	return info.opcode, info.data, info.err
}

func (p *pool) Close() {
	p.poolSize = 0
	pc := p.connSet.front
	p.connSet.len = 0
	p.connSet.front, p.connSet.back = nil, nil
	close(p.waitCh)

	for ; pc != nil; pc = pc.next {
		pc.c.Close()
	}
}

func Upgrade(conn Conn) (Conn, error) {
	dc, ok := conn.(*defaultConn)
	if !ok {
		// 如果不能断言为defaultConn，说明已经升级过了
		return dc, errors.New("upgrade connection error: conn is already upgraded")
	}

	p := &pool{
		data: queue.New(),
		// TODO
	}

	p.poolSize++
	p.put(&poolConn{c: dc, t: time.Now()})
	go p.recv(dc)

	return p, nil
}

func Join(dst Conn, conn Conn) (c Conn, err error) {
	dc, ok := conn.(*defaultConn)
	if !ok {
		return dc, errors.New("join conn to pool error: conn is already a pool")
	}

	p, ok := dst.(*pool)
	if !ok {
		dst, _ = Upgrade(dst)
		p = dst.(*pool)
		return
	}

	p.put(&poolConn{c: dc, t: time.Now()})
	p.poolSize++
	go p.recv(dc)

	return p, nil
}

func (p *pool) dial() (conn Conn, err error){
	remoteAddr := p.Addr()
	if remoteAddr == "" {
		return
	}

	conn, err = Dial(remoteAddr)
	if err != nil {
		return
	}

	opcode := byte(makeOpcode())
	localAddr := qconv.Qs2b(p.localAddr)
	err = conn.Write(opcode, localAddr)

	return
}

func (p *pool) get() (c *poolConn, err error) {
	n := p.poolSize - p.corePoolSize
	for i := 0; i < n && p.connSet.len != 0; i++ {
		pc := p.connSet.back
		if time.Now().Sub(pc.t) < p.ttl {
			break
		}
		pc.c.Close()
		p.connSet.popBack()
		p.poolSize--
	}

	if p.wait {
		<-p.waitCh
	}

	pc := p.connSet.back
	if pc != nil {
		p.connSet.popBack()
		c = pc
		return
	}

	if !p.wait && p.poolSize >= p.maxPoolSize || !p.canDial{
		return nil, errors.New("get conn error: connection pool exhausted")
	}

	conn, err := p.dial()
	if err != nil {
		p.waitCh <- struct{}{}
		return nil, err
	}

	c = &poolConn{
		c: conn.(*defaultConn),
		t: time.Now(),
	}
	p.poolSize++

	return
}

func (p *pool) put(pc *poolConn) {
	pc.t = time.Now()
	p.connSet.pushFront(pc)

	if p.connSet.len > p.corePoolSize {
		pc = p.connSet.back
		pc.c.Close()
		p.poolSize--
		p.connSet.popBack()
	} else {
		pc = nil
	}

	if p.wait {
		p.waitCh <- struct{}{}
	}
}

func (p *pool) recv(c *defaultConn) {

	for {
		opcode, data, err := c.Recv()
		tmp := &info{opcode: opcode, data: data, err: err}
		// 把数据压入队列
		p.data.Push(tmp)
		if err != nil {
			c.Close()
			break
		}
	}
}

type idleList struct {
	len         int
	front, back *poolConn
}

type poolConn struct {
	c          *defaultConn
	t          time.Time
	next, prev *poolConn
}

func (l *idleList) pushFront(pc *poolConn) {
	pc.next = l.front
	pc.prev = nil
	if l.len == 0 {
		l.back = pc
	} else {
		l.front.prev = pc
	}
	l.front = pc
	l.len++
	return
}

func (l *idleList) popFront() {
	pc := l.front
	l.len--
	if l.len == 0 {
		l.front, l.back = nil, nil
	} else {
		pc.next.prev = nil
		l.front = pc.next
	}
	pc.next, pc.prev = nil, nil
}

func (l *idleList) popBack() {
	pc := l.back
	l.len--
	if l.len == 0 {
		l.front, l.back = nil, nil
	} else {
		pc.prev.next = nil
		l.back = pc.prev
	}
	pc.next, pc.prev = nil, nil
}

func makeOpcode() int {
	dispatchCode := frame.CtrlConnCode // 0 ~ 7

	version := frame.ProtocolVersion // 0 ~ 15

	return ((1<<3)|dispatchCode)<<4 | version
}
