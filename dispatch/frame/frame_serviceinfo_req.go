// Time : 2020/8/4 9:10
// Author : MashiroC

// frame
package frame

// frame_serviceinfo_req.go something

//func NewServiceInfoReq(service string) Frame {
//	reqId := ids.New()
//	f := &serviceInfoReqFrame{
//		reqId:   reqId,
//		service: service,
//	}
//
//	f.m = Datas{
//		"reqId":   reqId,
//		"service": service,
//	}
//
//	return f
//}
//
//type serviceInfoReqFrame struct {
//	reqId   string
//	service string
//	m       Datas
//}
//
//func (f *serviceInfoReqFrame) Marshal() []byte {
//	//b, err := opcode.Encode(f.Opcode(), f.m)
//	//if err != nil {
//	//	panic(err)
//	//}
//	//return b
//	return nil
//}
//
//func (f *serviceInfoReqFrame) UnMarshal() Datas {
//	return f.m
//}
//
//func (f *serviceInfoReqFrame) Opcode() uint8 {
//	return opcode.SignInfoReq
//}
