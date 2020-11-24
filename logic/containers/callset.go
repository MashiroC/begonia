// Time : 2020/8/7 16:41
// Author : MashiroC

// containers
package containers

//import (
//	"github.com/MashiroC/begonia/server/avros"
//	"github.com/linkedin/goavro/v2"
//	"reflect"
//	"strconv"
//	"sync"
//)
//
//// callset.go something
//
//type CallSet struct {
//	l sync.RWMutex
//	m map[string]*RemoteFunc
//}
//
//type RemoteFunc struct {
//	in        interface{}
//	name      string
//	service   string
//	inSchema  string
//	outSchema string
//	inCodec   *goavro.Codec
//	outCodec  *goavro.Codec
//	vFun      reflect.Method
//}
//
//func (f *RemoteFunc) Call(bParams []byte) (b []byte, err error) {
//	native, _, _ := f.inCodec.NativeFromBinary(bParams)
//	datas := native.(map[string]interface{})
//	v := reflect.ValueOf(f.in)
//
//	in := make([]reflect.Value, 0, len(datas)+1)
//
//	in = append(in, v)
//
//	for i := 1; i < len(datas)+1; i++ {
//		v := datas["p"+strconv.FormatInt(int64(i), 10)]
//		//TODO:这里有个类型转换性能问题
//		switch v.(type) {
//		case int32:
//			v = int(v.(int32))
//		}
//		in = append(in, reflect.ValueOf(v))
//	}
//
//	res := f.vFun.Func.Call(in)
//
//	var resErr error
//	if len(res) != 0 {
//		if res[len(res)-1].String() == "<error Value>" {
//			inErr := res[len(res)-1].Interface()
//			if inErr != nil {
//				resErr = inErr.(error)
//			}
//			res = res[:len(res)-1]
//		}
//	}
//
//	outs := make([]interface{}, len(res))
//	for i, v := range res {
//		outs[i] = v.Interface()
//	}
//	avroRes := avros.ReSharp(outs)
//
//	m := avroRes.(map[string]interface{})
//	if resErr != nil {
//		m["err"] = map[string]interface{}{"string": resErr.Error()}
//	} else {
//		m["err"] = nil
//	}
//
//	b, err = f.outCodec.BinaryFromNative(nil, avroRes)
//	return
//}
//
//func NewCallSet() *CallSet {
//	return &CallSet{m: make(map[string]*RemoteFunc)}
//}
//
//func (s *CallSet) Add(service string, param map[string]interface{}, in interface{}, vFun reflect.Method) (err error) {
//	name := param["fun"].(string)
//	inSchema := param["inSchema"].(string)
//	outSchema := param["outSchema"].(string)
//	inCodec, err := goavro.NewCodec(inSchema)
//	if err != nil {
//		return err
//	}
//	outCodec, err := goavro.NewCodec(outSchema)
//	if err != nil {
//		return err
//	}
//	key := service + "." + name
//	s.l.Lock()
//	s.m[key] = &RemoteFunc{
//		in:        in,
//		name:      name,
//		service:   service,
//		inSchema:  inSchema,
//		outSchema: outSchema,
//		inCodec:   inCodec,
//		outCodec:  outCodec,
//		vFun:      vFun,
//	}
//	s.l.Unlock()
//	return
//}
//
//func (s *CallSet) GetFunc(service, fun string) *RemoteFunc {
//	key := service + "." + fun
//	s.l.RLock()
//	defer s.l.RUnlock()
//	return s.m[key]
//}
