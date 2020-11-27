package center

import (
	"fmt"
	"github.com/MashiroC/begonia/app/option"
	"github.com/MashiroC/begonia/core"
	"github.com/MashiroC/begonia/dispatch"
	"github.com/MashiroC/begonia/logic"
	"log"
)

// starter.go something
// bootStartByCenter 根据center cluster模式启动
func bootstart(optionMap map[string]interface{}) Center {

	//ctx, cancel := context.WithCancel(context.Background())
	c := &rCenter{
		//ctx:    ctx,
		//cancel: cancel,
	}

	var addr string
	if addrIn, ok := optionMap["managerAddr"]; ok {
		addr = addrIn.(string)
	}

	var dp dispatch.Dispatcher
	dp = dispatch.NewSetByDefaultCluster()
	go dp.Listen(addr)

	c.lg = logic.NewMix(dp)

	core.C = core.NewSubService()

	fmt.Println("  ____                              _        \n |  _ \\                            (_)       \n | |_) |  ___   __ _   ___   _ __   _   __ _ \n |  _ <  / _ \\ / _` | / _ \\ | '_ \\ | | / _` |\n | |_) ||  __/| (_| || (_) || | | || || (_| |\n |____/  \\___| \\__, | \\___/ |_| |_||_| \\__,_|\n                __/ |                        \n               |___/                         ")
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
func New(mode string, optionFunc ...option.WriteFunc) (cli Center) {
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
