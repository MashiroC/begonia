package register

import (
	"context"
	"fmt"
)

type Service struct {
	Name string
	Mode string
	Funs []FunInfo
}

type FunInfo struct {
	Name      string
	InSchema  string
	OutSchema string
}

func NewCoreRegister() *CoreRegister {
	return &CoreRegister{
		services: newStore(),
	}
}

type CoreRegister struct {
	services *registerServiceStore
}

func (r *CoreRegister) Register(ctx context.Context, si Service) (err error) {
	v := ctx.Value("info")
	info := v.(map[string]string)
	connID := info["connID"]

	err = r.services.Add(connID, si.Name, si.Funs)
	return err
}

func (r *CoreRegister) ServiceInfo(serviceName string) (si Service, err error) {
	service, ok := r.services.Get(serviceName)
	if !ok {
		err = fmt.Errorf("server [%s] not found", serviceName)
		return
	}

	si.Funs = service.funs
	si.Name = serviceName
	si.Mode = "avro"
	return
}
