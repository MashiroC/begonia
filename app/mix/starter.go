package mix

import (
	"fmt"
	"github.com/MashiroC/begonia/app"
	"github.com/MashiroC/begonia/app/client"
	"github.com/MashiroC/begonia/app/server"
	"log"
)

// BootStartByCenter 根据center cluster模式启动
func BootStart(optionMap map[string]interface{}) *MixNode {

	fmt.Println("  ____                              _        \n |  _ \\                            (_)       \n | |_) |  ___   __ _   ___   _ __   _   __ _ \n |  _ <  / _ \\ / _` | / _ \\ | '_ \\ | | / _` |\n | |_) ||  __/| (_| || (_) || | | || || (_| |\n |____/  \\___| \\__, | \\___/ |_| |_||_| \\__,_|\n                __/ |                        \n               |___/                         ")

	log.Printf("begonia client start with [%s] mode\n", app.ServiceAppMode)

	// TODO:给dispatch初始化

	s := server.BootStart(optionMap)
	lg := server.GetLogic(s)

	c := client.BootStartWithLogic(optionMap, lg.Client)
	m := &MixNode{
		cli:    c,
		server: s,
	}

	return m
}
