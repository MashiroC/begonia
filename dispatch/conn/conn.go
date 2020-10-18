// Time : 2020/10/6 0:40
// Author : Kieran

// conn
package conn

import (
	"bufio"
	"net"
)

const (
	defaultNetwork = "tcp4"

	// 第一个length的最大值
	baseLenMax = 255
	// 上述变量的byte
	baseLenMaxByte = byte(baseLenMax)

	// length的最大值
	extendLengthMax = 255 * 255
)

// conn.go something

type Conn interface {
	Write(opcode byte, data []byte) error
	Recv() (opcode byte, data []byte, err error)
	Close()
}

func Dial(addr string) (c Conn, err error) {

	nc, err := net.Dial(defaultNetwork, addr)
	if err != nil {
		return
	}

	c = warp(nc)
	return
}

func Listen(addr string) (acceptCh chan Conn, errCh chan error) {
	errCh = make(chan error)
	acceptCh = make(chan Conn)

	lt, err := net.Listen(defaultNetwork, addr)
	if err != nil {
		errCh <- err
		close(errCh)
		close(acceptCh)
		return
	}

	go func(lt net.Listener) {
		for {
			nc, err := lt.Accept()
			if err != nil {
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

func warp(nc net.Conn) (c Conn) {
	r := bufio.NewReader(nc)
	w := bufio.NewWriter(nc)
	rw := bufio.NewReadWriter(r, w)
	c = &defaultConn{
		nc:  nc,
		buf: rw,
	}
	return
}
