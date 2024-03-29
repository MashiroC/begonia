package server

// server_ast.go ast实现的api

import (
	"context"
	"fmt"
	coreRegister "github.com/MashiroC/begonia/core/register"
	"github.com/MashiroC/begonia/internal/register"
	"github.com/MashiroC/begonia/logic"
	"log"
)

type astDo = func(ctx context.Context, fun string, param []byte) (result []byte, err error)

// CodeGenService 代码生成实现的服务
type CodeGenService interface {

	// Do 调用服务
	Do(ctx context.Context, fun string, param []byte) (result []byte, err error)

	// FuncList 返回要注册的函数
	FuncList() []coreRegister.FunInfo
}

// astServer ast树代码生成的ast Server api
type astServer struct {
	lg     *logic.Service
	ctx    context.Context
	cancel context.CancelFunc

	store    *astServiceStore
	register register.Register
}

func (r *astServer) Register(name string, service interface{}, registerFunc ...string) {

	cgs, ok := service.(CodeGenService)
	if !ok {
		panic("please use code-gen")
	}

	// local store
	if err := r.store.store(name, cgs.Do); err != nil {
		panic(err)
	}

	fs := cgs.FuncList()


	if err := r.register.Register(name, fs); err != nil {
		log.Println("register func in start error:", err)
	}

	var flag bool
	if !r.register.IsLocal() {
		r.lg.Hook("dispatch.link", func(connID string) {
			if !flag {
				return
			}
			if err := r.register.Register(name, fs); err != nil {
				log.Println("register func in restart error:", err)
			}
		})
	}
	flag = true
}

func (r *astServer) Wait() {
	<-r.ctx.Done()
}

func (r *astServer) handleMsg(ctx context.Context, msg *logic.Call, wf logic.ResultFunc) {

	do, err := r.store.get(msg.Service)
	if err != nil {
		wf(&logic.CallResult{
			Err: fmt.Errorf("app.Server store get error: %w", err),
		})
		return
	}

	data, err := do(ctx, msg.Fun, msg.Param)
	if err != nil {
		wf(&logic.CallResult{
			Err: err,
		})
		return
	}

	wf(&logic.CallResult{Result: data})
}
