package core

// coder.go 保存着一些硬编码的函数的coder

import "github.com/MashiroC/begonia/internal/coding"

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
