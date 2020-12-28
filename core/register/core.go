package register

import "log"

func (r *CoreRegister) GetToID(serviceName string) (connID string, ok bool) {
	service, ok := r.services.Get(serviceName)
	if !ok {
		return
	}
	connID = service.connID

	return
}

func (r *CoreRegister) HandleConnClose(connID string, err error) {
	log.Printf("conn [%s] closed, unlink server\n", connID)
	r.services.Unlink(connID)
}
