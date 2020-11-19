package core

import "fmt"

func (s *SubService) register(connID string, param ServiceInfo) (err error) {
	err = s.services.Add(connID, param.Service, param.Funs)
	return err
}

func (s *SubService) serviceInfo(serviceName string) (si ServiceInfo, err error) {
	service, ok := s.services.Get(serviceName)
	if !ok {
		err = fmt.Errorf("service [%s] not found", serviceName)
		return
	}

	si.Funs = service.funs
	si.Service = serviceName
	return
}
