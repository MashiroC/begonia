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
	maxPoolSize  int           // 连接池拥有的最大连接数
	corePoolSize int           // 连接池至少持有的连接数
	poolSize     int           // 连接池当前拥有的连接数
	canDial      bool          // 当连接不够时，是否通过dial获取连接
	wait         bool          // 当连接数超过maxPoolSize时，是否需要排队
	waitCh       chan struct{} // 通过channel实现排队
	signal       chan struct{} // 通知Recv解除阻塞状态的信号
	ttl          time.Duration // 当连接池现有连接数超过核心连接数时，多余连接能存活的最大时间

	connSet    idleList
	data       *queue.Queue
	dataLock   sync.Mutex
	mu         sync.Mutex
	localAddr  string
	remoteAddr string
}

func (p *pool) Addr() string {
	return p.remoteAddr
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

begin:
	p.dataLock.Lock()

	// 当没有接收到数据时，进入阻塞状态
	if p.data.IsEmpty() {
		p.dataLock.Unlock()
		<-p.signal
		p.dataLock.Lock()
	}

	ele := p.data.PopBack()
	p.dataLock.Unlock()

	info, ok := ele.(*info)
	if !ok {
		goto begin
	}

	return info.opcode, info.data, info.err
}

func (p *pool) Close() {
	p.mu.Lock()

	p.poolSize = 0
	p.connSet.len = 0
	pc := p.connSet.front
	p.connSet.front, p.connSet.back = nil, nil
	close(p.waitCh)

	p.mu.Unlock()

	for ; pc != nil; pc = pc.next {
		pc.c.Close()
	}
}

// Upgrade 将一条普通连接升级为连接池
func Upgrade(conn Conn) (Conn, error) {
	dc, ok := conn.(*defaultConn)
	if !ok {
		// 如果不能断言为defaultConn，说明已经升级过了
		return dc, errors.New("upgrade connection error: conn is already upgraded")
	}

	// TODO: 加载配置
	p := &pool{
		signal:     make(chan struct{}),
		data:       queue.New(),
		localAddr:  dc.nc.LocalAddr().String(),
		remoteAddr: dc.nc.RemoteAddr().String(),
	}

	p.poolSize++
	p.put(&poolConn{c: dc, t: time.Now()})
	go p.recv(dc)

	return p, nil
}

// Join 将一条普通连接加入到连接池
func Join(dst Conn, conn Conn) (c Conn, err error) {
	dc, ok := conn.(*defaultConn)
	if !ok {
		return dc, errors.New("join conn to pool error: conn is already a pool")
	}

	p, ok := dst.(*pool)
	if !ok {
		dst, _ = Upgrade(dst)
		p = dst.(*pool)
	}

	p.put(&poolConn{c: dc, t: time.Now()})
	p.mu.Lock()
	p.poolSize++
	p.mu.Unlock()
	go p.recv(dc)

	return p, nil
}

// dial 通过dial建立一条连接，并向accept方发送一个请求升级连接的报文
func (p *pool) dial() (conn Conn, err error) {
	remoteAddr := p.Addr()
	if remoteAddr == "" {
		return nil, errors.New("dial up connection error: remote address doesn't exist")
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

	// 如果设置了排队等待，当空闲连接不够时
	// 会进入排队状态
	// waitCh的buffer等于maxPoolSize
	if p.wait {
		<-p.waitCh
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	// 当连接池现有连接数大于corePoolSize时，检测多余的连接是否超过最大生存时间
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

	// 从连接池获取一条连接
	pc := p.connSet.back
	if pc != nil {
		p.connSet.popBack()
		c = pc
		return
	}

	// 如果连接池中没有空闲连接
	// 当连接数大于maxPoolSize，且不支持排队等待，抛出异常
	if !p.wait && p.poolSize >= p.maxPoolSize {
		return nil, errors.New("get conn error: connection pool exhausted")
	}

	// 如果不支持dial的方式获取连接，抛出异常
	if !p.canDial {
		return nil, errors.New("get conn error: connection pool exhausted and can't dial up")
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

	p.mu.Lock()
	p.connSet.pushFront(pc)
	if p.connSet.len > p.corePoolSize {
		pc = p.connSet.back
		pc.c.Close()
		p.poolSize--
		p.connSet.popBack()
	} else {
		pc = nil
	}
	p.mu.Unlock()

	if p.wait {
		p.waitCh <- struct{}{}
	}
}

func (p *pool) recv(c *defaultConn) {

	for {
		opcode, data, err := c.Recv()
		tmp := &info{opcode: opcode, data: data, err: err}
		if err != nil {
			c.Close()
			break
		}

		p.dataLock.Lock()
		isWaiting := p.data.IsEmpty()
		p.data.Push(tmp)
		p.dataLock.Unlock()

		if isWaiting {
			p.signal <- struct{}{}
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

// makeOpcode 组装用于升级连接的opcode
func makeOpcode() int {
	dispatchCode := frame.CtrlConnCode // 0 ~ 7

	version := frame.ProtocolVersion // 0 ~ 15

	return ((1<<3)|dispatchCode)<<4 | version
}
