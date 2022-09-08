package mock

import (
	"errors"
	"fmt"
	"github.com/MashiroC/begonia/tool/qarr"
	"reflect"
)

var (
	ErrRepeatRegister = errors.New("repeatedly register mockFunc")
)

// Mocker 定义了mocker仓库的行为
// 一个mock函数应当在mocker仓库中先注册，为其添加规则，再调用以获取返回值
type Mocker interface {
	// Register 注册mock函数
	Register(obj interface{}, optionString ...string)
	// IsExist 判断函数名为funcName的mock函数是否已注册
	IsExist(funcName string) bool
	// Except 为函数名为funcName的mock函数添加规则
	Except(funcName string, params []interface{}, out []interface{}) error
	// Call 调用函数名为funcName的mock函数
	Call(funcName string, params ...interface{}) (res interface{}, err error)
}

type mockStore struct {
	// funcName -> excepts
	mockFuncs map[string]*mockFunc
}

func NewMockStore() *mockStore {
	return &mockStore{
		mockFuncs: make(map[string]*mockFunc),
	}
}

// IsExist 判断函数名为funcName是否已注册
func (m *mockStore) IsExist(funcName string) (exist bool) {
	_, exist = m.mockFuncs[funcName]

	return
}

// Call 调用函数名为funcName的mock函数，params为入参
func (m *mockStore) Call(funcName string, params ...interface{}) (interface{}, error) {
	ec, err := m.mockFuncs[funcName].FindMatch(params...)
	if err != nil {
		return nil, err
	}

	calls, err := ec.Call(params...)
	if err != nil {
		return nil, err
	}

	// 与begonia的reflect工具包里的ToInterfaces函数逻辑对齐
	l := len(calls)
	if l > 1 {
		return calls, err
	} else if l == 1 {
		return calls[0], err
	} else {
		return true, err
	}
}

// Register 注册mock函数。这里有两种使用方式
// 1. obj为结构体，optionString为要注册的函数函数名(为空时不筛选)。注册结构体下所有公开函数为mock函数原型
// 2. obj为函数，optionString为要注册的函数函数名(该参数长度只能为1)。注册obj为mock函数原型
func (m *mockStore) Register(obj interface{}, optionString ...string) {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Struct:
		m.registerByStruct(obj, optionString...)

	case reflect.Func:
		if len(optionString) != 1 {
			panic("function is registerSuccess with zero or more function names")
		}

		m.registerByFunc(optionString[0], obj)

	default:
		panic("illegal obj param")
	}
}

func (m *mockStore) registerByStruct(service interface{}, registerFunc ...string) {
	t := reflect.TypeOf(service)

	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)

		if len(registerFunc) != 0 && !qarr.StringsIn(registerFunc, method.Name) {
			continue
		}

		if _, exist := m.mockFuncs[method.Name]; exist {
			panic(ErrRepeatRegister)
		}

		m.register(method.Name, method.Func.Interface(), service)
	}
}

func (m *mockStore) registerByFunc(funcName string, f interface{}) {
	if _, exist := m.mockFuncs[funcName]; exist {
		panic(ErrRepeatRegister)
	}

	m.register(funcName, f)
}

func (m *mockStore) register(funcName string, funcType interface{}, obj ...interface{}) {
	var o interface{} = nil
	if len(obj) == 1 {
		o = obj[0]
	}

	// 检查是否有不支持的出入参类型
	t := reflect.TypeOf(funcType)
	for i := 0 + len(obj); i < t.NumIn(); i++ {
		if !validKind(t.In(i).Kind()) {
			panic(fmt.Errorf("unsupported remote func input %d param [%v]", i, t.In(i).Kind().String()))
		}
	}
	for i := 0; i < t.NumOut(); i++ {
		if !validKind(t.Out(i).Kind()) {
			panic(fmt.Errorf("unsupported remote func output %d param [%v]", i, t.In(i).Kind().String()))
		}
	}

	m.mockFuncs[funcName] = &mockFunc{
		obj: o,
		fun: reflect.ValueOf(funcType),
		ec:  newExcepts(),
	}
}

// validKind 检查是否是mock支持参数
func validKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.Invalid, reflect.Chan, reflect.Func, reflect.Interface, reflect.UnsafePointer:
		return false

	default:
		return true
	}
}

// Except 为函数名为funcName的匿名函数添加规则
// eg:
// 1. mocker.Except("GetUid", []interface{}{"aaa", 23}, []interface{}{"2333", true})
// 2. mocker.Except("GetUid", []interface{}{"aaa", mock.NewAnyMatch()}, []interface{}{"2333", true})
// 3. mocker.Except("GetUid", []interface{}{
//		mock.NewCustomMatch(func(s string, i int) bool {
//			// logic
//		}),
//	}, []interface{}{"2333", true})
// 4. mocker.Except("GetUid", []interface{}{"aaa", 123},
//		[]interface{}{mock.RetFunc(func(params ...interface{}) (rets []interface{}, err error) {
//			// logic
//	})})
func (m *mockStore) Except(funcName string, params []interface{}, out []interface{}) error {
	mf := m.mockFuncs[funcName]

	ec, err := NewExcept(mf.fun.Interface(), params, out, mf.obj != nil)
	if err != nil {
		return err
	}

	mf.ec = append(mf.ec, ec)
	return nil
}

type mockFunc struct {
	obj interface{}
	fun reflect.Value
	ec  excepts
}

func (m mockFunc) FindMatch(params ...interface{}) (*Except, error) {
	return m.ec.FindMatch(params...)
}
