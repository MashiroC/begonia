// Time : 2020/9/26 21:17
// Author : Kieran

// option
package option

// option.go something

type OptionFunc func(optionMap map[string]interface{})

func ManagerAddr(addr string) OptionFunc {
	return OptionFunc(func(optionMap map[string]interface{}) {
		optionMap["managerAddr"] = addr
	})
}