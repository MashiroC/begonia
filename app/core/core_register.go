package core

import "fmt"

type RegisterParam struct {
}

type RegisterResult struct {
}

func (s *SubService) Register(param ServiceInfo) (err error) {
	fmt.Println("call Register!")

	return nil
}
