package conn

import (
	"github.com/MashiroC/begonia/dispatch"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/dispatch/router"
	"github.com/MashiroC/begonia/tool/qconv"
	"log"
)

// Pool 使用ctrlCode = 0b001 将dispatch的普通连接升级为连接池
func Pool(dp dispatch.Dispatcher) func() (code int, fun router.CtrlHandleFunc) {
	return func() (code int, fun router.CtrlHandleFunc) {
		return frame.CtrlConnCode, func(connID string, typ int, data []byte) {
			addr := qconv.Qb2s(data)
			err := dp.Upgrade(connID, addr)

			if err != nil {
				log.Println(err)
			}
		}
	}
}
