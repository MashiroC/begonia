package core

import "begonia2/opcode/coding"

type ServiceInfo struct {
	Service string           `avro:"service"`
	Funs    []coding.FunInfo `avro:"funs"`
}