package containers

import (
	"fmt"
	"github.com/linkedin/goavro/v2"
	"sync"
)

type ServiceSet struct {
	sLock    sync.RWMutex
	services map[string]Service

	uLock sync.Mutex
	uuids map[string][]string
}

func NewServiceSet() *ServiceSet {
	return &ServiceSet{
		sLock:    sync.RWMutex{},
		services: make(map[string]Service),
		uLock:    sync.Mutex{},
		uuids:    make(map[string][]string),
	}
}

type Service struct {
	Name string
	Funs []Fun
}

type Fun struct {
	Name      string
	InSchema  string
	inCodec   *goavro.Codec
	OutSchema string
	outCodec  *goavro.Codec
}

func (s *ServiceSet) Sign(uuid string, service Service) (err error) {
	s.sLock.Lock()

	if _, exist := s.services[service.Name]; exist {
		err = fmt.Errorf("service [%s] exist!", service.Name)
		return
	}

	for i := 0; i < len(service.Funs); i++ {
		service.Funs[i].inCodec, err = goavro.NewCodec(service.Funs[i].InSchema)
		if err != nil {
			return
		}
		service.Funs[i].outCodec, err = goavro.NewCodec(service.Funs[i].OutSchema)
		if err != nil {
			return
		}

	}

	s.services[service.Name] = service

	s.sLock.Unlock()

	s.uLock.Lock()
	slice, exist := s.uuids[uuid]
	if exist {
		s.uuids[uuid] = append(slice, service.Name)
	} else {
		s.uuids[uuid] = []string{service.Name}
	}
	s.uLock.Unlock()
	return
}

func (s *ServiceSet) UnSign(uuid string) (err error) {

	s.uLock.Lock()
	slice, exist := s.uuids[uuid]
	if !exist {
		err = fmt.Errorf("uuid [%s] services not exist!", uuid)
	}
	delete(s.uuids, uuid)
	s.uLock.Unlock()

	s.sLock.Lock()
	for _, service := range slice {
		if _, exist := s.services[service]; !exist {
			err = fmt.Errorf("service [%s] not exist!", service)
			break
		}

		delete(s.services, service)
	}
	s.sLock.Unlock()

	return
}

func (s *ServiceSet) GetService(name string) (ser Service, exist bool) {
	s.sLock.RLock()
	defer s.sLock.RUnlock()

	ser, exist = s.services[name]
	return
}

func (s *ServiceSet) GetFunInCodec(service, fun string) (codec *goavro.Codec, err error) {
	codec, _, err = s.GetFunCodec(service, fun)
	return
}

func (s *ServiceSet) GetFunOutCodec(service, fun string) (codec *goavro.Codec, err error) {
	_, codec, err = s.GetFunCodec(service, fun)
	return
}

func (s *ServiceSet) GetFunCodec(service, fun string) (in *goavro.Codec, out *goavro.Codec, err error) {
	s.sLock.RLock()
	defer s.sLock.RUnlock()

	res, exist := s.services[service]
	if !exist {
		err = fmt.Errorf("service [%s] not exist!", service)
		return
	}

	flag := false
	for _, f := range res.Funs {
		if fun == f.Name {
			in = f.inCodec
			out = f.outCodec
			flag = true
			break
		}
	}
	if !flag {
		err = fmt.Errorf("fun not found!")
	}
	return
}
