package client

import (
	"context"
	"fmt"
	"github.com/MashiroC/begonia/dispatch"
	"github.com/MashiroC/begonia/logic"
	"time"
)

// starter.go something

// BootStartByCenter 根据center cluster模式启动
func BootStartByCenter(optionMap map[string]interface{}) *rClient {

	fmt.Println("  ____                              _        \n |  _ \\                            (_)       \n | |_) |  ___   __ _   ___   _ __   _   __ _ \n |  _ <  / _ \\ / _` | / _ \\ | '_ \\ | | / _` |\n | |_) ||  __/| (_| || (_) || | | || || (_| |\n |____/  \\___| \\__, | \\___/ |_| |_||_| \\__,_|\n                __/ |                        \n               |___/                         ")

	ctx, cancel := context.WithCancel(context.Background())
	c := &rClient{
		ctx:    ctx,
		cancel: cancel,
	}

	// TODO:给dispatch初始化

	optionMap["pingpongTime"] = 10 * time.Second

	var dp dispatch.Dispatcher
	dp = dispatch.NewLinkedByDefaultCluster()

	if err := dp.Link(optionMap); err != nil {
		panic(err)
	}

	c.lg = logic.NewClient(dp)

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
