// Time : 2020/9/26 21:46
// Author : Kieran

// logic
package logic

// entity.go something

const (
	CoreServiceInfo = "ServiceInfo"
)

type Call struct {
	Service string
	Fun     string
	Param   []byte
}

type CallResult struct {
	Result []byte
	Err    string
}
