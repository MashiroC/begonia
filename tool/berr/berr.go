// Package berr 自定义的error系统
package berr

import (
	"fmt"
)

// BErr error的实现结构体
type BErr struct {
	Package string // Package 包名，也可以是系统的名 需要遵循 "aa" 或 "aa.bb" 的命名方式
	Typ     string // Typ 错误类型，一般是发生错误的函数名，或者当前在做的哪一段逻辑出错
	Msg     string // Msg 错误信息
}

func (err BErr) Error() string {
	// fmt.Sprintf("%s %s error: %s")
	//          [  pkg ] [typ]       [           msg            ]
	// example: dispatch link error: you are link in a black hole

	return err.Package + " " + err.Typ + " error: " + err.Msg
}

// New 创建一个新的error，如果传入的msg为空，则会创建一个nil
func New(pkg, typ, msg string) error {
	if msg == "" {
		return nil
	}
	return BErr{
		Package: pkg,
		Typ:     typ,
		Msg:     msg,
	}
}

// NewF 是 New 的 Formatter 版本
func NewF(pkg, typ, msgF string, param ...interface{}) error {
	return BErr{
		Package: pkg,
		Typ:     typ,
		Msg:     fmt.Sprintf(msgF, param...),
	}
}

// Warp 包装一个error，把包装的error作为当前error的msg
func Warp(pkg, typ string, err error) error {
	return BErr{
		Package: pkg,
		Typ:     typ,
		Msg:     err.Error(),
	}
}
