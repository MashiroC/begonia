// Package client 的 api 层
package client

import (
	"context"
	"github.com/MashiroC/begonia/core"
	"github.com/MashiroC/begonia/internal"
	"github.com/MashiroC/begonia/internal/coding"
	"github.com/MashiroC/begonia/logic"
	"log"
	"reflect"
)

// logic_service.go something

type Client interface {
	Service(name string) (Service, error)
	FunSync(serviceName, funName string) (RemoteFunSync, error)
	FunAsync(serviceName, funName string) (RemoteFunAsync, error)
	Wait()
	Close()
}

// FunInfo 远程函数的一个封装
type Fun struct {
	Name     string       // 远程函数名
	InCoder  coding.Coder // 远程函数入参的编码器
	OutCoder coding.Coder // 远程函数出参的编码器
}

func NewClient(lg *logic.Client) *rClient {
	return &rClient{
		lg:     lg,
		ctx:    nil,
		cancel: nil,
	}
}

// rClient 客户端的github.com/MashiroC/begonia实现
type rClient struct {
	lg *logic.Client
	ctx    context.Context
	cancel context.CancelFunc
}

// Service 获取一个服务
func (r *rClient) Service(serviceName string) (s Service, err error) {

	// TODO:这里要换成注册器

	if internal.ServiceAppMode == internal.ServiceAppModeAst {
		s = newAstService(serviceName, r)
	} else if internal.ServiceAppMode == internal.ServiceAppModeReflect {
		res := r.lg.CallSync(core.Call.ServiceInfo(serviceName))

		if res.Err != nil {
			err = res.Err
			return
		}

		fs := core.Result.ServiceInfo(res.Result)

		s = r.newService(serviceName, fs)

	} else {
		panic("eeeeeeeeeerror!")
	}
	return
}

func (r *rClient) newService(name string, funs []Fun) Service {

	f := make(map[string]Fun, len(funs))

	for i := 0; i < len(funs); i++ {
		f[funs[i].Name] = funs[i]
	}

	log.Printf("client get service [%s] success, func list: %s", name, reflect.ValueOf(f).MapKeys())

	return &rService{
		name: name,
		funs: f,
		c:    r,
	}
}

func (r *rClient) FunSync(serviceName, funName string) (rf RemoteFunSync, err error) {
	s, err := r.Service(serviceName)
	if err != nil {
		return
	}
	rf, err = s.FuncSync(funName)
	return
}

func (r *rClient) FunAsync(serviceName, funName string) (rfa RemoteFunAsync, err error) {
	s, err := r.Service(serviceName)
	if err != nil {
		return
	}
	rfa, err = s.FuncAsync(funName)
	return
}

func (r *rClient) Wait() {
	select {
	case <-r.ctx.Done():
	}
}

func (r *rClient) Close() {
	r.lg.Close()
	r.cancel()
}
