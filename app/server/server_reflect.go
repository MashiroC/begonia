package server

// server_reflect.go 反射实现的api

import (
	"context"
	"github.com/MashiroC/begonia/core"
	"github.com/MashiroC/begonia/internal/coding"
	"github.com/MashiroC/begonia/logic"
	"github.com/MashiroC/begonia/tool/berr"
	"github.com/MashiroC/begonia/tool/qconv"
	"github.com/MashiroC/begonia/tool/reflects"
	"reflect"
)

// rService 反射的 reflect server api
type rService struct {
	lg     *logic.Service
	ctx    context.Context
	cancel context.CancelFunc

	store           *serviceStore
	isLocalRegister bool
}

func (r *rService) Register(name string, service interface{}) {

	fs, ms, reSharps := coding.Parse("avro", service)

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
	}

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

func (r *rService) Wait() {
	<-r.ctx.Done()
}

func (r *rService) handleMsg(msg *logic.Call, wf logic.ResultFunc) {

	if r.isLocalRegister && msg.Service == core.ServiceName && msg.Fun == "ServiceInfo" {
		res, err := core.C.Invoke("", "", "ServiceInfo", msg.Param)
		wf.Result(&logic.CallResult{
			Result: res,
			Err:    err,
		})
	}

	fun, err := r.store.get(msg.Service, msg.Fun)
	if err != nil {
		wf.Result(&logic.CallResult{
			Err: berr.Warp("app.server", "handle get func", err),
		})
		return
	}
	data, err := fun.in.Decode(msg.Param)
	if err != nil {
		wf.Result(&logic.CallResult{
			Err: berr.Warp("app.server", "handle", err),
		})
		return
	}

	//TODO:这个反射调用后面再想办法改改
	inVal := []reflect.Value{reflect.ValueOf(fun.obj)}
	if fun.hasContext {
		ctx := context.WithValue(r.ctx, "info", map[string]string{"reqID": wf.ReqID, "connID": wf.ConnID})
		inVal = append(inVal, reflect.ValueOf(ctx))
	}

	inVal = append(inVal, reflects.ToValue(data.(map[string]interface{}), fun.reSharp)...)

	outVal := fun.method.Func.Call(inVal)

	m := reflects.FromValue(outVal)
	lastKey := "out" + qconv.I2S(len(outVal))
	v := m[lastKey]
	if vErr, ok := v.(error); ok {
		delete(m, lastKey)
		m["err"] = vErr.Error()
	} else {
		m["err"] = nil
	}

	b, err := fun.out.Encode(m)
	if err != nil {
		// out的schema是解析的函数，这里不应该有error，如果有直接panic出来，然后去修
		panic(err)
	}

	wf.Result(&logic.CallResult{Result: b})
}
