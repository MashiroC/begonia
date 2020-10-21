package core

import "begonia2/opcode/coding"

var (
	serviceInfoCallCoder coding.Coder
	serviceInfoCoder     coding.Coder

	successCoder coding.Coder
)

func init() {
	var err error

	serviceInfoCallCoder, err = coding.NewAvro(serviceInfoCallRawSchema)
	if err != nil {
		panic(err)
	}

	serviceInfoCoder, err = coding.NewAvro(serviceInfoRawSchema)
	if err != nil {
		panic(err)
	}

	successCoder = &coding.SuccessCoder{}

}