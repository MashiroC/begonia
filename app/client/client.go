// Package client 的 api 层
package client

import (
	"context"
	"github.com/MashiroC/begonia/app"
	"github.com/MashiroC/begonia/app/coding"
	"github.com/MashiroC/begonia/config"
	cRegister "github.com/MashiroC/begonia/core/register"
	"github.com/MashiroC/begonia/internal/register"
	"github.com/MashiroC/begonia/logic"
	"github.com/MashiroC/begonia/tool/retry"
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
	mode   app.ServiceAppModeTyp
	lg     *logic.Client
	ctx    context.Context
	cancel context.CancelFunc

	register register.Register
}

// Service 获取一个服务
func (r *rClient) Service(serviceName string) (s Service, err error) {

	if r.mode == app.Ast {
		s = r.newAstService(serviceName)
	} else {
		s = r.newService(serviceName)
	}

	return
}

func (r *rClient) newService(name string) Service {

	s := &rService{
		name: name,
		c:    r,
	}

	f, err := r.parseServiceInfo(name)
	if err != nil {
		go retry.Always("getService", func() bool {
			f, err = r.parseServiceInfo(name)
			if err != nil {
				return false
			}

			s.funs = f
			return true
		}, config.C.App.GetServiceRetrySeconds)
	} else {
		s.funs = f
	}

	return s
}

func (r *rClient) parseServiceInfo(name string) (fMap map[string]Fun, err error) {
	var fs []cRegister.FunInfo

	fs, err = r.register.Get(name)
	if err != nil {
		log.Println("error in get services:", err)
		return
	}

	funs := make([]Fun, 0, len(fs))

	for _, f := range fs {
		var inCoder, outCoder coding.Coder

		inCoder, err = coding.NewAvro(f.InSchema)
		if err != nil {
			return
		}

		outCoder, err = coding.NewAvro(f.OutSchema)
		if err != nil {
			return
		}
		funs = append(funs, Fun{
			Name:     f.Name,
			InCoder:  inCoder,
			OutCoder: outCoder,
		})
	}

	fMap = make(map[string]Fun, len(funs))

	for i := 0; i < len(funs); i++ {
		fMap[funs[i].Name] = funs[i]
	}

	log.Printf("client get service [%s] success, func list: %s\n", name, reflect.ValueOf(fMap).MapKeys())

	return
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
