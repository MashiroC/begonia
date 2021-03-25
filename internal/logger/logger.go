package logger

import (
	"github.com/MashiroC/begonia/logic"
	"github.com/MashiroC/begonia/tool/log"
)

// 日志中心
type LoggerService interface {
	// 写入
	Save(name string, logMgs *log.Log) error
}

// remoteLoggerService 远程日志服务
// 远程日志服务 -> 日志库 -> 多个日志实例
type remoteLoggerService struct {
	lg    *logic.Client
	store *loggerServiceStore
}

func NewRemoteLoggerService(lg *logic.Client) LoggerService {
	r := &remoteLoggerService{
		lg: lg,
	}
	r.store = newLoggerServiceStore(r)
	return r
}

// 实现信息发送
func (r *remoteLoggerService) Save(name string, logMgs *log.Log) (err error) {
	// 如果危险等级大于了Info,则马上同步
	if logMgs.GetLevel() > log.LevelInfo {
		if err = r.store.Store(name, logMgs); err != nil {
			return
		}
		if err = r.sendMsg(name); err != nil {
			return
		}
	}

	if err = r.store.Store(name, logMgs); err != nil {
		return
	}
	return nil
}

func (r *remoteLoggerService) sendMsg(name string) (err error) {
	var in _CenterLogServiceServiceSaveIn
	in.F1,_ = r.store.Get(name)
	if len(in.F1) == 0 {
		// 如果当前没有日志，那就不需要发送了
		return nil
	}
	b, err := _CenterLogServiceServiceSaveInCoder.Encode(in)
	if err != nil {
		panic(err)
	}

	res := r.lg.CallSync(&logic.Call{
		Service: "LogService",
		Fun:     "Save",
		Param:   b,
	})

	if res.Err != nil {
		err = res.Err
		return
	}

	var out _CenterLogServiceServiceSaveOut
	err = _CenterLogServiceServiceSaveOutCoder.DecodeIn(res.Result, &out)
	if err != nil {
		panic(err)
	}

	return
}
