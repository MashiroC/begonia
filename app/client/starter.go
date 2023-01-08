package client

import (
	"context"
	"fmt"
	"github.com/MashiroC/begonia/app"
	"github.com/MashiroC/begonia/config"
	"github.com/MashiroC/begonia/dispatch"
	"github.com/MashiroC/begonia/internal/register"
	"github.com/MashiroC/begonia/logic"
	"github.com/MashiroC/begonia/tool/retry"
	"github.com/MashiroC/begonia/tracing"
	"log"
)

// starter.go something

// BootStartByCenter 根据center cluster模式启动
func BootStartByCenter(optionMap map[string]interface{}) *rClient {

	fmt.Println("  ____                              _        \n |  _ \\                            (_)       \n | |_) |  ___   __ _   ___   _ __   _   __ _ \n |  _ <  / _ \\ / _` | / _ \\ | '_ \\ | | / _` |\n | |_) ||  __/| (_| || (_) || | | || || (_| |\n |____/  \\___| \\__, | \\___/ |_| |_||_| \\__,_|\n                __/ |                        \n               |___/                         ")

	mode := app.ParseMode(optionMap)
	log.Printf("begonia client start with [%s] mode\n", mode.String())

	ctx, cancel := context.WithCancel(context.Background())
	c := &rClient{
		mode:   mode,
		ctx:    ctx,
		cancel: cancel,
	}

	// TODO:给dispatch初始化

	var addr string
	if addrIn, ok := optionMap["addr"]; ok {
		addr = addrIn.(string)
	} else {
		panic("addr must exist")
	}

	log.Printf("begonia client will link to [%s]\n", addr)

	var dp dispatch.Dispatcher
	dp = dispatch.NewLinkedByDefaultCluster()

	err := retry.Do("start", func() (ok bool) {
		err := dp.Start(addr)
		if err != nil {
			log.Println("dp start error: ", err)
		}
		return err == nil
	},
		config.C.Dispatch.ConnectionRetryCount,
		config.C.Dispatch.ConnectionIntervalSecond)
	if err != nil {
		panic(err)
	}

	var tracer tracing.Tracer
	if in, ok := optionMap["tracing"]; ok {
		tracer = in.(tracing.Tracer)
	}

	c.lg = logic.NewClient(dp, tracer)

	c.register = register.NewRemoteRegister(c.lg)

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
		// 修改配置之前的一系列调用全部都是按默认配置来的
	*/

	return c
}

func BootStartWithLogic(optionMap map[string]interface{}, lg *logic.Client) *rClient {
	fmt.Println("  ____                              _        \n |  _ \\                            (_)       \n | |_) |  ___   __ _   ___   _ __   _   __ _ \n |  _ <  / _ \\ / _` | / _ \\ | '_ \\ | | / _` |\n | |_) ||  __/| (_| || (_) || | | || || (_| |\n |____/  \\___| \\__, | \\___/ |_| |_||_| \\__,_|\n                __/ |                        \n               |___/                         ")

	log.Printf("begonia client start with [%s] mode\n", app.ServiceAppMode)

	ctx, cancel := context.WithCancel(context.Background())
	c := &rClient{
		ctx:      ctx,
		cancel:   cancel,
		register: register.NewRemoteRegister(lg),
	}

	// TODO:给dispatch初始化

	c.lg = lg

	c.lg.Handle("dispatch.frame", c.lg.DpHandler)

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
		// 修改配置之前的一系列调用全部都是按默认配置来的
	*/

	return c
}
