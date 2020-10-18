// Time : 2020/9/29 20:24
// Author : Kieran

// logic
package logic

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// waitchans.go something

type waitCallback = func(*CallResult)

type WaitChans struct {
	chLock sync.RWMutex
	ch     map[string]chan *CallResult
}

func NewWaitChans() *WaitChans {
	return &WaitChans{
		chLock: sync.RWMutex{},
		ch:     make(map[string]chan *CallResult),
	}
}

func (w *WaitChans) Callback(reqId string, cr *CallResult) (err error) {
	w.chLock.RLock()
	ch, exist := w.ch[reqId]
	w.chLock.RUnlock()
	if !exist {
		err = fmt.Errorf("reqId [%s] not found!", reqId)
		return
	}

	ch <- cr
	return
}

func (w *WaitChans) AddCallback(ctx context.Context, reqId string, callback waitCallback) {
	timeout, _ := context.WithTimeout(ctx, 9*time.Second) //TODO:抽成配置的

	ch := make(chan *CallResult)

	w.chLock.Lock()
	w.ch[reqId] = ch
	w.chLock.Unlock()

	go w.goWait(reqId, timeout.Done(), ctx.Done(), callback, ch)

}

func (w *WaitChans) goWait(reqId string, timeout, parent <-chan struct{}, cb waitCallback, ch chan *CallResult) {
	var f *CallResult
	var errStr string
	select {
	case <-timeout:
		// 超时
		errStr = "request timeout"
	case <-parent:
		// 外部取消了
		errStr = "request cancel"
	case f = <-ch:
		// 收到回调
	}

	w.chLock.Lock()
	delete(w.ch, reqId)
	w.chLock.Unlock()

	if errStr != "" {
		f = &CallResult{
			Result: nil,
			Err:    errStr,
		}
	}
	cb(f)
}
