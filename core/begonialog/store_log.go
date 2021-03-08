package centerlog

import (
	"fmt"
	"sync"
)

type logServiceStore struct {
	lock sync.Mutex

	// serviceName => L
	m        map[string]*logService
}

func newStore() *logServiceStore {
	return &logServiceStore{
		m: make(map[string]*logService),
	}
}

// 获取
func (s *logServiceStore) Get(service string) (rs *logService, ok bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	rs, ok = s.m[service]
	return
}

// 增加
func (s *logServiceStore) Add(serviceName string,l *logService) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.m[serviceName]; ok {
		err = fmt.Errorf("server [%s] existed", serviceName)
		return
	}
	s.m[serviceName]=l
	return
}

func (s *logServiceStore) Del(name string) {
	s.lock.Lock()
	defer s.lock.Unlock()


	delete(s.m, name)

}
