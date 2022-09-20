package mock

import (
	"errors"
	"reflect"
)

var (
	ErrNoExceptSatisfy    = errors.New("no Except satisfy")
	ErrNotFuncParam       = errors.New("fun param must be func")
	ErrCustomMatchIllegal = errors.New("CustomMatch input params isn't same with mockFun")
)

type excepts []*Except

func newExcepts() excepts {
	return make([]*Except, 0)
}

// FindMatch 找到匹配的except，如果没有返回错误 ErrNoExceptSatisfy
func (e excepts) FindMatch(params ...interface{}) (*Except, error) {
	for i := range e {
		if e[i].Matches(params...) {
			return e[i], nil
		}
	}
	return nil, ErrNoExceptSatisfy
}

type Except struct {
	fun            reflect.Type // 函数原型存根
	ignoreReceiver bool         // 是否忽略接受者。若为true则入参时不需要传入receiver
	matches        []Matcher    // 入参匹配校验。其顺序与 fun 的入参顺序对应，特别的使用 CustomMatch 除外
	getRetsFunc    RetFunc      // 获取返回值函数
}

// NewExcept 创建一个except
// func -> 函数签名
// params -> 入参期望，个数与函数签名的入参个数一致(ignoreReceiver为true时不需要receiver)
//			 且 一一对应，每个参数可以为定值或 Matcher 接口。
//  		 特别的，可以只传入一个funcAllMatch
// out -> 出参期望，个数与函数签名的出参个数一致。
//        特别的，可以只传入一个 RetFunc
// ignoreReceiver -> 是否忽略receiver
func NewExcept(fun interface{}, params []interface{}, out []interface{}, ignoreReceiver bool) (*Except, error) {
	ec := &Except{}

	if reflect.TypeOf(fun).Kind() != reflect.Func {
		return &Except{}, ErrNotFuncParam
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
			return &Except{}, errors.New("fun isn't have receiver")
		}

		ec.matches = append(ec.matches, NewAnyMatch())
		ec.ignoreReceiver = true
		ig = 1
	}

	// 解析获取match
	if len(params) == 1 && isCustomMatch(params[0]) {
		// 对于 CustomMatch 特殊处理
		cm := params[0].(*CustomMatch)
		if cm.f.Type().NumIn() != ec.fun.NumIn()-ig {
			return &Except{}, ErrCustomMatchIllegal
		}

		// 检验入参个数与类型是否一致
		customFunType := cm.f.Type()
		for i := 0; i < customFunType.NumIn(); i++ {
			if customFunType.In(i).Kind() != ec.fun.In(i+ig).Kind() {
				return &Except{}, ErrCustomMatchIllegal
			}
		}

		ec.matches = []Matcher{cm}
	} else {
		if len(params) != ec.fun.NumIn()-ig {
			return &Except{}, errors.New("the length is inconsistent with the number of input arguments to the function")
		}

		for i := range params {
			switch v := params[i].(type) {
			case Matcher:
				ec.matches = append(ec.matches, v)

			case nil:
				ec.matches = append(ec.matches, NewNilMatch())

			default:
				ec.matches = append(ec.matches, NewEqualMatch(v))
			}
		}
	}

	// 解析获取 RetFunc
	if len(out) == 1 && isRetFunc(out[0]) {
		// 如果是 RetFunc ，直接使用
		ec.getRetsFunc = out[0].(RetFunc)
	} else {
		// 检验出参个数是否一致
		if len(out) != ec.fun.NumOut() {
			return &Except{}, errors.New("the length is inconsistent with the number of output arguments to the function")
		}

		o := make([]interface{}, len(out))
		copy(o, out)
		ec.getRetsFunc = func(params ...interface{}) (rets []interface{}, err error) {
			return o, nil
		}
	}

	return ec, nil
}

func isCustomMatch(i interface{}) bool {
	_, ok := i.(*CustomMatch)
	return ok
}

func isRetFunc(i interface{}) bool {
	_, ok := i.(RetFunc)
	return ok
}

// Matches 检查是否匹配入参期望
func (e *Except) Matches(params ...interface{}) bool {
	// 如果是 CustomMatch 则特殊处理
	if e.isCustomMatch() {
		return e.matches[0].Match(params)
	}

	// 如果忽略receiver，需要补全receiver参数
	if e.ignoreReceiver {
		params = append([]interface{}{nil}, params...)
	}

	numIn := len(e.matches)

	if !e.fun.IsVariadic() { // 如果没有可变参数
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

func (e *Except) isCustomMatch() bool {
	if len(e.matches) != 1 {
		return false
	}

	return isCustomMatch(e.matches[0])
}

// Call 获取该except的返回值
func (e *Except) Call(params ...interface{}) (res []interface{}, err error) {
	return e.getRetsFunc(params...)
}

// RetFunc 若一个except的out参数只有一个 RetFunc，则会调用该函数并将其结果作为 Call 的返回值
type RetFunc func(params ...interface{}) (rets []interface{}, err error)
