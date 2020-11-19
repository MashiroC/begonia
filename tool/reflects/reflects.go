// Package reflects 反射相关工具包
package reflects

import (
	"begonia2/app/coding"
	"begonia2/tool/qconv"
	"reflect"
	"strconv"
)

// ToValue 将一个 map 转化为一个 reflect.Value 的数组
// 该函数会抽取map中key为 “in” + i 的value，组装为数组。i的范围为 1 ~ ∞
func ToValue(m map[string]interface{}, resharp []coding.ReSharpFunc) (res []reflect.Value) {

	res = make([]reflect.Value, 0, 2)

	var i = 1
	for {

		v, ok := m["f"+qconv.I2S(i)]
		if !ok {
			break
		}

		if resharp != nil && len(resharp) >= i && resharp[i-1] != nil {
			v = resharp[i-1](v)
		}

		res = append(res, reflect.ValueOf(v))

		i++
	}

	return
}

// ToInterfaces 将一个 map 转化为一个 interface{}
// 该函数会抽取map中key为 “out” + i 的value，组装为数组。i的范围为 1 ~ ∞
// 根据最终的数组长度来做不同的返回：
//  >1 返回数组
//  =1 返回数组第一个值
//  =0 返回 true
func ToInterfaces(m map[string]interface{}) (res interface{}) {

	tmp := make([]interface{}, 0, 2)

	var i int64 = 1
	for {
		v, ok := m["f"+strconv.FormatInt(i, 10)]
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

// FromValue 将一个 reflect.Value 的数组转化为 map
// 该函数会将数组中每一个值转化为 "out"+i - interface{} i的范围为 1 ~ ∞
func FromValue(values []reflect.Value) (m map[string]interface{}) {

	m = make(map[string]interface{})

	for i, v := range values {
		m["f"+strconv.FormatInt(int64(i+1), 10)] = v.Interface()
	}

	return
}
