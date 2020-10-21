package core

import "fmt"

func (s *SubService) register(connID string, param ServiceInfo) (err error) {
	err = s.services.Add(connID, param.Service, param.Funs)
	fmt.Println(s.services.m)
	return err
}
