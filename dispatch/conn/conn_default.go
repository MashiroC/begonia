package conn

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/MashiroC/begonia/config"
	"net"
	"sync"
	"time"
)

var (
	chPool       sync.Pool
	bytesPool    sync.Pool
	twoBytesPool sync.Pool
)

func init() {
	chPool = sync.Pool{New: func() interface{} {
		return make(chan int)
	}}

	bytesPool = sync.Pool{New: func() interface{} {
		return make([]byte, 0, 1024)
	}}

	twoBytesPool = sync.Pool{New: func() interface{} {
		return make([]byte, 2)
	}}
}

// defaultConn 默认的conn实现，单条tcp连接
type defaultConn struct {
	nc  net.Conn
	buf *bufio.Reader
	l   sync.Mutex
}

func (d *defaultConn) Addr() string {
	return d.nc.RemoteAddr().String()
}

func (d *defaultConn) Write(opcode byte, data []byte) (err error) {

	// 计算 payload length 与 extend payload length
	var size []byte
	if len(data) >= extendLengthMax {
		err = fmt.Errorf("conn write error: payload len [%d] oversize", len(data))
		return
	} else if len(data) >= baseLenMax {
		tmp := twoBytesPool.Get().([]byte)
		binary.BigEndian.PutUint16(tmp, uint16(len(data)))
		size = []byte{255}
		size = append(size, tmp...)
		twoBytesPool.Put(tmp)
	} else {
		size = append(size, byte(len(data)))
	}

	// 组装opcode length data
	tmp := bytesPool.Get().([]byte)
	defer bytesPool.Put(tmp[:0])
	tmp = append(tmp, opcode)
	tmp = append(tmp, size...)
	tmp = append(tmp, data...)

	d.l.Lock()
	defer d.l.Unlock()

	if _, err = d.nc.Write(tmp); err != nil {
		err = fmt.Errorf("conn write error: %w", err)
		return
	}

	return
}

func (d *defaultConn) Recv() (opcode byte, data []byte, err error) {
	// 检测一个有没有panic出来错误 有的话把连接关了
	defer func() {
		if errIn := recover(); err != nil {
			d.Close()
			err = errIn.(error)
			err = fmt.Errorf("conn recv error: %w", err)
			return
		}
	}()

	// 读opcode 不需要等超时 直接等第一个包来就行
	// 除了包头的opcode 剩下的都要等超时

	/*
		    4      4         8       0 || 16
		{opcode}{version}{length}{extendLength}
	*/

	// 拿opcode
	opcode, err = d.readByte()
	handlerErr(err)

	// 拿payload length
	baseLen, err := d.buf.ReadByte()
	handlerErr(err)
	payloadLen := uint(baseLen)

	// baseLen如果是255 读extend length
	if baseLen == baseLenMaxByte {
		// 这里读了两个byte 然后转化成int
		var extendLen []byte
		extendLen, err = d.read(2)
		handlerErr(err)
		payloadLen = uint(binary.BigEndian.Uint16(extendLen))
		// 我们不支持超过一定大小的包
		if payloadLen >= extendLengthMax {
			err = errors.New("payload len oversize")
			return
		}
	}

	if payloadLen == 0 {
		panic("payload length zero")
	}

	// 拿数据，线程安全，内存安全
	data, err = d.read(payloadLen)
	//n, err := buf.Read(data)
	handlerErr(err)

	return
}

// read 读一定长度的数据 超时时间在配置里
// 设计成这样是因为之前改一个高并发的bug的时候 发现了一个问题
// 在高并发场景下 client发了一个包 这边一次性read不完
// 报序列化错误 导致后面所有的包都乱序
func (d *defaultConn) read(len uint) (data []byte, err error) {

	data = make([]byte, len)
	n, err := d.readWithTimeout(data)
	if err != nil {
		return
	}

	// 一次没读够指定的len 继续读
	for n < int(len) {
		overSize := make([]byte, int(len)-n)
		size, err := d.readWithTimeout(overSize)
		if err != nil {
			return nil, err
		}

		for i := 0; i < size; i++ {
			data[n+i] = overSize[i]
		}
		n += size
	}

	return
}

// readWithTimeout 带超时时间的读，超时时间在config包里
func (d *defaultConn) readWithTimeout(b []byte) (n int, err error) {

	ch := chPool.Get().(chan int)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.C.Conn.ReadTimeout)*time.Second)
	defer cancel()
	go func() {
		n, err = d.buf.Read(b)
		ch <- n
		chPool.Put(ch)
	}()
	select {
	case n = <-ch:
	case <-ctx.Done():
		err = errors.New("read time out")
	}

	return
}

// readByte 读一个byte
func (d *defaultConn) readByte() (data byte, err error) {
	data, err = d.buf.ReadByte()
	return
}

func (d *defaultConn) Close() {
	d.nc.Close()
}

func handlerErr(err error) {
	if err != nil {
		panic(err)
	}
}
