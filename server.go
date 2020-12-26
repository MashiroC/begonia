package begonia

import (
	"github.com/MashiroC/begonia/app/option"
	"github.com/MashiroC/begonia/app/server"
)

// Service 服务端的接口
type Service interface {
	Register(name string, service interface{})
	Wait()
}

// New 初始化，获得一个service对象，传入一个mode参数，以及一个option的不定参数
func NewServer(optionFunc ...option.WriteFunc) (s Service) {
	optionMap := defaultServiceConfig()

	for _, f := range optionFunc {
		f(optionMap)
	}

	in := server.BootStartByManager(optionMap)
	return in.(Service)
}

func defaultServiceConfig() map[string]interface{} {
	m := make(map[string]interface{})

	// TODO:加入配置
	m["addr"] = ":12306"

	return m
}
