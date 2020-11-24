// Time : 2020/8/3 20:04
// Author : MashiroC

// containers
package containers

import (
	"github.com/MashiroC/begonia/dispatch/frame"
	"context"
	"fmt"
	"sync"
	"time"
)

// waitchans.go something

type waitCallback = func(frame.Frame, error)

type WaitChans struct {
	chLock sync.RWMutex
	ch     map[string]chan frame.Frame
}

func NewWaitChans() *WaitChans {
	return &WaitChans{
		//waitLock: sync.RWMutex{},
		//wait:     make(map[string]waitCallback),
		chLock: sync.RWMutex{},
		ch:     make(map[string]chan frame.Frame),
	}
}

func (w *WaitChans) Callback(reqId string, f frame.Frame) (err error) {
	w.chLock.RLock()
	ch, exist := w.ch[reqId]
	w.chLock.RUnlock()
	if !exist {
		err = fmt.Errorf("reqId [%s] not found!", reqId)
		return
	}

	ch <- f
	return
}

func (w *WaitChans) AddCallback(ctx context.Context, reqId string, callback waitCallback) {
	timeout, _ := context.WithTimeout(ctx, 9*time.Second) //TODO:抽成配置的

	ch := make(chan frame.Frame)

	w.chLock.Lock()
	w.ch[reqId] = ch
	w.chLock.Unlock()

	go w.goWait(reqId, timeout.Done(), ctx.Done(), callback, ch)

}

func (w *WaitChans) goWait(reqId string, timeout, parent <-chan struct{}, cb waitCallback, ch chan frame.Frame) {
	var f frame.Frame
	var err error
	select {
	case <-timeout:
		// 超时
		err = fmt.Errorf("request timeout!")
	case <-parent:
		// 外部取消了
		err = fmt.Errorf("request cancel!")
	case f = <-ch:
		// 收到回调
	}

	w.chLock.Lock()
	delete(w.ch, reqId)
	w.chLock.Unlock()

	cb(f, err)
}
