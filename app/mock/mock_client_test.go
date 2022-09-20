package mock

import (
	"github.com/MashiroC/begonia/app/client"
	"reflect"
	"sync"
	"testing"
)

func TestMClient(t *testing.T) {
	mC := NewMockClient()

	type registerArgs struct {
		serviceName  string
		service      interface{}
		registerFunc []string
	}
	registerTests := []struct {
		name      string
		args      registerArgs
		wantPanic bool
	}{
		{
			name: "test register by func",
			args: registerArgs{
				serviceName:  "mock",
				service:      testNotVariadicISUFFunc,
				registerFunc: []string{"ISUF"},
			},
			wantPanic: true,
		},
		{
			name: "test register by receiver",
			args: registerArgs{
				serviceName:  "mock",
				service:      &testReceiver{},
				registerFunc: []string{"NotVariadicISUFFunc", "VariadicFunc"},
			},
			wantPanic: false,
		},
	}
	for _, tt := range registerTests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if re := recover(); (re != nil) != tt.wantPanic {
					t.Errorf("RegisterMock() want panic = %v but got = %v", tt.wantPanic, re != nil)
				}
			}()

			mC.RegisterMock(tt.args.serviceName, tt.args.service, tt.args.registerFunc...)
		})
	}

	var (
		service client.Service
		nvSync  client.RemoteFunSync
		vaAsync client.RemoteFunAsync

		err error
	)

	t.Run("test FunSync not exist", func(t *testing.T) {
		_, err = mC.FunSync("not exist", "not exist")
		if (err != nil) != true {
			t.Errorf("get FuncSync wantErr = %v, got = %v", true, err != nil)
			return
		}
	})

	t.Run("test FunAsync not exist", func(t *testing.T) {
		_, err = mC.FunAsync("not exist", "not exist")
		if (err != nil) != true {
			t.Errorf("get FunAsync wantErr = %v, got = %v", true, err != nil)
			return
		}
	})

	t.Run("test Service", func(t *testing.T) {
		service, err = mC.Service("mock")
		if err != nil {
			t.Errorf("get service fail, err = %v", err)
			return
		}
	})

	t.Run("test FuncSync", func(t *testing.T) {
		nvSync, err = service.FuncSync("NotVariadicISUFFunc")
		if err != nil {
			t.Errorf("get FuncSync ISUF fail, err = %v", err)
			return
		}
	})

	t.Run("test FunAsync", func(t *testing.T) {
		vaAsync, err = mC.FunAsync("mock", "VariadicFunc")
		if err != nil {
			t.Errorf("get FuncASync NotVariadicISUFFunc fail")
			return
		}
	})

	var (
		paramI = 12
	)

	type syncArgs struct {
		sync   client.RemoteFunSync
		params []interface{}
	}
	syncTests := []struct {
		name    string
		args    syncArgs
		want    interface{}
		wantErr bool
	}{
		{
			name: "test sync",
			args: syncArgs{
				sync:   nvSync,
				params: []interface{}{paramI, "isuf", uint(6), float64(3.34)},
			},
			want:    paramI + 1,
			wantErr: false,
		},
		{
			name: "test sync wrong numbers of input params",
			args: syncArgs{
				sync:   nvSync,
				params: []interface{}{paramI},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test sync wrong input param type",
			args: syncArgs{
				sync:   nvSync,
				params: []interface{}{paramI, "isuf", uint(6), "wrong type"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test sync wrong input params numbers",
			args: syncArgs{
				sync:   nvSync,
				params: []interface{}{paramI, "isuf", uint(6)},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range syncTests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := tt.args.sync(tt.args.params...)
			if (err != nil) != tt.wantErr {
				t.Errorf("call FuncSync wantErr = %v, got = %v, err = %v", tt.wantErr, err != nil, err)
			}
			if !reflect.DeepEqual(res, tt.want) {
				t.Errorf("call FuncSync, want = %v, got = %v", tt.want, res)
			}
		})
	}

	var (
		wg sync.WaitGroup
	)

	type asyncArgs struct {
		async  client.RemoteFunAsync
		params []interface{}
	}
	asyncTests := []struct {
		name    string
		args    asyncArgs
		want    interface{}
		wantErr bool
	}{
		{
			name: "test async",
			args: asyncArgs{
				async:  vaAsync,
				params: []interface{}{"async", 1, 2, 3},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "test async missing input params",
			args: asyncArgs{
				async:  vaAsync,
				params: []interface{}{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test async wrong input param type 1",
			args: asyncArgs{
				async:  vaAsync,
				params: []interface{}{uint(2), 3, 5},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range asyncTests {
		t.Run(tt.name, func(t *testing.T) {
			wg.Add(1)

			tt.args.async(func(res interface{}, err error) {
				if (err != nil) != tt.wantErr {
					t.Errorf("call FuncAsync wantErr = %v, got = %v, err = %v", tt.wantErr, err != nil, err)
				}
				if !reflect.DeepEqual(res, tt.want) {
					t.Errorf("call FuncAsync, want = %v, got = %v", tt.want, res)
				}

				wg.Done()
			}, tt.args.params...)
		})
	}

	wg.Wait()
}
