// Package qconv 是快速类型转换库，封装了一些常见的类型转换函数。
package qconv

import (
	"strconv"
	"unsafe"
)

// Qs2b 快速的 string 转 []byte，非内存安全
func Qs2b(str string) []byte {
	return *((*[]byte)((unsafe.Pointer(&str))))
}

// Qb2s 快速的 []byte 转 string，非内存安全
func Qb2s(b []byte) string {
	return *((*string)(unsafe.Pointer(&b)))
}

// I2S 十进制 int 转 string
func I2S(i int) string {
	return strconv.FormatInt(int64(i), 10)
}
