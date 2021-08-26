// Package conn 底层的连接的抽象
package conn

import (
	"bufio"
	"fmt"
	"net"
)

const (

	// 默认的连接协议 tcp ipv4
	defaultNetwork = "tcp4"

	// 第一个length的最大值
	baseLenMax = 255

	// 上述变量的byte
	baseLenMaxByte = byte(baseLenMax)

	// length的最大值
	extendLengthMax = 255 * 255
)

// Conn 连接的抽象接口
type Conn interface {
	Write(opcode byte, data []byte) error        // 写数据，线程安全
	Recv() (opcode byte, data []byte, err error) // 读数据
	Addr() string
	Close() // 关闭连接
}

// Dial 对一个地址建立一条tcp连接
func Dial(addr string) (c Conn, err error) {

	nc, err := net.Dial(defaultNetwork, addr)
	if err != nil {
		err = fmt.Errorf("conn dial error: %w", err)
		return
	}

	c = warp(nc)
	return
}

// Listen 监听一个地址
// 返回两个管道 一个是建立成功的连接的管道，一个是错误的管道
// 如果监听出现错误，errCh写一个错误进去，然后关闭两个管道
func Listen(addr string) (acceptCh chan Conn, errCh chan error) {

	errCh = make(chan error, 10)
	acceptCh = make(chan Conn, 10)

	lt, err := net.Listen(defaultNetwork, addr)
	if err != nil {
		err = fmt.Errorf("conn listen error: %w", err)
		errCh <- err
		close(errCh)
		close(acceptCh)
		return
	}

	go func(lt net.Listener) {
		for {
			nc, err := lt.Accept()
			if err != nil {
				err = fmt.Errorf("conn accept error: %w", err)
				errCh <- err
				close(errCh)
				close(acceptCh)
				return
			}
			c := warp(nc)
			acceptCh <- c
		}

	}(lt)

	return
}

// warp 包装一个net.Conn未为github.com/MashiroC/begonia.conn.Conn
func warp(nc net.Conn) (c Conn) {

	r := bufio.NewReader(nc)

	c = &defaultConn{
		nc:  nc,
		buf: r,
	}

	return
}
