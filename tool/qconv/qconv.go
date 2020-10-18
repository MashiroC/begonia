// Time : 2020/10/6 1:31
// Author : Kieran

// qconv
package qconv

import "unsafe"

// qconv.go something

func Qs2b(str string) []byte {
	return *((*[]byte)((unsafe.Pointer(&str))))
}

func Qb2s(b []byte) string {
	return *((*string)(unsafe.Pointer(&b)))
}
