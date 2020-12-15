package service

import (
	"errors"
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
		err = errors.New("service not exist")
	}
	return
}

func (s *astServiceStore) store(service string, fun astDo) (err error) {
	s.l.Lock()
	defer s.l.Unlock()
	if _, ok := s.v[service]; ok {
		return errors.New("service exist")
	}
	s.v[service] = fun
	return
}
