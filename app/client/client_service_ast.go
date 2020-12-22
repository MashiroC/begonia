package client

import (
	"github.com/MashiroC/begonia/logic"
)

type astService struct {
	name string
	c    *rClient
}

func newAstService(name string,c *rClient) *astService {
	return &astService{
		name: name,
		c:    c,
	}
}

func (a *astService) FuncSync(name string) (rf RemoteFunSync, err error) {
	rf = func(params ...interface{}) (result interface{}, err error) {
		res := a.c.lg.CallSync(&logic.Call{
			Service: a.name,
			Fun:     name,
			Param:   params[0].([]byte),
		})
		result = res.Result
		err = res.Err
		return
	}
	return
}

func (a *astService) FuncAsync(name string) (rf RemoteFunAsync, err error) {

	rf = func(callback AsyncCallback, params ...interface{}) {

		// 对于代码生成 传进来的就只有bytes了

		a.c.lg.CallAsync(&logic.Call{
			Service: a.name,
			Fun:     name,
			Param:   params[0].([]byte),
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
