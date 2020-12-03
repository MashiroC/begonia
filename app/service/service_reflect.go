package service

// service_reflect.go 反射实现的api

import (
	"context"
	"github.com/MashiroC/begonia/app/coding"
	"github.com/MashiroC/begonia/core"
	"github.com/MashiroC/begonia/logic"
	"github.com/MashiroC/begonia/tool/berr"
	"github.com/MashiroC/begonia/tool/qconv"
	"github.com/MashiroC/begonia/tool/reflects"
	"reflect"
)

// rService 反射的 reflect service api
type rService struct {
	lg     *logic.Service
	ctx    context.Context
	cancel context.CancelFunc

	store *serviceStore
}

func (r *rService) Register(name string, service interface{}) {

	// TODO:注册后 把函数注册到本地
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
			in:      inCoder,
			out:     outCoder,
			obj:     service,
			reSharp: reSharps[i],
			method:  ms[i],
		})
	}

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

func (r *rService) Wait() {
	<-r.ctx.Done()
}

func (r *rService) handleMsg(msg *logic.Call, wf logic.ResultFunc) {
	fun, err := r.store.get(msg.Service, msg.Fun)
	if err != nil {
		wf.Result(&logic.CallResult{
			Err: berr.Warp("app.service", "handle get func", err),
		})
		return
	}
	data, err := fun.in.Decode(msg.Param)
	if err != nil {
		wf.Result(&logic.CallResult{
			Err: berr.Warp("app.service", "handle", err),
		})
		return
	}

	//TODO:这个反射调用后面再想办法改改
	inVal := []reflect.Value{reflect.ValueOf(fun.obj)}
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
