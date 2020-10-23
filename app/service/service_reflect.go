package service

// service_reflect.go 反射实现的api

import (
	"begonia2/app/coding"
	"begonia2/app/core"
	"begonia2/logic"
	"begonia2/tool/qconv"
	"begonia2/tool/reflects"
	"context"
	"log"
	"reflect"
)

// rService 反射的 reflect service api
type rService struct {
	lg     logic.Service
	ctx    context.Context
	cancel context.CancelFunc

	store *serviceStore
}

func (r *rService) Register(name string, service interface{}) {

	// TODO:注册后 把函数注册到本地
	fs, ms := coding.Parse("avro", service)

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
			in:     inCoder,
			out:    outCoder,
			obj:    service,
			method: ms[i],
		})
	}

	res := r.lg.CallSync(core.Call.Register(name, fs))
	// TODO:handler error
	if res.Err != "" {
		panic(res.Err)
	}

	var ok bool
	success.DecodeIn(res.Result, &ok)

	if ok {
		return
	}
}

func (r *rService) Wait() {
	<-r.ctx.Done()
}

func (r *rService) work() {
	for {
		msg, wf := r.lg.RecvCall()

		go r.handleMsg(msg, wf)

	}
}

func (r *rService) handleMsg(msg *logic.Call, wf logic.ResultFunc) {
	fun, err := r.store.get(msg.Service, msg.Fun)
	if err != nil {
		log.Println("get fun err")
		wf.Result(&logic.CallResult{
			Err: "get fun err",
		})
		return
	}
	data, err := fun.in.Decode(msg.Param)
	if err != nil {
		log.Println("decode err")
		wf.Result(&logic.CallResult{
			Err: "decode error",
		})
		return
	}

	//TODO:这个反射调用后面再想办法改改
	inVal := []reflect.Value{reflect.ValueOf(fun.obj)}
	inVal = append(inVal, reflects.ToValue(data.(map[string]interface{}))...)

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
		// 这个error不应该有的
		panic(err)
	}

	wf.Result(&logic.CallResult{Result: b})
}