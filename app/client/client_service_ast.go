package client

import (
	"context"
	"github.com/MashiroC/begonia/logic"
)

type astService struct {
	name string
	c    *rClient
}

func (r *rClient) newAstService(name string) *astService {
	return &astService{
		name: name,
		c:    r,
	}
}

func (a *astService) FuncSync(name string) (rf RemoteFunSync, err error) {
	rf = func(params ...interface{}) (result interface{}, err error) {
		var ctx context.Context
		var p []byte

		switch params[0].(type) {
		case context.Context:
			ctx = params[0].(context.Context)
			p = params[1].([]byte)
		case []byte:
			ctx = context.TODO()
			p = params[0].([]byte)
		}

		res := a.c.lg.CallSync(ctx, &logic.Call{
			Service: a.name,
			Fun:     name,
			Param:   p,
		})
		result = res.Result
		err = res.Err
		return
	}
	return
}

func (a *astService) FuncAsync(name string) (rf RemoteFunAsync, err error) {

	rf = func(callback AsyncCallback, params ...interface{}) {
		var ctx context.Context
		var p []byte

		switch params[0].(type) {
		case context.Context:
			ctx = params[0].(context.Context)
			p = params[1].([]byte)
		case []byte:
			ctx = context.TODO()
			p = params[0].([]byte)
		}
		// 对于代码生成 传进来的就只有bytes了

		a.c.lg.CallAsync(ctx, &logic.Call{
			Service: a.name,
			Fun:     name,
			Param:   p,
		}, func(result *logic.CallResult) {
			if result.Err != nil {
				callback(nil, result.Err)
				return
			}
			// 对出参解码

			callback(result.Result, err)
		})
	}

	return
}
