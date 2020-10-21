package logic

import (
	"sync"
	"time"
)

type ReqSet struct {
	m map[string]*reqSetEntry
	l sync.Mutex

	overtime time.Duration
}

type reqSetEntry struct {
	connID string
	t      time.Time
}

func NewReqSet(overtime int) *ReqSet {
	return &ReqSet{
		m:        make(map[string]*reqSetEntry),
		l:        sync.Mutex{},
		overtime: time.Duration(overtime),
	}
}

func (s *ReqSet) Add(reqID, connID string) {
	s.l.Lock()
	defer s.l.Unlock()

	s.m[reqID] = &reqSetEntry{
		connID: connID,
		t:      time.Now().Add(s.overtime * time.Second),
	}

}
func (s *ReqSet) Get(reqID string) (connID string, ok bool) {

	s.l.Lock()
	defer s.l.Unlock()

	v, ok := s.m[reqID]
	if !ok {
		return
	}

	if time.Now().After(v.t) {
		ok = false
		return
	}

	connID = v.connID
	ok = true

	return
}

func (s *ReqSet) Remove(reqID string) {

}

func (s *ReqSet) CleanUp() {
	s.l.Lock()
	defer s.l.Unlock()

	now := time.Now()
	for k, v := range s.m {
		if now.After(v.t) {
			delete(s.m, k)
		}
	}
}
