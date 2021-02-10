package coding

import (
	"github.com/mitchellh/mapstructure"
	"github.com/modern-go/reflect2"
	"reflect"
	"unsafe"
)

func parseReSharpFunc(m reflect.Method) []ReSharpFunc {
	re := make([]ReSharpFunc, m.Type.NumIn()-1)

	for i := 0; i < len(re); i++ {
		re[i] = reSharp(m.Type.In(i + 1))
	}

	return re
}

func reSharp(t reflect.Type) (sharpFunc ReSharpFunc) {
	switch t.Kind() {
	case reflect.Int8:
		sharpFunc = func(in interface{}) interface{} {
			return int8(in.(int))
		}
	case reflect.Int16:
		sharpFunc = func(in interface{}) interface{} {
			return int16(in.(int))
		}
	case reflect.Int32:
		sharpFunc = func(in interface{}) interface{} {
			return int32(in.(int))
		}
	case reflect.Slice:
		if t.Elem().Kind() != reflect.Uint8 {
			sharpFunc = func(in interface{}) interface{} {
				tmp := in.([]interface{})
				s := reflect.MakeSlice(t, 0, 2)
				for i := 0; i < len(tmp); i++ {
					v := reflect.ValueOf(tmp[i])
					s = reflect.Append(s, v)
				}
				//reflect.AppendSlice(s,)
				return s.Interface()
			}
		}
	case reflect.Ptr:
		var resharp ReSharpFunc
		resharp = reSharp(t.Elem())
		sharpFunc = func(in interface{}) interface{} {
			if resharp != nil {
				in = resharp(in)
			}
			v := reflect2.TypeOf(in).PackEFace(unsafe.Pointer(&in))
			return v
		}
	case reflect.Struct:
		sharpFunc = func(in interface{}) interface{} {
			v := reflect.New(t)
			obj := v.Interface()
			err := mapstructure.Decode(in, obj)
			if err != nil {
				panic(err)
			}

			return v.Elem().Interface()
		}
	case reflect.Map:
		var resharp ReSharpFunc
		if t.Elem().Kind() == reflect.Struct {
			resharp = reSharp(t.Elem())
		}
		sharpFunc = func(in interface{}) interface{} {
			out := reflect.MakeMap(t)
			m := in.(map[string]interface{})
			for k, v := range m {
				if resharp != nil {
					v = resharp(v)
				}

				out.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
			}
			return out.Interface()
		}
	default:
		return nil
	}
	return
}
