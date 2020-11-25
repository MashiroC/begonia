package logic

import (
	"context"
	"github.com/MashiroC/begonia/config"
	"github.com/MashiroC/begonia/tool/berr"
	"sync"
	"time"
)

// waitchans.go 等待管道

// waitCallback 等待的回调函数
type waitCallback = func(*CallResult)

// WaitChans 等待管道，拥有注册回调、回调的方法
type WaitChans struct {
	chLock sync.RWMutex                // 锁
	ch     map[string]chan *CallResult // 存储的map
}

// NewWaitChans 创建一个实例
func NewWaitChans() *WaitChans {
	return &WaitChans{
		chLock: sync.RWMutex{},
		ch:     make(map[string]chan *CallResult),
	}
}

// Callback 注册后调后，根据reqID来调用回调
func (w *WaitChans) Callback(reqID string, cr *CallResult) (err error) {
	w.chLock.RLock()
	ch, exist := w.ch[reqID]
	w.chLock.RUnlock()

	if !exist {
		err = berr.NewF("waitChan", "callback", "reqID [%s] not found", reqID)
		return
	}

	ch <- cr
	return
}

// AddCallback 添加一个回调
func (w *WaitChans) AddCallback(ctx context.Context, reqID string, callback waitCallback) {
	timeout, _ := context.WithTimeout(ctx, time.Duration(config.C.Logic.RequestTimeOut)*time.Second)

	ch := make(chan *CallResult)

	w.chLock.Lock()
	w.ch[reqID] = ch
	w.chLock.Unlock()

	go w.goWait(reqID, timeout.Done(), ctx.Done(), callback, ch)

}

// goWait 这个需要开一个新协程 来等待结果或者超时
func (w *WaitChans) goWait(reqID string, timeout, parent <-chan struct{}, cb waitCallback, ch chan *CallResult) {

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
	delete(w.ch, reqID)
	w.chLock.Unlock()

	if errStr != "" {
		f = &CallResult{
			Result: nil,
			Err:    berr.New("waitChat", "callback", errStr),
		}
	}

	cb(f)

}
