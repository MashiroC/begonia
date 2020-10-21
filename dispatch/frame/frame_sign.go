// Time : 2020/8/7 16:36
// Author : MashiroC

// frame
package frame

// frame_sign.go something

//type signFrame struct {
//	service string
//	funs    []Datas
//	check   []byte
//	m       Datas
//	v       []byte
//}
//
//func NewSign(service string, funs []Datas) Frame {
//	return &signFrame{
//		service: service,
//		funs:    funs,
//	}
//}
//
//func (f *signFrame) Marshal() []byte {
//	if f.m == nil {
//		f.m = Datas{
//			"service": f.service,
//			"funs":    f.funs,
//		}
//		s := sha256.New()
//		b, _ := json.Marshal(f.m)
//		s.Write(b)
//		f.check = s.Sum(nil)
//		f.m["check"] = f.check
//	}
//
//	if f.v == nil {
//		f.v, _ = opcode.Encode(f.Opcode(), f.m)
//	}
//
//	return f.v
//}
//
//func (f *signFrame) UnMarshal() Datas {
//	if f.m == nil {
//		f.m = Datas{
//			"service": f.service,
//			"funs":    f.funs,
//		}
//		s := sha256.New()
//		b, _ := json.Marshal(f.m)
//		s.Write(b)
//		f.check = s.Sum(nil)
//		f.m["check"] = f.check
//	}
//
//	return f.m
//}
//
//func (f *signFrame) Opcode() uint8 {
//	return opcode.Register
//}
