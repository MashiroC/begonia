package core

import (
	"fmt"
	"github.com/MashiroC/begonia/app/coding"
	"sync"
)

type registerServiceStore struct {
	l sync.RWMutex
	m map[string]*registerService

	connIndexes map[string][]string
	connLock    sync.Mutex
}

type registerService struct {
	connID string
	name   string
	funs   []coding.FunInfo
}

func newStore() *registerServiceStore {
	return &registerServiceStore{
		m:           make(map[string]*registerService),
		connIndexes: make(map[string][]string),
	}
}

func (s *registerServiceStore) Get(service string) (rs *registerService, ok bool) {
	s.l.RLock()
	defer s.l.RUnlock()

	rs, ok = s.m[service]
	return
}

func (s *registerServiceStore) Add(connID, serviceName string, funs []coding.FunInfo) (err error) {
	s.l.Lock()
	s.connLock.Lock()
	defer s.l.Unlock()
	defer s.connLock.Unlock()

	if _, ok := s.m[serviceName]; ok {
		err = fmt.Errorf("service [%s] existed", serviceName)
		return
	}

	s.m[serviceName] = &registerService{
		connID: connID,
		name:   serviceName,
		funs:   funs,
	}

	v, ok := s.connIndexes[connID]
	if ok {
		s.connIndexes[connID] = append(v, serviceName)
	} else {
		s.connIndexes[connID] = []string{serviceName}
	}
	return
}

func (s *registerServiceStore) Unlink(connID string) {
	s.l.Lock()
	s.connLock.Lock()
	defer s.l.Unlock()
	defer s.connLock.Unlock()

	services, ok := s.connIndexes[connID]
	if !ok {
		return
	}

	for _, service := range services {
		delete(s.m, service)
	}

}
