// Time : 2020/10/20 16:30
// Author : Kieran

// app
package app

import "begonia2/opcode/coding"

// coreservice_entry.go something

type SignInfo struct {
	Service string `avro:"service"`
	Funs []coding.FunInfo `avro:"funs"`
}