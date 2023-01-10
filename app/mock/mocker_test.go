package mock

import (
	"reflect"
	"testing"
)

func TestMockStoreRegisterByFunc(t *testing.T) {
	type args struct {
		funcName string
		f        interface{}
	}
	tests := []struct {
		name            string
		args            args
		registerSuccess bool
	}{
		{
			name: "test notVariadicISUFFunc",
			args: args{
				funcName: "ISUF",
				f:        testNotVariadicISUFFunc,
			},
			registerSuccess: true,
		},
		{
			name: "test notVariadicMSTFunc",
			args: args{
				funcName: "MST",
				f:        testNotVariadicMSTFunc,
			},
			registerSuccess: true,
		},
		{
			name: "test variadicFunc",
			args: args{
				funcName: "variadicFunc",
				f:        testVariadicFunc,
			},
			registerSuccess: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMockStore()
			m.registerByFunc(tt.args.funcName, tt.args.f)

			if exist := m.IsExist(tt.args.funcName); exist != tt.registerSuccess {
				t.Errorf("RegisterByFunc() want registerSuccess = %v, but got = %v", tt.registerSuccess, exist)
			}
		})
	}
}

func TestMockStoreRegisterByStruct(t *testing.T) {
	type args struct {
		service      interface{}
		registerFunc []string
	}
	tests := []struct {
		name      string
		args      args
		existFunc []string
	}{
		{
			name: "test normal",
			args: args{
				service:      &testReceiver{},
				registerFunc: nil,
			},
			existFunc: []string{"NotVariadicISUFFunc", "NotVariadicMSTFunc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMockStore()
			m.registerByStruct(tt.args.service, tt.args.registerFunc...)

			for _, s := range tt.existFunc {
				if !m.IsExist(s) {
					t.Errorf("RegisterByStruct() want register func = %v but not register", s)
				}
			}
		})
	}
}

func TestMockStoreRegister(t *testing.T) {
	type args struct {
		obj          interface{}
		optionString []string
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{
			name: "test register func",
			args: args{
				obj:          testNotVariadicISUFFunc,
				optionString: []string{"ISUF"},
			},
			wantPanic: false,
		},
		{
			name: "test panic illegal funcName param",
			args: args{
				obj:          testNotVariadicISUFFunc,
				optionString: []string{},
			},
			wantPanic: true,
		},
		{
			name: "test panic illegal input param type",
			args: args{
				obj:          func(c chan bool) {},
				optionString: []string{"illegal"},
			},
			wantPanic: true,
		},
		{
			name: "test panic illegal output param type",
			args: args{
				obj: func() interface{} {
					return nil
				},
				optionString: []string{"illegal"},
			},
			wantPanic: true,
		},
		{
			name: "test register receiver",
			args: args{
				obj:          &testReceiver{},
				optionString: []string{"NotVariadicISUFFunc"},
			},
			wantPanic: false,
		},
		{
			name: "test panic illegal obj",
			args: args{
				obj:          1,
				optionString: nil,
			},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if re := recover(); (re != nil) != tt.wantPanic {
					t.Errorf("Register() want panic = %v but got = %v", tt.wantPanic, re != nil)
				}
			}()

			m := NewMockStore()
			m.Register(tt.args.obj, tt.args.optionString...)
		})
	}
}

func TestMockStoreRegisterRepeated(t *testing.T) {
	mFunc := NewMockStore()
	mFunc.Register(testNotVariadicISUFFunc, "ISUF")

	mRec := NewMockStore()
	mRec.Register(&testReceiver{})

	type args struct {
		ms           *mockStore
		obj          interface{}
		optionString []string
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{
			name: "test panic register func repeated",
			args: args{
				ms:           mFunc,
				obj:          testNotVariadicMSTFunc,
				optionString: []string{"ISUF"},
			},
			wantPanic: true,
		},
		{
			name: "test panic register receiver repeated",
			args: args{
				ms:           mRec,
				obj:          &testReceiver{},
				optionString: nil,
			},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if re := recover(); (re != nil) != tt.wantPanic {
					t.Errorf("Register() want panic = %v but got = %v", tt.wantPanic, re != nil)
				}
			}()

			tt.args.ms.Register(tt.args.obj, tt.args.optionString...)
		})
	}
}

func TestMockStoreExcept(t *testing.T) {
	m := NewMockStore()
	m.Register(&testReceiver{})

	type args struct {
		funcName string
		params   []interface{}
		out      []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test normal",
			args: args{
				funcName: "NotVariadicISUFFunc",
				params:   []interface{}{NewAnyMatch(), NewAnyMatch(), NewAnyMatch(), NewAnyMatch()},
				out:      []interface{}{666},
			},
			wantErr: false,
		},
		{
			name: "test except err",
			args: args{
				funcName: "NotVariadicISUFFunc",
				params:   nil,
				out:      nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := m.Except(tt.args.funcName, tt.args.params, tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf("Except() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMockStoreCall(t *testing.T) {
	m := NewMockStore()
	m.Register(&testReceiver{})
	m.Except("NotVariadicISUFFunc",
		[]interface{}{NewEqualMatch(1), NewEqualMatch("match"), NewAnyMatch(), NewAnyMatch()},
		[]interface{}{666})
	m.Except("NotVariadicISUFFunc",
		[]interface{}{NewAnyMatch(), NewEqualMatch("match"), NewAnyMatch(), NewAnyMatch()},
		[]interface{}{9})
	m.Register(func() (s string, i int) {
		return "", 0
	}, "SI")
	m.Except("SI",
		[]interface{}{},
		[]interface{}{"match", 8})
	m.Register(func() {
		return
	}, "NULL")
	m.Except("NULL",
		[]interface{}{},
		[]interface{}{})

	type args struct {
		funcName string
		params   []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "test call with result1",
			args: args{
				funcName: "NotVariadicISUFFunc",
				params:   []interface{}{1, "match", uint(6), float64(1.12)},
			},
			want:    666,
			wantErr: false,
		},
		{
			name: "test call with result2",
			args: args{
				funcName: "NotVariadicISUFFunc",
				params:   []interface{}{90, "match", uint(6), float64(1.12)},
			},
			want:    9,
			wantErr: false,
		},
		{
			name: "test call with err",
			args: args{
				funcName: "NotVariadicISUFFunc",
				params:   []interface{}{90, "fail", uint(6), float64(1.12)},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test call with []interface",
			args: args{
				funcName: "SI",
				params:   []interface{}{},
			},
			want:    []interface{}{"match", 8},
			wantErr: false,
		},
		{
			name: "test call with no return",
			args: args{
				funcName: "NULL",
				params:   []interface{}{},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := m.Call(tt.args.funcName, tt.args.params...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Call() got = %v, want %v", got, tt.want)
			}
		})
	}
}
