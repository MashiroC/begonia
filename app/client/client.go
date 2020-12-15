// Package client 的 api 层
package client

import (
	"context"
	"fmt"
	"github.com/MashiroC/begonia/app"
	"github.com/MashiroC/begonia/core"
	"github.com/MashiroC/begonia/logic"
)

// logic_service.go something

// rClient 客户端的github.com/MashiroC/begonia实现
type rClient struct {
	lg *logic.Client
	//pool *conn.Pool
	ctx    context.Context
	cancel context.CancelFunc
}

// Service 获取一个服务
func (r *rClient) Service(serviceName string) (s Service, err error) {

	res := r.lg.CallSync(core.Call.ServiceInfo(serviceName))

	if res.Err != nil {
		err = res.Err
		return
	}

	fs := core.Result.ServiceInfo(res.Result)
	fmt.Println(fs)
	fmt.Println(res)

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

func (r *rClient) Close() {
	r.lg.Close()
	r.cancel()
}
