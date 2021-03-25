package server

// server_ast.go ast实现的api

import (
	"context"
	"fmt"
	coreRegister "github.com/MashiroC/begonia/core/register"
	"github.com/MashiroC/begonia/internal/logger"
	"github.com/MashiroC/begonia/internal/register"
	"github.com/MashiroC/begonia/logic"
	"github.com/MashiroC/begonia/tool/log"
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

	logService logger.LoggerService // 日志中心服务
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

	r.register.Register(name, fs)
	// 注册进日志服务中心
	if r.logService != nil {
		l:=log.DefaultNewLogger()
		l.Info(name+" join in")
		r.logService.Save(name, l.SetFields(log.Fields{"server": name}))

	}
}

func (r *astServer) Wait() {
	<-r.ctx.Done()
}

func (r *astServer) handleMsg(msg *logic.Call, wf logic.ResultFunc) {

	do, err := r.store.get(msg.Service)
	if err != nil {
		wf.Result(&logic.CallResult{
			Err: fmt.Errorf("app.Server store get error: %w", err),
		})
		return
	}

	ctx := context.WithValue(r.ctx, "info", map[string]string{"reqID": wf.ReqID, "connID": wf.ConnID})

	data, err := do(ctx, msg.Fun, msg.Param)
	if err != nil {
		wf.Result(&logic.CallResult{
			Err: fmt.Errorf("app.Server handle error: %w", err),
		})
		return
	}

	wf.Result(&logic.CallResult{Result: data})
}
