package centerlog

import (
	"github.com/MashiroC/begonia/tool/log"
	"github.com/MashiroC/begonia/tool/qconv"
	"io/ioutil"
	"os"
)

type logService struct {
	Name string // 服务名
	Log  *log.Log
}

var path string

func newLogService() *logService {
	l := &logService{Log: log.DefaultNewLogger()}
	// 初始化文件路径
	if len(path) == 0 {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		path = dir + "/store/"
	}

	if err := os.Chdir(path); err != nil {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			panic(err)
		}
	}
	return l
}

// 注册志愿服务中心
func (l *logService) SetLogServicePath(name string) {
	l.Log.FilePath = path + name + ".log"
}

// 获取服务的log信息
func (l *logService) GetLogServiceInfo() ([]byte, error) {
	readAll, err := ioutil.ReadAll(l.Log.GetFile())
	if err != nil {
		return nil, err
	}
	return readAll, nil
}

// 写入文件
//TODO: 日志同步
func (l *logService) Write(msg []byte) (int, error) {
	l.Log.Skip = 6
	return l.Log.Output(qconv.Qb2s(msg))
}
