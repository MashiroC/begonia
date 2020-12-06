package service

// service_ast.go ast实现的api

import (
	"context"
	"github.com/MashiroC/begonia/app/coding"
	"github.com/MashiroC/begonia/core"
	"github.com/MashiroC/begonia/logic"
	"github.com/MashiroC/begonia/tool/berr"
)

type astDo = func(fun string, param []byte) (result []byte, err error)

type CodeGenFunc struct {
}

type CodeGenService interface {
	Do(fun string, param []byte) (result []byte, err error)
	AvroSchema() string
	FuncList() []coding.FunInfo
}

// astService ast树代码生成的ast service api
type astService struct {
	lg     *logic.Service
	ctx    context.Context
	cancel context.CancelFunc

	store           *astServiceStore
	isLocalRegister bool
}

func (r *astService) Register(name string, service interface{}) {

	cgs, ok := service.(CodeGenService)
	if !ok {
		panic("please use code-gen")
	}

	// local store
	if err := r.store.store(name, cgs.Do); err != nil {
		panic(err)
	}

	fs := cgs.FuncList()

	// register
	if r.isLocalRegister {
		register := core.Call.Register(name, fs)
		_, err := core.C.Invoke("", "", register.Fun, register.Param)
		if err != nil {
			panic(err)
		}
	} else {
		res := r.lg.CallSync(core.Call.Register(name, fs))
		// TODO:handler error
		if res.Err != nil {
			panic(res.Err)
		}

		var ok bool
		_ = success.DecodeIn(res.Result, &ok)

		if ok {
			return
		}
	}

}

func (r *astService) Wait() {
	<-r.ctx.Done()
}

func (r *astService) handleMsg(msg *logic.Call, wf logic.ResultFunc) {

	if r.isLocalRegister && msg.Service == core.ServiceName && msg.Fun == "ServiceInfo" {
		res, err := core.C.Invoke("", "", "ServiceInfo", msg.Param)
		wf.Result(&logic.CallResult{
			Result: res,
			Err:    err,
		})
	}

	do, err := r.store.get(msg.Service)
	if err != nil {
		wf.Result(&logic.CallResult{
			Err: berr.Warp("app.service", "handle get func", err),
		})
		return
	}

	data, err := do(msg.Fun, msg.Param)
	if err != nil {
		wf.Result(&logic.CallResult{
			Err: berr.Warp("app.service", "handle", err),
		})
		return
	}

	wf.Result(&logic.CallResult{Result: data})
}
