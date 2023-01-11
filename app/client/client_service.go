package client

import (
	"context"
	"fmt"
	"github.com/MashiroC/begonia/app/coding"
	"github.com/MashiroC/begonia/logic"
	"github.com/MashiroC/begonia/tool/reflects"
)

// client_service.go something

// Service 客户端获取的远程服务的抽象
type Service interface {
	// 同步
	FuncSync(name string) (RemoteFunSync, error)
	// 异步
	FuncAsync(name string) (RemoteFunAsync, error)
}

// RemoteFunSync 同步远程函数
type RemoteFunSync func(params ...interface{}) (result interface{}, err error)

// RemoteFunAsync 异步远程函数
type RemoteFunAsync func(callback AsyncCallback, params ...interface{})

// AsyncCallback 异步回调
type AsyncCallback = func(interface{}, error)

type rService struct {
	name string
	funs map[string]Fun
	c    *rClient
}

func (r*rService) FuncSubscribe(key string) {
	// server

}

func (r *rService) FuncSync(name string) (rf RemoteFunSync, err error) {
	fun, exist := r.funs[name]

	rf = func(params ...interface{}) (result interface{}, err error) {
		if r.funs == nil {
			err = fmt.Errorf("nil service, please check if get service [%s] success", r.name)
			return
		}

		if !exist {
			fun, exist = r.funs[name]
			if !exist {
				err = fmt.Errorf("nil func, fun [%s] not found", name)
				return
			}
		}

		ch := make(chan *logic.CallResult)

		ctx := context.TODO()
		if len(params) > 0 {
			if v, ok := params[0].(context.Context); ok {
				ctx = v
				params = params[1:]
			}
		}

		b, err := fun.InCoder.Encode(coding.ToAvroObj(params))
		if err != nil {
			err = fmt.Errorf("input type error: %w", err)
			return
		}

		r.c.lg.CallAsync(ctx, &logic.Call{
			Service: r.name,
			Fun:     name,
			Param:   b,
		}, func(res *logic.CallResult) {
			ch <- res
		})

		tmp := <-ch
		if tmp.Err != nil {
			err = tmp.Err
			return
		}

		// 对出参解码
		out, err := fun.OutCoder.Decode(tmp.Result)

		result = reflects.ToInterfaces(out.(map[string]interface{}))
		return
	}
	return
}

func (r *rService) FuncAsync(name string) (rf RemoteFunAsync, err error) {

	fun, exist := r.funs[name]

	rf = func(callback AsyncCallback, params ...interface{}) {

		if !exist {
			fun, exist = r.funs[name]
			if !exist {
				err = fmt.Errorf("nil func, fun [%s] not found", name)
				return
			}
		}
		// 对入参编码
		b, err := fun.InCoder.Encode(coding.ToAvroObj(params))
		if err != nil {
			//TODO: 当传入参数和要求类型不符时的错误返回
			panic(err)
		}

		ctx := context.TODO()
		if len(params) > 0 {
			if v, ok := params[0].(context.Context); ok {
				ctx = v
				params = params[1:]
			}
		}

		r.c.lg.CallAsync(ctx, &logic.Call{
			Service: r.name,
			Fun:     name,
			Param:   b,
		}, func(result *logic.CallResult) {
			if result.Err != nil {
				callback(nil, result.Err)
				return
			}
			// 对出参解码
			out, err := fun.OutCoder.Decode(result.Result)

			res := reflects.ToInterfaces(out.(map[string]interface{}))

			callback(res, err)
		})
	}

	return
}