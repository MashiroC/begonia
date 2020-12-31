package logic

import (
	"context"
	"fmt"
	"github.com/MashiroC/begonia/config"
	"sync"
	"time"
)

// callbacks.go 等待管道

// CallbackStore 回调仓库，拥有注册回调、回调的方法
type CallbackStore struct {
	chLock sync.RWMutex                // 锁
	ch     map[string]chan *CallResult // 存储的map
}

// NewWaitChans 创建一个实例
func NewWaitChans() *CallbackStore {
	return &CallbackStore{
		chLock: sync.RWMutex{},
		ch:     make(map[string]chan *CallResult),
	}
}

// Callback 注册后调后，根据reqID来调用回调
func (w *CallbackStore) Callback(reqID string, cr *CallResult) (err error) {
	w.chLock.RLock()
	ch, exist := w.ch[reqID]
	w.chLock.RUnlock()

	if !exist {
		err = fmt.Errorf("callbacks callback error: reqID [%s] not exist", reqID)
		return
	}

	ch <- cr
	return
}

// AddCallback 添加一个回调
func (w *CallbackStore) AddCallback(ctx context.Context, reqID string, callback Callback) {
	timeout, _ := context.WithTimeout(ctx, time.Duration(config.C.Logic.RequestTimeOut)*time.Second)

	ch := make(chan *CallResult)

	w.chLock.Lock()
	w.ch[reqID] = ch
	w.chLock.Unlock()

	go w.goWait(reqID, timeout.Done(), ctx.Done(), callback, ch)

}

// goWait 这个需要开一个新协程 来等待结果或者超时
func (w *CallbackStore) goWait(reqID string, timeout, parent <-chan struct{}, cb Callback, ch chan *CallResult) {

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
			Err:    fmt.Errorf("callbacks wait error: %s", errStr),
		}
	}

	cb(f)

}
