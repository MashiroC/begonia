package mock

import (
	"fmt"
	"github.com/MashiroC/begonia/app/client"
	"github.com/MashiroC/begonia/tool/qarr"
	"reflect"
)

type MockClient interface {
	client.Client

	// RegisterMock 注册mock函数
	RegisterMock(serviceName string, obj interface{}, optionString ...string)
	// GetServiceMocker 获取服务mock对象
	GetServiceMocker(serviceName string) Mocker
}

// NewMockClient 获取一个mock客户端，调用远程函数时优先调用已注册的mock函数
func NewMockClient(c ...client.Client) MockClient {
	var cli client.Client = nilClient{}
	if len(c) > 0 {
		cli = c[0]
	}

	return &mClient{
		c:          cli,
		mockStores: make(map[string]Mocker),
	}
}

type mClient struct {
	c client.Client

	// serviceName -> service's mockStore
	mockStores map[string]Mocker
}

// RegisterMock 注册mock函数，其使用方法与一个server注册服务是一样的
// 其会自动在mocker仓库中注册并添加规则，规则的返回值通过调用注册的函数获取
func (mC *mClient) RegisterMock(serviceName string, service interface{}, registerFunc ...string) {
	mocker := mC.GetServiceMocker(serviceName)

	// 判断是不是结构体
	t := reflect.TypeOf(service)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic("service must be struct or struct ptr")
	}

	// 注册mock函数
	mocker.Register(service, registerFunc...)

	var err error

	// 添加except
	receiver := reflect.TypeOf(service)
	for i := 0; i < receiver.NumMethod(); i++ {
		method := receiver.Method(i)

		if len(registerFunc) != 0 && !qarr.StringsIn(registerFunc, method.Name) {
			continue
		}

		// 将method调用封装在 RetFunc 里
		rf := convMethodToRetFunc(service, method)

		// 每个入参都始终匹配。注意：由于不需要入参receiver，故这里长度为numIn-1
		numIn := method.Type.NumIn()
		params := make([]interface{}, 0, numIn-1)
		for j := 0; j < numIn-1; j++ {
			params = append(params, NewAnyMatch())
		}

		err = mocker.Except(method.Name, params, []interface{}{rf})
		if err != nil {
			panic(err)
		}
	}
}

func convMethodToRetFunc(receiver interface{}, method reflect.Method) RetFunc {
	methodType := method.Type

	// 失败时的零值出参
	zeroRets := make([]interface{}, 0)
	for j := 0; j < methodType.NumOut(); j++ {
		zeroRets = append(zeroRets, reflect.Zero(methodType.Out(j)).Interface())
	}

	numIn := methodType.NumIn()

	return func(params ...interface{}) (rets []interface{}, err error) {
		values := make([]reflect.Value, 0, numIn+1)
		values = append(values, reflect.ValueOf(receiver))

		if !methodType.IsVariadic() {
			if len(params) != numIn-1 {
				return zeroRets, fmt.Errorf("needed %d params but input %d params", numIn, len(params))
			}

			for k, param := range params {
				// 检验入参个数与类型是否一致
				if reflect.TypeOf(param).Kind() != methodType.In(k+1).Kind() {
					return zeroRets, fmt.Errorf("needed %d param's type is %v but input %v type param",
						k, reflect.TypeOf(param).Kind(), methodType.In(k+1).Kind())
				}

				values = append(values, reflect.ValueOf(param))
			}
		} else {
			if len(params) < numIn-1 {
				return zeroRets, fmt.Errorf("needed at least %d params but input %d params", numIn-1, len(params))
			}

			for k, param := range params {
				// 检验入参个数与类型是否一致
				if k+1 >= numIn-1 {
					if reflect.TypeOf(param).Kind() != methodType.In(numIn-1).Elem().Kind() {
						return zeroRets, fmt.Errorf("needed %d param's type is %v but input %v type param",
							k, reflect.TypeOf(param).Kind(), methodType.In(k+1).Kind())
					}
				} else {
					if reflect.TypeOf(param).Kind() != methodType.In(k+1).Kind() {
						return zeroRets, fmt.Errorf("needed %d param's type is %v but input %v type param",
							k, reflect.TypeOf(param).Kind(), methodType.In(k+1).Kind())
					}
				}

				values = append(values, reflect.ValueOf(param))
			}
		}

		// 防止调用call时panic的处理
		defer func() {
			if re := recover(); re != nil {
				err = fmt.Errorf("%v", re)
			}
		}()

		calls := method.Func.Call(values)
		rets = make([]interface{}, 0, len(calls))
		for _, call := range calls {
			rets = append(rets, call.Interface())
		}

		return rets, err
	}
}

// GetServiceMocker 获取服务mock对象
func (mC *mClient) GetServiceMocker(serviceName string) Mocker {
	mocker, exist := mC.mockStores[serviceName]
	if !exist {
		mocker = NewMockStore()
		mC.mockStores[serviceName] = mocker
	}

	return mocker
}

func (mC *mClient) Service(name string) (service client.Service, err error) {
	s := &mService{
		name: name,
		mC:   mC,
	}

	// 如果没有注册mock service，找remote service
	if _, exist := mC.mockStores[name]; !exist {
		_, err = mC.c.Service(name)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (mC *mClient) FunSync(serviceName, funName string) (rf client.RemoteFunSync, err error) {
	service, err := mC.Service(serviceName)
	if err != nil {
		return nil, err
	}

	return service.FuncSync(funName)
}

func (mC *mClient) FunAsync(serviceName, funName string) (rf client.RemoteFunAsync, err error) {
	service, err := mC.Service(serviceName)
	if err != nil {
		return nil, err
	}

	return service.FuncAsync(funName)
}

func (mC *mClient) Wait() {
	mC.c.Wait()
}

func (mC *mClient) Close() {
	mC.c.Close()
}

type mService struct {
	name string
	mC   *mClient
}

func (s *mService) FuncSync(name string) (client.RemoteFunSync, error) {
	// 先找有没有对应的mock函数，如果有则优先使用
	if mS, exist := s.mC.mockStores[s.name]; exist && mS.IsExist(name) {
		return func(params ...interface{}) (result interface{}, err error) {
			return mS.Call(name, params...)
		}, nil
	}

	// 没有mock函数，使用远程函数
	funSync, err := s.mC.c.FunSync(s.name, name)
	if err != nil {
		return nil, err
	}

	return func(params ...interface{}) (result interface{}, err error) {
		return funSync(params...)
	}, nil
}

func (s *mService) FuncAsync(name string) (client.RemoteFunAsync, error) {
	// 先找有没有对应的mock函数，如果有则优先使用
	if mS, exist := s.mC.mockStores[s.name]; exist && mS.IsExist(name) {
		return func(callback client.AsyncCallback, params ...interface{}) {
			callback(mS.Call(name, params...))
		}, nil
	}

	// 没有mock函数，使用远程函数
	funAsync, err := s.mC.c.FunAsync(s.name, name)
	if err != nil {
		return nil, err
	}

	return func(callback client.AsyncCallback, params ...interface{}) {
		funAsync(callback, params...)
	}, nil
}

// nilClient 实现 client.Client，但对其的任何远程函数调用都会返回"未找到该服务/函数"
// 调用 Service 方法将返回 nilServer, nil
// 其他的方法调用中不会做任何事
type nilClient struct{}

func (n nilClient) Service(serviceName string) (client.Service, error) {
	return nilServer{}, nil
}

func (n nilClient) FunSync(serviceName, funName string) (client.RemoteFunSync, error) {
	return nilServer{}.FuncSync(funName)
}

func (n nilClient) FunAsync(serviceName, funName string) (client.RemoteFunAsync, error) {
	return nilServer{}.FuncAsync(funName)
}

func (n nilClient) Wait() {
}

func (n nilClient) Close() {
}

// nilServer 实现 client.Service，但对其的任何远程函数调用都会返回"未找到该服务/函数"
type nilServer struct{}

func (n nilServer) FuncSync(name string) (client.RemoteFunSync, error) {
	return nil, fmt.Errorf("RemoteFun [%s] not found", name)
}

func (n nilServer) FuncAsync(name string) (client.RemoteFunAsync, error) {
	return nil, fmt.Errorf("RemoteFun [%s] not found", name)
}
