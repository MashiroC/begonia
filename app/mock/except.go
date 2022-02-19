package mock

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrNoExceptSatisfy = errors.New("no except satisfy")
	ErrNotFuncParam    = errors.New("fun param must be func")
	ErrFAMIllegal      = errors.New("funcAllMatch input params isn't same with mockFun")
)

type excepts []*except

func newExcepts() excepts {
	return make([]*except, 0)
}

// FindMatch 找到匹配的except，如果没有返回错误 ErrNoExceptSatisfy
func (e excepts) FindMatch(params ...interface{}) (*except, error) {
	for i := range e {
		if e[i].Matches(params...) {
			return e[i], nil
		}
	}
	return nil, ErrNoExceptSatisfy
}

type except struct {
	fun            reflect.Type
	ignoreReceiver bool
	matches        []Matcher
	getRetsFunc    RetFunc
}

// newExcept 创建一个except
// func -> 函数签名
// params -> 入参期望，个数与函数签名的入参个数一致(ignoreReceiver为true时不需要receiver)
//			 且 一一对应，每个参数可以为定值或 Matcher 接口。
//  		 特别的，可以只传入一个funcAllMatch
// out -> 出参期望，个数与函数签名的出参个数一致。
//        特别的，可以只传入一个 RetFunc
// ignoreReceiver -> 是否忽略receiver
func newExcept(fun interface{}, params []interface{}, out []interface{}, ignoreReceiver bool) (ec except, err error) {
	if reflect.TypeOf(fun).Kind() != reflect.Func {
		return except{}, ErrNotFuncParam
	}
	ec.fun = reflect.TypeOf(fun)

	// 如果忽略receiver
	ig := 0
	if ignoreReceiver {
		recType := ec.fun.In(0)
		if recType.Kind() == reflect.Ptr {
			recType = recType.Elem()
		}
		if recType.Kind() != reflect.Struct {
			fmt.Println(ec.fun.In(0).Kind())
			return except{}, errors.New("fun isn't have receiver")
		}

		ec.matches = append(ec.matches, Any())
		ec.ignoreReceiver = true
		ig = 1
	}

	// 解析获取match
	if len(params) == 1 && isFuncAllMatch(params[0]) {
		// 对于 funcAllMatch 特殊处理
		fam := params[0].(funcAllMatch)
		if fam.f.Type().NumIn() != ec.fun.NumIn()-ig {
			return except{}, ErrFAMIllegal
		}

		famType := fam.f.Type()
		for i := 0; i < famType.NumIn(); i++ {
			if famType.In(i).Kind() != ec.fun.In(i+ig).Kind() {
				return except{}, ErrFAMIllegal
			}
		}

		ec.matches = []Matcher{fam}
	} else {
		if len(params) != ec.fun.NumIn()-ig {
			return except{}, errors.New("the length is inconsistent with the number of input arguments to the function")
		}

		for i := range params {
			switch v := params[i].(type) {
			case Matcher:
				ec.matches = append(ec.matches, v)

			case nil:
				ec.matches = append(ec.matches, Nil())

			default:
				ec.matches = append(ec.matches, Equal(v))
			}
		}
	}

	// 解析获取 RetFunc
	if len(out) == 1 && isRetFunc(out[0]) {
		// 如果是 RetFunc ，直接使用
		ec.getRetsFunc = out[0].(RetFunc)
	} else {
		if len(out) != ec.fun.NumOut() {
			return except{}, errors.New("the length is inconsistent with the number of output arguments to the function")
		}

		o := make([]interface{}, len(out))
		copy(o, out)
		ec.getRetsFunc = func(params ...interface{}) (rets []interface{}, err error) {
			return o, nil
		}
	}

	return ec, nil
}

func isFuncAllMatch(i interface{}) bool {
	_, ok := i.(funcAllMatch)
	return ok
}

func isRetFunc(i interface{}) bool {
	_, ok := i.(RetFunc)
	return ok
}

// Matches 检查是否匹配入参期望
func (e *except) Matches(params ...interface{}) bool {
	// 是 funcAllMatch 则特殊处理
	if e.isFuncAllMatch() {
		return e.matches[0].Match(params)
	}

	// 如果忽略receiver，需要补全receiver参数
	if e.ignoreReceiver {
		params = append([]interface{}{nil}, params...)
	}

	numIn := len(e.matches)

	if !e.fun.IsVariadic() {
		if len(params) != numIn {
			return false
		}

		for i := range params {
			if !e.matches[i].Match(params[i]) {
				return false
			}
		}
		return true
	} else {
		if len(params) < numIn-1 {
			return false
		}

		for i := range params {
			if i < numIn-1 {
				if !e.matches[i].Match(params[i]) {
					return false
				}
				continue
			}
		}

		variaLen := len(params) - (numIn - 1)
		variadic := reflect.MakeSlice(e.fun.In(numIn-1), 0, variaLen)
		for _, param := range params[numIn-1:] {
			variadic = reflect.Append(variadic, reflect.ValueOf(param))
		}

		return e.matches[numIn-1].Match(variadic.Interface())
	}
}

func (e *except) isFuncAllMatch() bool {
	if len(e.matches) != 1 {
		return false
	}

	_, ok := e.matches[0].(funcAllMatch)
	return ok
}

// Call 获取该except的返回值
func (e *except) Call(params ...interface{}) (res []interface{}, err error) {
	return e.getRetsFunc(params...)
}

type RetFunc func(params ...interface{}) (rets []interface{}, err error)

// ZeroRetFunc 根据传入的函数生成一个 RetFunc ，该RetFunc返回值全是零值f
func ZeroRetFunc(fun interface{}) RetFunc {
	t := reflect.TypeOf(fun)

	out := make([]interface{}, 0, t.NumOut())
	for i := 0; i < t.NumOut(); i++ {
		out = append(out, reflect.Zero(t.Out(i)).Interface())
	}

	return func(params ...interface{}) (rets []interface{}, err error) {
		return out, nil
	}
}
