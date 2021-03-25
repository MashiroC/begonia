package logger

import (
	"github.com/MashiroC/begonia/tool/log"
	"sync"
)

type loggerServiceStore struct {
	l sync.RWMutex

	// serviceName => buf
	m map[string]*storeInstance
	r *remoteLoggerService
}

func newLoggerServiceStore(r *remoteLoggerService) *loggerServiceStore {
	return &loggerServiceStore{
		l: sync.RWMutex{},
		m: make(map[string]*storeInstance),
		r: r,
	}
}

func (s *loggerServiceStore) Get(service string) (msg []Msg, ok bool) {
	s.l.RLock()
	defer s.l.RUnlock()

	ins, ok := s.m[service]
	if !ok {
		return
	}
	msg = ins.msg
	// 获得信息后，马上重置
	// do -> listen -> get
	ins.resetAll()
	return
}

func (s *loggerServiceStore) Store(serverName string, l *log.Log) (err error) {
	s.l.Lock()
	defer s.l.Unlock()
	if ins, ok := s.m[serverName]; ok {
		ins.msg = append(ins.msg, logToMsg(serverName, l))
		return
	}

	// 如果没有，证明没有开始
	instance := newStoreInstance(logToMsg(serverName, l))
	s.m[serverName] = instance
	// 开始定时
	go instance.order()
	// 开始监听
	go s.listenToSend(serverName, instance)
	return
}

//TODO:这里肯定有消耗，感觉可以直接传log

// log结构体转Msg
func logToMsg(serverName string, l *log.Log) Msg {
	return Msg{
		ServerName: serverName,
		Level:      int(l.GetLevel()),
		Fields:     l.Data,
		Time:       l.TimeNow.Unix(),
		Callers:    l.GetCaller(),
	}
}

//TODO: 系统错误收集
func (s *loggerServiceStore) listenToSend(serverName string, i *storeInstance) {
	for {
		select {
		case f := <-i.do:
			if f {
				if err := s.r.sendMsg(serverName); err != nil {
					return
				}
			} else {
				// 出错,输出剩下的
				s.r.sendMsg(serverName)
				return
			}
		}
	}
}
