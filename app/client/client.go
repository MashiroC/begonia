// Time : 2020/9/19 16:02
// Author : Kieran

// client
package client

import (
	"begonia2/app"
	"begonia2/logic"
	"context"
	"errors"
)

// service.go something

// Client 客户端的接口
type Client interface {
	Service(name string) (Service, error)
	FunSync(serviceName, funName string) (RemoteFunSync, error)
	FunAsync(serviceName, funName string) (RemoteFunAsync, error)
	Wait()
	Close()
}

// rClient 客户端的begonia实现
type rClient struct {
	lg logic.Client
	//pool *conn.Pool
	ctx    context.Context
	cancel context.CancelFunc
}

// Service 获取一个服务
func (r *rClient) Service(serviceName string) (s Service, err error) {

	res := r.lg.CallSync(app.Core.SignInfo(serviceName))

	if res.Err != "" {
		err = errors.New(res.Err)
		return
	}

	fs := app.Core.SignInfoResult(res.Result)

	s = r.newService(serviceName, fs)
	return
}

func (r *rClient) newService(name string, funs []app.FunInfo) Service {

	f := make(map[string]app.FunInfo, len(funs))

	for i := 0; i < len(funs); i++ {
		f[funs[i].Name] = funs[i]
	}
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

func (r *rClient) Close(){
	r.lg.Close()
	r.cancel()
}
