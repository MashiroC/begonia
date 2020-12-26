// Package app api层
package app

import "github.com/MashiroC/begonia/internal/coding"

// FunInfo 远程函数的一个封装
type FunInfo struct {
	Name     string       // 远程函数名
	InCoder  coding.Coder // 远程函数入参的编码器
	OutCoder coding.Coder // 远程函数出参的编码器
}
