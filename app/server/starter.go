package server

import (
	"context"
	"fmt"
	"github.com/MashiroC/begonia/app"
	cRegister "github.com/MashiroC/begonia/core/register"
	"github.com/MashiroC/begonia/dispatch"
	"github.com/MashiroC/begonia/internal/register"
	"github.com/MashiroC/begonia/logic"
	"log"
)

// starter.go something

// BootStartByManager 根据manager cluster模式启动
func BootStartByManager(optionMap map[string]interface{}) (s Server) {

	fmt.Println("  ____                              _        \n |  _ \\                            (_)       \n | |_) |  ___   __ _   ___   _ __   _   __ _ \n |  _ <  / _ \\ / _` | / _ \\ | '_ \\ | | / _` |\n | |_) ||  __/| (_| || (_) || | | || || (_| |\n |____/  \\___| \\__, | \\___/ |_| |_||_| \\__,_|\n                __/ |                        \n               |___/                         ")

	log.Printf("begonia Server start with [%s] mode\n", app.ServiceAppMode)

	ctx, cancel := context.WithCancel(context.Background())
	var isLocal bool
	// 读配置
	var addr string
	if addrIn, ok := optionMap["addr"]; ok {
		addr = addrIn.(string)
	} else {
		panic("addr must exist")
	}

	var isP2P bool
	if dpTyp, ok := optionMap["dpTyp"]; ok && dpTyp == "p2p" {
		isP2P = true
	}

	// 创建 dispatch
	var dp dispatch.Dispatcher
	if isP2P {
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

	var waitChans *logic.CallbackStore
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

	coreRegister := cRegister.NewCoreRegister()

	var rg register.Register
	if isP2P {
		rg = register.NewLocalRegister(coreRegister)

	} else {
		rg = register.NewRemoteRegister(lg.Client)
	}

	// 创建实例
	if app.ServiceAppMode == app.Ast {
		ast := &astService{}
		ast.ctx = ctx
		ast.cancel = cancel

		ast.lg = lg
		ast.lg.HandleRequest = ast.handleMsg

		// 创建服务存储的数据结构
		ast.store = newAstServiceStore()

		ast.register = rg

		s = ast
	} else {
		r := &rService{}
		r.ctx = ctx
		r.cancel = cancel

		r.lg = lg
		r.lg.HandleRequest = r.handleMsg
		r.isLocalRegister = isLocal

		// 创建服务存储的数据结构
		r.store = newServiceStore()

		r.register = rg

		s = r
	}

	if isP2P {
		s.Register("REGISTER", coreRegister, "Register", "ServiceInfo")
		optionMap["REGISTER"] = coreRegister
	}

	return s
}

func GetLogic(s Server) *logic.Service {
	switch in := s.(type) {
	case *astService:
		return in.lg
	case *rService:
		return in.lg
	default:
		panic("error")
	}
}
