// Time : 2020/9/26 19:47
// Author : Kieran

// coding
package coding

// coding.go something

type Coder interface {
	Encode(data interface{}) ([]byte, error)
	Decode([]byte) (data interface{}, err error)
	DecodeIn([]byte, interface{}) error
}

type FunInfo struct {
	Name      string
	Mode      string
	InSchema  string
	OutSchema string
}

func Parse(mode string, in interface{}) (c Coder,fi []FunInfo) {

	return
}
