package register

//type Register interface {
//	Register(name string, info []coding.FunInfo) (err error)
//	Get(name string) (fs []coding.FunInfo, err error)
//}
//
//// LocalRegister 本地注册器
//type LocalRegister struct {
//
//}
//
//func (r *LocalRegister) Register(name string, info []coding.FunInfo) (err error) {
//	register := core.Call.Register(name, info)
//	_, err = core.C.Invoke("", "", register.Fun, register.Param)
//	return
//}
//
//func (r *LocalRegister) Get(name string) (fs []coding.FunInfo, err error) {
//	si := core.Call.ServiceInfo(name)
//	result, err := core.C.Invoke("", "", si.Fun, si.Param)
//	core.Result.ServiceInfo(result)
//}
//
//// RemoteRegister 远程注册器
//type RemoteRegister struct {
//
//}
