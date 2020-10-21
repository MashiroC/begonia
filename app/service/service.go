// Time : 2020/9/19 16:02
// Author : Kieran

// client
package service

import (
	"begonia2/app/core"
	"begonia2/logic"
	"begonia2/opcode/coding"
	"begonia2/tool/reflects"
	"context"
	"log"
	"reflect"
)

var success coding.SuccessCoder

// service.go something

// Service 服务端的接口
type Service interface {
	Register(name string, service interface{})
	Wait()
}

// rService 反射的 reflect service api
type rService struct {
	lg     logic.Service
	ctx    context.Context
	cancel context.CancelFunc

	coders *coderSet
}

func (r *rService) Register(name string, service interface{}) {
	//coder := coding.AvroCoder

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
		r.coders.store(name, f.Name, reflectFun{
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
		msg, wf := r.lg.RecvMsg()

		go r.handleMsg(msg, wf)

	}
}

func (r *rService) handleMsg(msg *logic.Call, wf logic.WriteFunc) {
	fun,err := r.coders.get(msg.Service, msg.Fun)
	if err != nil {
		log.Println("get fun err")
		wf(&logic.CallResult{
			Err: "get fun err",
		})
		return
	}
	data, err := fun.in.Decode(msg.Param)
	if err != nil {
		log.Println("decode err")
		wf(&logic.CallResult{
			Err: "decode error",
		})
		return
	}

	//TODO:这个反射调用后面再想办法改改
	inVal := []reflect.Value{reflect.ValueOf(fun.obj)}
	inVal = append(inVal, reflects.ToValue(data.(map[string]interface{}))...)

	outVal := fun.method.Func.Call(inVal)

	m := reflects.FromValue(outVal)

	b, err := fun.out.Encode(m)
	if err != nil {
		// 这个error不应该有的
		panic(err)
	}

	wf(&logic.CallResult{Result: b})
}

// astService ast树代码生成的ast service api
type astService struct {
}
