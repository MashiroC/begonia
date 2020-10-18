package containers

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type RespSet struct {
	callbackMap map[string]respEntry
	l           sync.RWMutex

	timeoutChan chan RespTimeoutCall
}

func NewRespSet() *RespSet {
	return &RespSet{
		callbackMap: make(map[string]respEntry),
		l:           sync.RWMutex{},
		timeoutChan: make(chan RespTimeoutCall),
	}

}

type RespTimeoutCall struct {
	ConnUuid string
	Service  string
	Fun      string
}

type respEntry struct {
	cancel   context.CancelFunc
	reqId    string
	service  string
	fun      string
	connUuid string
}

func (s *RespSet) SignCallback(ctx context.Context, reqId, service, fun, connUuid string) (err error) {
	timeoutCtx, _ := context.WithTimeout(ctx, 1000000*time.Second)
	cancelCtx, cancel := context.WithCancel(ctx)
	s.l.Lock()
	s.callbackMap[reqId] = respEntry{
		cancel:   cancel,
		reqId:    reqId,
		service:  service,
		fun:      fun,
		connUuid: connUuid,
	}
	s.l.Unlock()
	// 开一个协程管超时
	go func() {
		select {
		case <-timeoutCtx.Done():
			// 超时了
			s.timeoutChan <- RespTimeoutCall{
				ConnUuid: connUuid,
				Service:  service,
				Fun:      fun,
			}
		case <-cancelCtx.Done():
			// 正常请求 这个被取消了 不用任何处理
		}
		// 从map里删掉这个
		s.l.Lock()
		delete(s.callbackMap, reqId)
		s.l.Unlock()
	}()

	return
}

func (s *RespSet) GetCallbackByReqId(reqId string) (connUuid, service, fun string, err error) {
	s.l.RLock()
	defer s.l.RUnlock()

	e, exist := s.callbackMap[reqId]
	if !exist {
		err = fmt.Errorf("reqId [%s] not found!", reqId)
		return
	}

	connUuid = e.connUuid
	service = e.service
	fun = e.fun
	return
}

func (s *RespSet) CancelCallback(reqId string) (err error) {
	s.l.RLock()
	e, exist := s.callbackMap[reqId]
	s.l.RUnlock()
	if !exist {
		err = fmt.Errorf("reqId [%s] not found!", reqId)
		return
	}
	e.cancel()
	return
}

func (s *RespSet) TimeoutChan() chan RespTimeoutCall {
	return s.timeoutChan
}
