package logger

import (
	"sync"
	"time"
)

const WaitTime = time.Second * 10

// 储存的实体
type storeInstance struct {
	msg     []Msg
	timer   time.Ticker // 定时，每100s同步一次
	msgNum  int         // 当前消息条数
	numChan chan int
	do      chan bool

	lock sync.Mutex
}

// 初始化储存的实例
func newStoreInstance(m ...Msg) *storeInstance {
	return &storeInstance{
		msg:     m,
		timer:   *time.NewTicker(WaitTime),
		msgNum:  0,
		numChan: make(chan int),
		do:      make(chan bool),
		lock:    sync.Mutex{},
	}
}

func (i *storeInstance) order() {
	defer func() {
		i.do <- false
	}()
	for {
		select {
		case num := <-i.numChan:
			// 可以进行发送
			if num == 100 {
				i.do <- true
			}
		case <-i.timer.C:
			// 每100s，可进行发送
			i.do <- true

		}
	}
}

// 增加消息条数
func (i *storeInstance) addMsgNum() {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.msgNum += 1
	i.numChan <- i.msgNum
}

// 重置所有
func (i *storeInstance) resetAll() {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.msgNum = 0
	i.timer.Reset(WaitTime)
	i.msg = []Msg{}
}
