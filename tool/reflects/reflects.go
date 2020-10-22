// Time : 2020/9/30 21:36
// Author : Kieran

// reflects
package reflects

import (
	"reflect"
	"strconv"
)

// reflects.go something

func ToValue(m map[string]interface{}) (res []reflect.Value) {
	res = make([]reflect.Value, 0, 2)
	var i int64 = 1
	for {
		v, ok := m["in"+strconv.FormatInt(i, 10)]
		if !ok {
			break
		}
		res = append(res, reflect.ValueOf(v))
		i++
	}

	return
}

func ToInterfaces(m map[string]interface{}) (res interface{}) {
	tmp := make([]interface{}, 0, 2)
	var i int64 = 1
	for {
		v, ok := m["out"+strconv.FormatInt(i, 10)]
		if !ok {
			break
		}
		tmp = append(tmp, v)
		i++
	}

	l := len(tmp)

	if l > 1 {
		return tmp
	} else if l == 1 {
		return tmp[0]
	} else {
		return true
	}

	return
}

func FromValue(values []reflect.Value) (m map[string]interface{}) {
	m = make(map[string]interface{})
	for i, v := range values {
		m["out"+strconv.FormatInt(int64(i+1), 10)] = v.Interface()
	}
	return
}
