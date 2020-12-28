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
	cli appClient.Client
	server appServer.Server
}

func (m *MixNode) Service(name string) (s appClient.Service, err error) {
	panic("implement me")
}

func (m *MixNode) FunSync(serviceName, funName string) (rf appClient.RemoteFunSync, err error) {
	panic("implement me")
}

func (m *MixNode) FunAsync(serviceName, funName string) (rf appClient.RemoteFunAsync, err error) {
	panic("implement me")
}

func (m *MixNode) Wait() {
	panic("implement me")
}

func (m *MixNode) Close() {
	panic("implement me")
}

func (m *MixNode) Register(name string, service interface{}) {
	panic("implement me")
}
