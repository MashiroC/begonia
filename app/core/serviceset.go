package core

import (
	"begonia2/opcode/coding"
	"fmt"
	"sync"
)

type serviceSet struct {
	l sync.RWMutex
	m map[string]*registerService
}

type registerService struct {
	connID string
	name   string
	funs   []coding.FunInfo
}

func newServiceSet() *serviceSet {
	return &serviceSet{
		m: make(map[string]*registerService),
	}
}

func (s *serviceSet) Get(service string) (rs *registerService, ok bool) {
	s.l.RLock()
	defer s.l.RUnlock()

	rs, ok = s.m[service]
	return
}

func (s *serviceSet) Add(connID, serviceName string, funs []coding.FunInfo) (err error) {
	s.l.Lock()
	defer s.l.Unlock()

	if _, ok := s.m[serviceName]; ok {
		err = fmt.Errorf("service [%s] existed", serviceName)
		return
	}

	s.m[serviceName] = &registerService{
		connID: connID,
		name:   serviceName,
		funs:   funs,
	}

	return
}
