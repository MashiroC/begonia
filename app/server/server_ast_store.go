package server

import (
	"fmt"
	"sync"
)

type astServiceStore struct {
	l sync.Mutex
	v map[string]astDo
}

func newAstServiceStore() *astServiceStore {
	return &astServiceStore{v: make(map[string]astDo)}
}

func (s *astServiceStore) get(service string) (do astDo, err error) {
	s.l.Lock()
	defer s.l.Unlock()
	var ok bool
	do, ok = s.v[service]
	if !ok {
		err = fmt.Errorf("service [%s] not exist", service)
	}
	return
}

func (s *astServiceStore) store(service string, fun astDo) (err error) {
	s.l.Lock()
	defer s.l.Unlock()
	if _, ok := s.v[service]; ok {
		return fmt.Errorf("service [%s] exist, you cannot store it", service)
	}
	s.v[service] = fun
	return
}
