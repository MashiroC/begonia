// 日志服务
package begonialog

import (
	"github.com/MashiroC/begonia"
	"github.com/MashiroC/begonia/app/option"
	"github.com/MashiroC/begonia/tool/log"
	"io/ioutil"
	"os"
)

//go:generate begonia -s -c -r ../

func StartLogService(addr string) {
	s := begonia.NewServer(option.Addr(addr), option.LogService())
	setCoreLog()
	//CoreLog.Log.Print("test")
	s.Register("LogService", &CoreLog)
	go s.Wait()
}

var CoreLog L

func setCoreLog() {
	CoreLog = L{
		Log: *log.NewLogger(),
	}
}

type L struct {
	Log log.Logger
}

// 获取全部的log信息
func (l *L) GetAllLog() ([]byte, error) {
	file, err := os.OpenFile(l.Log.FilePath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	readAll, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return readAll, nil
}
