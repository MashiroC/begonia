package begonia

import (
	"github.com/MashiroC/begonia/app/option"
	"github.com/MashiroC/begonia/app/server"
)

// Server 服务端的接口
type Server = server.Server

// New 初始化
func NewServer(optionFunc ...option.WriteFunc) (s Server) {
	optionMap := defaultServerOption()

	for _, f := range optionFunc {
		f(optionMap)
	}

	in := server.BootStart(optionMap)
	return in
}

func defaultServerOption() map[string]interface{} {
	m := make(map[string]interface{})

	// TODO:加入配置
	m["addr"] = ":12306"

	return m
}
