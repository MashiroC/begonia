package core

const (
	ServiceName = "CORE"
)

type SubService struct {
	services *serviceSet
}

func NewSubService() *SubService {
	return &SubService{services: newServiceSet()}
}

func (s *SubService) Invoke(connID, reqID string, fun string, param []byte) (result []byte, err error) {
	switch fun {
	case "register":
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
		return
	}

	result = []byte{1, 2, 3}
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