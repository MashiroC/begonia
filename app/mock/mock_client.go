package mock

import (
	"github.com/MashiroC/begonia/app/client"
)

type MockClient interface {
	client.Client

	// GetServiceMocker 获取服务mock对象
	GetServiceMocker(serviceName string) Mocker
}

// NewMockClient 获取一个mock客户端
func NewMockClient(c client.Client) MockClient {
	return &mClient{
		c:          c,
		mockStores: make(map[string]Mocker),
	}
}

type mClient struct {
	c client.Client

	// serviceName -> service's mockStore
	mockStores map[string]Mocker
}

// GetServiceMocker 获取服务mock对象
func (mC *mClient) GetServiceMocker(serviceName string) Mocker {
	mocker, exist := mC.mockStores[serviceName]
	if !exist {
		mC.mockStores[serviceName] = NewMockStore()
		mocker = mC.mockStores[serviceName]
	}

	return mocker
}

func (mC *mClient) Service(name string) (service client.Service, err error) {
	s := &mService{
		name: name,
		mC:   mC,
		ser:  nil,
	}

	// 如果没有注册mock service，找remote service
	if _, exist := mC.mockStores[name]; !exist {
		s.ser, err = mC.c.Service(name)
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
	ser  client.Service
}

func (s *mService) FuncSync(name string) (client.RemoteFunSync, error) {
	rf := func(params ...interface{}) (result interface{}, err error) {
		// 先找有没有对应的mock函数，如果有则优先使用
		if mS, exist := s.mC.mockStores[s.name]; exist && mS.IsExist(name) {
			return mS.Call(name, params...)
		}

		// 没有mock函数，使用远程函数
		funSync, err := s.ser.FuncSync(name)
		if err != nil {
			return nil, err
		}

		return funSync(params...)
	}

	return rf, nil
}

func (s *mService) FuncAsync(name string) (client.RemoteFunAsync, error) {
	rf := func(callback client.AsyncCallback, params ...interface{}) {
		// 先找有没有对应的mock函数，如果有则优先使用
		if mS, exist := s.mC.mockStores[s.name]; exist && mS.IsExist(name) {
			callback(mS.Call(name, params...))
		}

		// 没有mock函数，使用远程函数
		funAsync, err := s.ser.FuncAsync(name)
		if err != nil {
			callback(nil, err)
			return
		}

		funAsync(callback, params...)
		return
	}

	return rf, nil
}
