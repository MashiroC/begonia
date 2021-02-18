package conn

import (
	"fmt"
	"testing"
	"time"
)

const (
	addr     = ":14949"
	poolSize = 5
	caseNum  = 100
)

func poolRecv(c Conn) {
	var num = 0
	for {
		_,_, err := c.Recv()
		if err != nil {
			panic(err)
		}
		num++
		//fmt.Println(opcode, data)
	}
}

func init() {
	ltCh, errCh := Listen(addr)
	var ltPool Conn
	var err error
	go func() {
		for {
			select {
			case c := <-ltCh:
				if ltPool == nil {
					ltPool, err = Upgrade(c)
					go poolRecv(ltPool)
					if err != nil {
						panic(err)
					}
					continue
				}
				err = Join(ltPool, c)
				if err != nil {
					panic(err)
				}
			case err = <-errCh:
				panic(err)
			}
		}
	}()

}

func TestUpgrade(t *testing.T) {
	c, err := Dial(addr)
	if err != nil {
		panic(err)
	}
	fmt.Println("dial", c.Addr())
	p, err := Upgrade(c)
	if err != nil {
		panic(err)
	}
	for i := 0; i < poolSize; i++ {
		cTmp, err := Dial(addr)
		if err != nil {
			panic(err)
		}
		err = Join(p, cTmp)
		if err != nil {
			panic(err)
		}
	}
	for i := 0; i < caseNum; i++ {
		go func() {
			err = p.Write(byte(1), []byte{1, 2, 3})
			if err != nil {
				panic(err)
			}
		}()

	}
	time.Sleep(3 * time.Second)
}
