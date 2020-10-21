// Time : 2020/10/10 21:32
// Author : Kieran

// center
package center

import (
	"begonia2/app/option"
	"begonia2/config"
	"begonia2/dispatch"
	"begonia2/logic"
	"log"
)

// starter.go something
// bootStartByCenter 根据center cluster模式启动
func bootstart(optionMap map[string]interface{}) Center {

	//ctx, cancel := context.WithCancel(context.Background())
	c := &rCenter{
		//ctx:    ctx,
		//cancel: cancel,
		services: newServiceSet(),
	}

	// TODO:给dispatch初始化

	var addr string
	if addrIn, ok := optionMap["managerAddr"]; ok {
		addr = addrIn.(string)
	}

	var dp dispatch.Dispatcher
	dp = dispatch.NewByCenterCluster()
	go dp.Listen(addr)

	var rs *logic.ReqSet
	rs = logic.NewReqSet(config.C.Logic.RequestTimeOut)

	c.rs = rs

	c.lg = logic.NewMixWithReqSet(dp, rs)

	log.Println("begonia center started")
	//TODO: 发一个包，拉取配置

	/*

		先不去拉配置 后面再加

		// 假设这个getConfig是sub service的一个远程函数
		var getConfig = func(...interface{}) (interface{}, error) {
			return map[string]interface{}{}, nil
		}

		// 假设m就是拿到的远程配置
		m, err := getConfig()

		// TODO:根据拿到的远程配置来修改配置
		// do some thing
		fmt.Println(m, err)
		// 修改配置之前的一系列调用全部都是按默认配置来的
	*/
	return c
}

// New 初始化，获得一个service对象，传入一个mode参数，以及一个option的不定参数
func New(mode string, optionFunc ...option.OptionFunc) (cli Center) {
	optionMap := defaultClientConfig()

	for _, f := range optionFunc {
		f(optionMap)
	}

	switch mode {
	case "center":
		cli = bootstart(optionMap)
		// TODO:其他的模式和模式出问题的判断
	}

	return
}

func defaultClientConfig() map[string]interface{} {
	m := make(map[string]interface{})

	// TODO:加入配置

	return m
}
