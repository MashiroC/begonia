package service

import (
	"context"
	"fmt"
	"github.com/MashiroC/begonia/dispatch"
	"github.com/MashiroC/begonia/logic"
)

// starter.go something

// BootStartByManager 根据manager cluster模式启动
func BootStartByManager(optionMap map[string]interface{}) *rService {

	fmt.Println("  ____                              _        \n |  _ \\                            (_)       \n | |_) |  ___   __ _   ___   _ __   _   __ _ \n |  _ <  / _ \\ / _` | / _ \\ | '_ \\ | | / _` |\n | |_) ||  __/| (_| || (_) || | | || || (_| |\n |____/  \\___| \\__, | \\___/ |_| |_||_| \\__,_|\n                __/ |                        \n               |___/                         ")

	ctx, cancel := context.WithCancel(context.Background())

	// 读配置
	var addr string
	if addrIn, ok := optionMap["managerAddr"]; ok {
		addr = addrIn.(string)
	}

	// 创建 dispatch
	var dp dispatch.Dispatcher
	dp = dispatch.NewLinkedByDefaultCluster()

	if err := dp.Link(addr); err != nil {
		panic(err)
	}

	var waitChans *logic.WaitChans
	waitChans= logic.NewWaitChans()

	// 创建 logic
	var lg *logic.Service
	lg = logic.NewService(dp,waitChans)

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

	// 创建实例
	s := &rService{}
	s.ctx = ctx
	s.cancel = cancel

	s.lg = lg
	s.lg.HandleRequest=s.handleMsg

	// 创建服务存储的数据结构
	s.store = newServiceStore()

	return s
}
