package begonia

import (
	"github.com/MashiroC/begonia/app"
	"github.com/MashiroC/begonia/app/client"
	"github.com/MashiroC/begonia/app/option"
)

// Client 客户端的接口
type Client = client.Client

type RemoteFunc = client.RemoteFunSync
type RemoteFuncAsync = client.RemoteFunAsync

// NewClient 初始化，获得一个service对象，传入一个mode参数，以及一个option的不定参数
func NewClient(optionFunc ...option.WriteFunc) (cli Client) {
	optionMap := defaultClientConfig()

	for _, f := range optionFunc {
		f(optionMap)
	}

	cli = client.BootStartByCenter(optionMap)

	return
}

func NewClientWithAst(optionFunc ...option.WriteFunc) (cli Client) {
	app.ServiceAppMode = app.Ast

	return NewClient(optionFunc...)
}

func defaultClientConfig() map[string]interface{} {
	m := make(map[string]interface{})

	// TODO:加入默认配置
	m["addr"] = ":12306"

	return m
}
