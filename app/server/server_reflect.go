package server

// server_reflect.go 反射实现的api

import (
	"context"
	"errors"
	"fmt"
	"github.com/MashiroC/begonia/app/coding"
	cRegister "github.com/MashiroC/begonia/core/register"
	"github.com/MashiroC/begonia/internal/register"
	"github.com/MashiroC/begonia/logic"
	"github.com/MashiroC/begonia/tool/qconv"
	"github.com/MashiroC/begonia/tool/reflects"
	"log"
	"reflect"
)

// rServer 反射的 reflect Server api
type rServer struct {
	lg     *logic.Service
	ctx    context.Context
	cancel context.CancelFunc

	store           *serviceStore
	isLocalRegister bool

	register register.Register
}

func (r *rServer) Register(name string, service interface{}, registerFunc ...string) {

	fs, ms, reSharps := coding.Parse("avro", service, registerFunc)

	var registerFs []cRegister.FunInfo
	registerFs = make([]cRegister.FunInfo, 0, len(fs))

	for i, f := range fs {
		inCoder, err := coding.NewAvro(f.InSchema)
		if err != nil {
			panic(err)
		}
		outCoder, err := coding.NewAvro(f.OutSchema)
		if err != nil {
			panic(err)
		}
		r.store.store(name, f.Name, reflectFun{
			in:         inCoder,
			out:        outCoder,
			obj:        service,
			reSharp:    reSharps[i],
			method:     ms[i],
			hasContext: f.HasContext,
		})

		registerFs = append(registerFs, cRegister.FunInfo{
			Name:      f.Name,
			InSchema:  f.InSchema,
			OutSchema: f.OutSchema,
		})
	}

	//fmt.Println(registerFs)
	err := r.register.Register(name, registerFs)
	if err != nil {
		panic(err)
	}
}

func (r *rServer) Wait() {
	<-r.ctx.Done()
}

func (r *rServer) handleMsg(ctx context.Context,msg *logic.Call, wf logic.ResultFunc) {
	fun, err := r.store.get(msg.Service, msg.Fun)
	if err != nil {
		wf(&logic.CallResult{
			Err: fmt.Errorf("app.Server get func error: %w", err),
		})
		return
	}
	data, err := fun.in.Decode(msg.Param)
	if err != nil {
		wf(&logic.CallResult{
			Err: fmt.Errorf("app.Server handle error: %w", err),
		})
		return
	}

	//TODO:这个反射调用后面再想办法改改
	inVal := []reflect.Value{reflect.ValueOf(fun.obj)}
	if fun.hasContext {
		inVal = append(inVal, reflect.ValueOf(ctx))
	}

	inVal = append(inVal, reflects.ToValue(data.(map[string]interface{}), fun.reSharp)...)
	m, hasError := callWithRecover(fun.method, inVal)
	if hasError {
		wf(&logic.CallResult{Err: m["err"].(error)})
		return
	}

	b, err := fun.out.Encode(m)
	if err != nil {
		// out的schema是解析的函数，这里不应该有error，如果有直接panic出来，然后去修
		panic(err)
	}

	wf(&logic.CallResult{Result: b})
}

func callWithRecover(fun reflect.Method, inVal []reflect.Value) (m map[string]interface{}, hasError bool) {
	defer func() {
		if re := recover(); re != nil {
			hasError = true
			log.Printf("[RECOVER] panic in remote func call: %s\n",re)
			m["err"] = errors.New(fmt.Sprintf("server func call recover: %s", re))
		}
	}()
	outVal := fun.Func.Call(inVal)

	m = reflects.FromValue(outVal)
	lastKey := "F" + qconv.I2S(len(outVal))
	v := m[lastKey]
	if vErr, ok := v.(error); ok {
		m["err"] = vErr
		hasError = true
	}

	return
}
