package main

import (
	"github.com/MashiroC/begonia/app/center"
	"github.com/MashiroC/begonia/app/option"
	"github.com/MashiroC/begonia/tool/log"
)

func main() {

	//s := os.Args[len(os.Args)-1]

	//if s == "start" {
	addr := ":12306"

	log.InitLogger() // 初始化一个log
	c := center.New(option.Addr(addr))
	//begonialog.CoreLog.Log.Print("test")
	log.Logger.OutCaller()
	log.Logger.Info("begonia start")
	c.Wait()
	//}

}
