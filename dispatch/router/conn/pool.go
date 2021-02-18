package conn

import (
	"github.com/MashiroC/begonia/dispatch"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/tool/qconv"
	"log"
)

type HandleFunc = func(connID string, data []byte)

// Pool 使用ctrlCode = 0b001 将dispatch的普通连接升级为连接池
func Pool(dp dispatch.Dispatcher) func() (code int, fun HandleFunc) {
	return func() (code int, fun HandleFunc) {
		return frame.CtrlConnCode, func(connID string, data []byte) {
			addr := qconv.Qb2s(data)
			err := dp.Upgrade(connID, addr)

			if err != nil {
				log.Println(err)
			}
		}
	}
}