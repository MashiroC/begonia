package center

import (
	"context"
	"github.com/MashiroC/begonia/app/option"
	"github.com/MashiroC/begonia/app/server"
	cRegister "github.com/MashiroC/begonia/core/register"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/internal"
	"github.com/MashiroC/begonia/internal/proxy"
	"github.com/MashiroC/begonia/logic"
	"log"
)

// starter.go something
// bootStartByCenter 根据center cluster模式启动
func bootstart(optionMap map[string]interface{}) server.Server {

	internal.ServiceAppMode = internal.Ast

	s := server.BootStartByManager(optionMap)

	coreRegister := optionMap["REGISTER"].(*cRegister.CoreRegister)

	// ========== 初始化代理器 ==========

	p := proxy.NewHandler()

	lg := server.GetLogic(s)

	p.Check = func(connID string, f frame.Frame) (redirectConnID string, ok bool) {

		// Response不走proxy器
		if _, okk := f.(*frame.Response); okk {
			return
		}

		req := f.(*frame.Request)

		if req.Service != "REGISTER" {
			redirectConnID, ok = coreRegister.GetToID(req.Service)
			if !ok {
				panic("unknown bu ok error")
			}
		}
		return
	}

	p.AddAction(func(connID, redirectConnID string, f frame.Frame) {
		req := f.(*frame.Request)
		lg.Callbacks.AddCallback(context.TODO(), req.ReqID, func(result *logic.CallResult) {
			err := lg.Dp.SendTo(connID, frame.NewResponse(req.ReqID, result.Result, result.Err))
			// TODO: sendTo如果发送失败，加入到队列，这里先log一下
			if err != nil {
				log.Println(err)
			}
		})
	})

	p.AddAction(func(connID, redirectConnID string, f frame.Frame) {
		err := lg.Dp.SendTo(redirectConnID, f)
		if err != nil {
			log.Println(err)
		}
		// TODO:handler err not panic
		return
	})

	lg.Dp.Handle("proxy", p)

	// ========== END ==========

	// ========== 初始化核心服务 ==========

	log.Println("begonia center started")
	//TODO: 发一个包，拉取配置

	// ========== END ==========
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
		// 修改配置之前的一系列调用全部都是按默认配置来的
	*/

	return s
}

// New 初始化，获得一个service对象，传入一个mode参数，以及一个option的不定参数
func New(optionFunc ...option.WriteFunc) (s server.Server) {
	optionMap := defaultClientConfig()

	for _, f := range optionFunc {
		f(optionMap)
	}

	s = bootstart(optionMap)

	return
}

func defaultClientConfig() map[string]interface{} {
	m := make(map[string]interface{})

	m["dpTyp"] = "p2p"

	// TODO:加入配置

	return m
}
