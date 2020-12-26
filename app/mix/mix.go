package mix

import (
	appClient "github.com/MashiroC/begonia/app/client"
	appServer "github.com/MashiroC/begonia/app/server"
)

type Mix interface {
	appClient.Client
	appServer.Server
}

type MixNode struct {
	s appServer.Server
}

func (m *MixNode) Service(name string) (appClient.Service, error) {
	return m.c.Service(name)
}

func (m *MixNode) FunSync(serviceName, funName string) (appClient.RemoteFunSync, error) {
	return m.c.FunSync(serviceName, funName)
}

func (m *MixNode) FunAsync(serviceName, funName string) (appClient.RemoteFunAsync, error) {
	return m.c.FunAsync(serviceName, funName)
}

func (m *MixNode) Close() {
	m.c.Close()
}

func (m *MixNode) Register(name string, service interface{}) {
	m.s.Register(name, service)
}

func (m *MixNode) Wait() {
	m.s.Wait()
	m.c.Wait()
}
