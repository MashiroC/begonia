// Time : 2020/9/26 17:20
// Author : Kieran

// appservice
package service

import (
	"begonia2/dispatch"
	"begonia2/logic"
)

// starter.go something

// bootStartByManager 根据manager cluster模式启动
func bootStartByManager(optionMap map[string]interface{}) Service {
	s := &rService{}

	s.coders = newCoderSet()
	var addr string
	if addrIn, ok := optionMap["managerAddr"]; ok {
		addr = addrIn.(string)
	}

	var dp dispatch.Dispatcher
	dp = dispatch.NewByCenterCluster()
	dp.Link(addr)

	s.lg = logic.NewService(dp)

	go s.work()
	//TODO: 发一个包，拉取配置

	// 假设这个getConfig是sub service的一个远程函数
	//var getConfig = func(...interface{}) (interface{}, error) {
	//	return map[string]interface{}{}, nil
	//}
	//
	//// 假设m就是拿到的远程配置
	//m, err := getConfig()
	//
	//// TODO:根据拿到的远程配置来修改配置
	//// do some thing
	//fmt.Println(m, err)
	// 修改配置之前的一系列调用全部都是按默认配置来的

	return s
}

// New 初始化，获得一个service对象，传入一个mode参数，以及一个option的不定参数
func New(mode string, optionFunc ...OptionFunc) (s Service) {
	optionMap := defaultServiceConfig()

	for _, f := range optionFunc {
		f(optionMap)
	}

	switch mode {
	case "center":
		s = bootStartByManager(optionMap)
		// TODO:其他的模式和模式出问题的判断
	}

	return
}

func defaultServiceConfig() map[string]interface{} {
	m := make(map[string]interface{})

	// TODO:加入配置

	return m
}

type OptionFunc func(optionMap map[string]interface{})

func ManagerAddr(addr string) OptionFunc {
	return OptionFunc(func(optionMap map[string]interface{}) {
		optionMap["managerAddr"] = addr
	})
}
