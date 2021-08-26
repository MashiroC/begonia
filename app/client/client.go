// Package client 的 api 层
package client

import (
	"context"
	"github.com/MashiroC/begonia/app"
	"github.com/MashiroC/begonia/app/coding"
	"github.com/MashiroC/begonia/internal/register"
	"github.com/MashiroC/begonia/logic"
	"log"
	"reflect"
)

type Client interface {
	Service(name string) (s Service, err error)
	FunSync(serviceName, funName string) (rf RemoteFunSync, err error)
	FunAsync(serviceName, funName string) (rf RemoteFunAsync, err error)
	Wait()
	Close()
}

// FunInfo 远程函数的一个封装
type Fun struct {
	Name     string       // 远程函数名
	InCoder  coding.Coder // 远程函数入参的编码器
	OutCoder coding.Coder // 远程函数出参的编码器
}

// rClient 客户端的github.com/MashiroC/begonia实现
type rClient struct {
	mode app.ServiceAppModeTyp
	lg     *logic.Client
	ctx    context.Context
	cancel context.CancelFunc

	register register.Register
}

// Service 获取一个服务
func (r *rClient) Service(serviceName string) (s Service, err error) {

	fs, err := r.register.Get(serviceName)
	if err != nil {
		return
	}

	funs := make([]Fun, 0, len(fs))

	for _, f := range fs {
		inCoder, err := coding.NewAvro(f.InSchema)
		if err != nil {
			return nil, err
		}

		outCoder, err := coding.NewAvro(f.OutSchema)
		if err != nil {
			return nil, err
		}
		funs = append(funs, Fun{
			Name:     f.Name,
			InCoder:  inCoder,
			OutCoder: outCoder,
		})
	}

	if r.mode == app.Ast {
		s = r.newAstService(serviceName, r)
	} else {
		s = r.newService(serviceName, funs)
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
