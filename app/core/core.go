package core

const (
	ServiceName = "CORE"
)

var C *SubService

type SubService struct {
	services *serviceSet
}

func NewSubService() *SubService {
	return &SubService{services: newServiceSet()}
}

func (s *SubService) Invoke(connID, reqID string, fun string, param []byte) (result []byte, err error) {
	switch fun {
	case "Register":
		var si ServiceInfo
		err = serviceInfoCoder.DecodeIn(param, &si)
		if err != nil {
			panic(err)
		}

		err = s.register(connID,si)
		if err != nil {
			return
		}
		result, err = successCoder.Encode(true)
	case "ServiceInfo":
		var call serviceInfoCall
		var si ServiceInfo
		err = serviceInfoCallCoder.DecodeIn(param,&call)
		if err!=nil{
			return
		}

		si,err = s.serviceInfo(call.Service)
		if err!=nil{
			return
		}
		result, err = serviceInfoCoder.Encode(si)
		if err!=nil{
			return
		}
	default:
		panic("err")
	}

	return
}

func (s *SubService) GetToID(serviceName string) (connID string,ok bool){
	service,ok:=s.services.Get(serviceName)
	if !ok{
		return
	}
	connID=service.connID

	return
}