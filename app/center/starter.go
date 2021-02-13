// Package bgacenter 服务中心
package center

import (
	"context"
	"github.com/MashiroC/begonia/app"
	"github.com/MashiroC/begonia/app/option"
	"github.com/MashiroC/begonia/app/server"
	cRegister "github.com/MashiroC/begonia/core/register"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/internal/proxy"
	"github.com/MashiroC/begonia/logic"
	"log"
)

// bootstart 根据center cluster模式启动
func bootstart(optionMap map[string]interface{}) server.Server {

	app.ServiceAppMode = app.Ast

	s := server.BootStart(optionMap)

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

	lg.Dp.Hook("close", coreRegister.HandleConnClose)

	lg.Dp.Handle("proxy", p)
	// ========== END ==========

	log.Println("begonia bgacenter started")
	//TODO: 发一个包，拉取配置

	return s
}

// New 拿到一个server，该server中会初始化center相关的东西
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
