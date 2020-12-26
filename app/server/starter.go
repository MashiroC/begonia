package server

import (
	"context"
	"fmt"
	"github.com/MashiroC/begonia/core"
	"github.com/MashiroC/begonia/dispatch"
	"github.com/MashiroC/begonia/internal"
	"github.com/MashiroC/begonia/logic"
	"log"
)

// starter.go something

// BootStartByManager 根据manager cluster模式启动
func BootStartByManager(optionMap map[string]interface{}) Server {

	fmt.Println("  ____                              _        \n |  _ \\                            (_)       \n | |_) |  ___   __ _   ___   _ __   _   __ _ \n |  _ <  / _ \\ / _` | / _ \\ | '_ \\ | | / _` |\n | |_) ||  __/| (_| || (_) || | | || || (_| |\n |____/  \\___| \\__, | \\___/ |_| |_||_| \\__,_|\n                __/ |                        \n               |___/                         ")

	log.Printf("begonia Server start with [%s] mode\n", internal.ServiceAppMode)

	ctx, cancel := context.WithCancel(context.Background())
	var isLocal bool
	// 读配置
	var addr string
	if addrIn, ok := optionMap["addr"]; ok {
		addr = addrIn.(string)
	} else {
		panic("addr must exist")
	}

	// 创建 dispatch
	var dp dispatch.Dispatcher
	if dpTyp, ok := optionMap["dpTyp"]; ok && dpTyp == "p2p" {
		log.Printf("begonia Server will listen on [%s]", addr)
		dp = dispatch.NewSetByDefaultCluster()
		go dp.Listen(addr)
		isLocal = true
	} else {
		log.Printf("begonia Server will link to [%s]", addr)
		dp = dispatch.NewLinkedByDefaultCluster()
		if err := dp.Link(addr); err != nil {
			panic(err)
		}
	}

	var waitChans *logic.WaitChans
	waitChans = logic.NewWaitChans()

	// 创建 logic
	var lg *logic.Service
	lg = logic.NewService(dp, waitChans)

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
	if internal.ServiceAppMode == "ast" {
		s := &astService{}
		s.ctx = ctx
		s.cancel = cancel

		s.lg = lg
		s.lg.HandleRequest = s.handleMsg
		s.isLocalRegister = isLocal

		// 创建服务存储的数据结构
		s.store = newAstServiceStore()

		//if isLocal {
		//	core.C = core.NewSubService()
		//	s.Register(core.ServiceName, &Core{})
		//}

		return s
	} else {
		s := &rService{}
		s.ctx = ctx
		s.cancel = cancel

		s.lg = lg
		s.lg.HandleRequest = s.handleMsg
		s.isLocalRegister = isLocal

		// 创建服务存储的数据结构
		s.store = newServiceStore()

		//if isLocal {
		//	core.C = core.NewSubService()
		//	s.Register(core.ServiceName, &Core{})
		//}

		return s
	}

}

type Core struct {
}

func (c *Core) ServiceInfo(serviceName string) []byte {
	s := core.Call.ServiceInfo(serviceName)
	res, err := core.C.Invoke("", "", s.Fun, s.Param)
	if err != nil {
		panic(err)
	}
	return res
}
