package logger

import (
	centerlog "github.com/MashiroC/begonia/core/begonialog"
	"github.com/MashiroC/begonia/logic"
)

// 日志中心
type LoggerService interface {
	// 服务加入日志中心
	RegisterLogService(name string) error
	// 写入
	Write(name string, msg []byte) (string, error)
	// 查看日志
	Look(name string) ([]byte, error)
}

type localLoggerService struct {
	c *centerlog.CenterLogService
}

func NewLocalLoggerService(c *centerlog.CenterLogService) LoggerService {
	return &localLoggerService{c: c}
}
func (l *localLoggerService) RegisterLogService(name string) error {
	return l.c.RegisterLogService(name)
}
func (l *localLoggerService) Write(name string, msg []byte) (string, error) {
	return l.c.Write(name, msg)
}
func (l *localLoggerService) Look(name string) ([]byte, error) {
	return l.c.GetLog(name)
}

// remoteLoggerService 远程日志服务
type remoteLoggerService struct {
	lg *logic.Client
}

func NewRemoteLoggerService(lg *logic.Client) LoggerService {
	return &remoteLoggerService{
		lg: lg,
	}
}

func (l *remoteLoggerService) RegisterLogService(name string) (err error) {
	var in _CenterLogServiceServiceRegisterLogServiceIn
	in.F1 = name

	b, err := _CenterLogServiceServiceRegisterLogServiceInCoder.Encode(in)
	if err != nil {
		panic(err)
	}
	// RegisterLogService
	res := l.lg.CallSync(&logic.Call{
		Service: "LogService",
		Fun:     "RegisterLogService",
		Param:   b,
	})
	return res.Err
}
func (l *remoteLoggerService) Write(name string, msg []byte) (F1 string, err error) {
	var in _CenterLogServiceServiceWriteIn
	in.F1 = name
	in.F2 = msg
	b, err := _CenterLogServiceServiceWriteInCoder.Encode(in)
	if err != nil {
		panic(err)
	}

	res := l.lg.CallSync(&logic.Call{
		Service: "LogService",
		Fun:     "Write",
		Param:   b,
	})

	var out _CenterLogServiceServiceWriteOut
	err = _CenterLogServiceServiceWriteOutCoder.DecodeIn(res.Result, &out)
	if err != nil {
		panic(err)
	}

	F1 = out.F1

	return
}
func (l *remoteLoggerService) Look(name string) ([]byte, error) {
	var in _CenterLogServiceServiceGetLogIn
	in.F1 = name

	b, err := _CenterLogServiceServiceGetLogInCoder.Encode(in)
	if err != nil {
		panic(err)
	}

	res := l.lg.CallSync(&logic.Call{
		Service: "LogService",
		Fun:     "GetLog",
		Param:   b,
	})

	var out _CenterLogServiceServiceGetLogOut
	err = _CenterLogServiceServiceGetLogOutCoder.DecodeIn(res.Result, &out)
	if err != nil {
		panic(err)
	}

	return out.F1, nil
}
