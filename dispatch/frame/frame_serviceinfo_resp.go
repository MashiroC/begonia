// Time : 2020/8/4 9:10
// Author : MashiroC

// frame
package frame

// frame_serviceinfo_resp.go something

//type Service struct {
//	Name string
//	Funs []Fun
//}
//
//type Fun struct {
//	Name      string
//	InSchema  string
//	OutSchema string
//}
//
//func NewServiceInfoRespWithMap(reqId string,s Datas) Frame{
//	f:=&serviceInfoRespFrame{
//		reqId:   reqId,
//		//service: Service{},
//		m:       s,
//	}
//	return f
//}
//
//func NewServiceInfoResp(reqId, serviceName string, serviceFuns []Fun) Frame {
//	f := &serviceInfoRespFrame{
//		reqId: reqId,
//		service: Service{
//			Name: serviceName,
//			Funs: serviceFuns,
//		},
//	}
//
//	f.m = Datas{
//		"reqId":   f.reqId,
//		"service": f.service.Name,
//	}
//	funs := make([]Datas, 0, len(f.service.Funs))
//	for _, f := range f.service.Funs {
//		funs = append(funs, Datas{
//			"fun":       f.Name,
//			"inSchema":  f.InSchema,
//			"outSchema": f.OutSchema,
//		})
//	}
//	f.m["funs"] = funs
//	return f
//}
//
//type serviceInfoRespFrame struct {
//	reqId   string
//	service Service
//	m       Datas
//}
//
//func (f *serviceInfoRespFrame) Marshal() []byte {
//	//b, err := opcode.Encode(f.Opcode(), f.m)
//	//if err != nil {
//	//	panic(err)
//	//}
//	//return b
//return nil
//}
//
//func (f *serviceInfoRespFrame) UnMarshal() Datas {
//	return f.m
//}
//
//func (f *serviceInfoRespFrame) Opcode() uint8 {
//	return opcode.SignInfoResp
//}
