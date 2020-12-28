package server

// server_ast.go ast实现的api

import (
	"context"
	coreRegister "github.com/MashiroC/begonia/core/register"
	"github.com/MashiroC/begonia/internal/register"
	"github.com/MashiroC/begonia/logic"
	"github.com/MashiroC/begonia/tool/berr"
)

type astDo = func(ctx context.Context, fun string, param []byte) (result []byte, err error)

type CodeGenFunc struct {
}

type CodeGenService interface {
	Do(ctx context.Context, fun string, param []byte) (result []byte, err error)
	FuncList() []coreRegister.FunInfo
}

// astService ast树代码生成的ast Server api
type astService struct {
	lg     *logic.Service
	ctx    context.Context
	cancel context.CancelFunc

	store    *astServiceStore
	register register.Register
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

	r.register.Register(name, fs)
}

func (r *astService) Wait() {
	<-r.ctx.Done()
}

func (r *astService) handleMsg(msg *logic.Call, wf logic.ResultFunc) {

	do, err := r.store.get(msg.Service)
	if err != nil {
		wf.Result(&logic.CallResult{
			Err: berr.Warp("app.Server", "handle get func", err),
		})
		return
	}

	ctx := context.WithValue(r.ctx, "info", map[string]string{"reqID": wf.ReqID, "connID": wf.ConnID})

	data, err := do(ctx, msg.Fun, msg.Param)
	if err != nil {
		wf.Result(&logic.CallResult{
			Err: berr.Warp("app.Server", "handle", err),
		})
		return
	}

	wf.Result(&logic.CallResult{Result: data})
}
