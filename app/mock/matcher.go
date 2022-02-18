package mock

import (
	"reflect"
)

type Matcher interface {
	Match(i interface{}) bool
}

type funcMatch struct {
	f reflect.Value
}

func (f funcMatch) Match(x interface{}) bool {
	calls := f.f.Call([]reflect.Value{reflect.ValueOf(x)})
	return calls[0].Bool()
}

type funcAllMatch struct {
	f reflect.Value
}

// Func 返回一个匹配器，匹配规则自定义
func Func(fun func(interface{}) bool) Matcher {
	return funcMatch{f: reflect.ValueOf(fun)}
}

func (m funcAllMatch) Match(x interface{}) bool {
	params := x.([]interface{})

	numIn := m.f.Type().NumIn()
	inVals := make([]reflect.Value, 0, numIn)
	if !m.f.Type().IsVariadic() {
		if len(params) != numIn {
			return false
		}

		for i := range params {
			inVals = append(inVals, reflect.ValueOf(params[i]))
		}
	} else {
		if len(params) < numIn-1 {
			return false
		}

		for i := range params {
			if i < numIn-1 {
				inVals = append(inVals, reflect.ValueOf(params[i]))
			}
		}

		variaLen := len(params) - (numIn - 1)
		variadic := reflect.MakeSlice(m.f.Type().In(numIn-1), 0, variaLen)
		for _, param := range params[numIn-1:] {
			reflect.Append(variadic, reflect.ValueOf(param))
		}
		inVals = append(inVals, variadic)
	}

	calls := m.f.Call(inVals)
	return calls[0].Bool()
}

// FuncAll 返回一个匹配器，采用传入的所有参数进行自定义匹配
func FuncAll(fun interface{}) Matcher {
	vf := reflect.ValueOf(fun)
	if vf.Type().NumOut() != 1 || vf.Type().Out(0).Kind() != reflect.Bool {
		panic("illegal fun param, the fun out params must and only be `bool`")
	}

	return funcAllMatch{f: vf}
}

type anyMatch struct{}

func (a anyMatch) Match(i interface{}) bool {
	return true
}

// Any 返回一个始终匹配的匹配器
func Any() Matcher {
	return anyMatch{}
}

type nilMatch struct{}

func (n nilMatch) Match(i interface{}) bool {
	if i == nil {
		return true
	}

	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	}

	return false
}

// Nil 返回一个匹配器，如果接收到的值为 nil，则返回true
func Nil() Matcher {
	return nilMatch{}
}

type equalMatch struct {
	x interface{}
}

func (e equalMatch) Match(i interface{}) bool {
	return reflect.DeepEqual(e.x, i)
}

// Equal 返回一个匹配器，如果返回值反射深度相等，则返回true
func Equal(i interface{}) Matcher {
	return equalMatch{x: i}
}

type notMatcher struct {
	m Matcher
}

func (n notMatcher) Match(x interface{}) bool {
	return !n.m.Match(x)
}

// Not 反转其给定子匹配器的结果
func Not(i interface{}) Matcher {
	if m, ok := i.(Matcher); ok {
		return notMatcher{m}
	}
	return notMatcher{Equal(i)}
}
