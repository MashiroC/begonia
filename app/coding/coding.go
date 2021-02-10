// Package coding 编码相关的包
package coding

// coding.go something

// Coder 编码器
type Coder interface {
	Encode(data interface{}) ([]byte, error)
	Decode([]byte) (data interface{}, err error)
	DecodeIn([]byte, interface{}) error
}
