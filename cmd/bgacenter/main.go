package main

import (
	"github.com/MashiroC/begonia/app/begonialog"
	"github.com/MashiroC/begonia/app/center"
	"github.com/MashiroC/begonia/app/option"
)

func main() {

	//s := os.Args[len(os.Args)-1]

	//if s == "start" {
	addr := ":12306"

	c := center.New(option.Addr(addr))
	begonialog.CoreLog.Log.Print("test")
	c.Wait()
	//}

}
