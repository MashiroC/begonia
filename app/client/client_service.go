package client

import (
	"github.com/MashiroC/begonia/app"
	"github.com/MashiroC/begonia/app/coding"
	"github.com/MashiroC/begonia/logic/containers"
	"github.com/MashiroC/begonia/tool/berr"
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
	funs map[string]app.FunInfo
	c    *rClient
}

func (r *rService) FuncSync(name string) (rf RemoteFunSync, err error) {
	f, exist := r.funs[name]
	if !exist {
		err = berr.NewF("app.client", "get func", "remote func [%s] not exist", name)
		return
	}

	rf = func(params ...interface{}) (result interface{}, err error) {
		ch := make(chan *containers.CallResult)

		b, err := f.InCoder.Encode(coding.ToAvroObj(params))
		if err != nil {
			//TODO: 当传入参数和要求类型不符时的错误返回
			panic(err)
		}

		r.c.lg.CallAsync(&containers.Call{
			Service: r.name,
			Fun:     name,
			Param:   b,
		}, func(res *containers.CallResult) {
			ch <- res
		})

		tmp := <-ch
		if tmp.Err != nil {
			err = tmp.Err
			return
		}

		// 对出参解码
		out, err := f.OutCoder.Decode(tmp.Result)

		result = reflects.ToInterfaces(out.(map[string]interface{}))
		return
	}
	return
}

func (r *rService) FuncAsync(name string) (rf RemoteFunAsync, err error) {

	f, exist := r.funs[name]
	if !exist {
		err = berr.NewF("app.client", "get func async", "remote func [%s] not exist", name)
		return
	}

	rf = func(callback AsyncCallback, params ...interface{}) {

		// 对入参编码
		b, err := f.InCoder.Encode(params)
		if err != nil {
			//TODO: 当传入参数和要求类型不符时的错误返回
			panic(err)
		}

		r.c.lg.CallAsync(&containers.Call{
			Service: r.name,
			Fun:     name,
			Param:   b,
		}, func(result *containers.CallResult) {
			if result.Err != nil {
				callback(nil, result.Err)
				return
			}
			// 对出参解码
			out, err := f.OutCoder.Decode(result.Result)

			res := reflects.ToInterfaces(out.(map[string]interface{}))

			callback(res, err)
		})
	}

	return
}
