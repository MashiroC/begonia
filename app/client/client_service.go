// Time : 2020/9/19 16:03
// Author : Kieran

// client
package client

import (
	"begonia2/app"
	"begonia2/logic"
	"errors"
	"fmt"
)

// client_service.go something
type Service interface {
	// 同步
	FuncSync(name string) (RemoteFunSync, error)
	// 异步
	FuncAsync(name string) (RemoteFunAsync, error)
}

type RemoteFunSync func(params ...interface{}) (result interface{}, err error)

type RemoteFunAsync func(callback AsyncCallback, params ...interface{})

type AsyncCallback = func(interface{}, error)

type rService struct {
	name string
	funs map[string]app.FunInfo
	c    *rClient
}

func (r *rService) FuncSync(name string) (rf RemoteFunSync, err error) {
	f, exist := r.funs[name]
	if !exist {
		err = fmt.Errorf("remote func [%s] not exist!", name)
		return
	}

	rf = func(params ...interface{}) (result interface{}, err error) {
		ch := make(chan *logic.CallResult)

		b, err := f.InCoder.Encode(result)
		if err != nil {
			//TODO: 当传入参数和要求类型不符时的错误返回
			panic(err)
		}

		r.c.lg.CallAsync(&logic.Call{
			Service: r.name,
			Fun:     name,
			Param:   b,
		}, func(res *logic.CallResult) {
			ch <- res
		})

		tmp := <-ch
		if tmp.Err != "" {
			err = errors.New(tmp.Err)
			return
		}

		return f.OutCoder.Decode(tmp.Result)
	}
	return
}

func (r *rService) FuncAsync(name string) (rf RemoteFunAsync, err error) {

	f, exist := r.funs[name]
	if !exist {
		err = fmt.Errorf("remote func [%s] not exist!", name)
		return
	}

	rf = func(callback AsyncCallback, params ...interface{}) {

		// 对入参编码
		b, err := f.InCoder.Encode(params)
		if err != nil {
			//TODO: 当传入参数和要求类型不符时的错误返回
			panic(err)
		}

		r.c.lg.CallAsync(&logic.Call{
			Service: r.name,
			Fun:     name,
			Param:   b,
		}, func(result *logic.CallResult) {
			if result.Err != "" {
				callback(nil, errors.New(result.Err))
				return
			}
			// 对出参解码
			in, err := f.OutCoder.Decode(result.Result)
			callback(in, err)
		})
	}

	return
}
