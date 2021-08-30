package retry

import (
	"fmt"
	"github.com/MashiroC/begonia/config"
	"log"
	"time"
)

func Always(actionName string, fun func() bool, intervalSeconds int) {
	_ = Do(actionName, fun, -1, intervalSeconds)
}

func Do(actionName string, fun func() bool, times, intervalSeconds int) (err error) {
	var ok bool

	if times <= 0 {
		count := 1

		for !ok {
			if count > 1 {
				time.Sleep(time.Duration(intervalSeconds) * time.Second)
				log.Printf("retry for action: %s, times: %d, limit:always\n", actionName, count)
			}
			ok = fun()
			count++
		}

	} else {

		for i := 0; i < times && !ok; i++ {
			if i != 0 {
				log.Printf("retry for action: %s, times: %d, limit: %d\n", actionName, i+1, times)
				time.Sleep(time.Duration(config.C.Dispatch.ConnectionIntervalSecond) * time.Second)
			}
			ok = fun()
			if ok {
				break
			}
		}

		if !ok {
			err = fmt.Errorf("retry filed for action:%s", actionName)
		}

	}

	return
}
