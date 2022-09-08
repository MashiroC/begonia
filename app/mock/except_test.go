package mock

import (
	"reflect"
	"testing"
)

type testReceiver struct{}

type testStruct struct{}

func (r *testReceiver) NotVariadicISUFFunc(i int, s string, u uint, f float64) (x int) {
	return i + 1
}

func (r *testReceiver) NotVariadicMSTFunc(m map[string]string, s []int, t testStruct) (x string) {
	return
}

func (r *testReceiver) VariadicFunc(s string, ints ...int) (x bool) {
	return len(ints) > 0
}

func testNotVariadicISUFFunc(i int, s string, u uint, f float64) (x int) {
	return
}

func testNotVariadicMSTFunc(m map[string]string, s []int, t testStruct) (x string) {
	return
}

func testVariadicFunc(s string, ints ...int) (x bool) {
	return
}

func TestNewExcept(t *testing.T) {
	receiverMethod0 := reflect.TypeOf(&testReceiver{}).Method(0).Func.Interface()

	type args struct {
		fun            interface{}
		params         []interface{}
		out            []interface{}
		ignoreReceiver bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test not func",
			args: args{
				fun:            1,
				params:         nil,
				out:            nil,
				ignoreReceiver: false,
			},
			wantErr: true,
		},
		{
			name: "test condition - normal",
			args: args{
				fun: testNotVariadicISUFFunc,
				params: []interface{}{NewEqualMatch(1), NewEqualMatch("except"),
					NewEqualMatch(uint(6)), NewEqualMatch(float64(1.12138))},
				out:            []interface{}{666},
				ignoreReceiver: false,
			},
			wantErr: false,
		},
		{
			name: "test condition - missing excepted input parameters",
			args: args{
				fun: testNotVariadicISUFFunc,
				params: []interface{}{NewEqualMatch(1), NewEqualMatch("except"),
					NewEqualMatch(uint(6))},
				out:            []interface{}{666},
				ignoreReceiver: false,
			},
			wantErr: true,
		},
		{
			name: "test condition - overmuch excepted input parameters",
			args: args{
				fun: testNotVariadicISUFFunc,
				params: []interface{}{NewEqualMatch(1), NewEqualMatch("except"),
					NewEqualMatch(uint(6)), NewEqualMatch(float64(1.12138)), NewEqualMatch(1)},
				out:            []interface{}{666},
				ignoreReceiver: false,
			},
			wantErr: true,
		},
		{
			name: "test condition - missing excepted output parameters type",
			args: args{
				fun: testNotVariadicISUFFunc,
				params: []interface{}{NewEqualMatch(1), NewEqualMatch("except"),
					NewEqualMatch(uint(6)), NewEqualMatch(float64(1.12138))},
				out:            []interface{}{},
				ignoreReceiver: false,
			},
			wantErr: true,
		},
		{
			name: "test condition - overmuch excepted output parameters type",
			args: args{
				fun: testNotVariadicISUFFunc,
				params: []interface{}{NewEqualMatch(1), NewEqualMatch("except"),
					NewEqualMatch(uint(6)), NewEqualMatch(float64(1.12138))},
				out:            []interface{}{666, "except"},
				ignoreReceiver: false,
			},
			wantErr: true,
		},
		{
			name: "test condition - ignoreReceiver",
			args: args{
				fun: testNotVariadicISUFFunc,
				params: []interface{}{NewEqualMatch(1), NewEqualMatch("except"),
					NewEqualMatch(uint(6)), NewEqualMatch(float64(1.12138))},
				out:            []interface{}{666},
				ignoreReceiver: true,
			},
			wantErr: true,
		},
		{
			name: "test receiver condition - normal",
			args: args{
				fun: receiverMethod0,
				params: []interface{}{NewEqualMatch(1), NewEqualMatch("except"),
					NewEqualMatch(uint(6)), NewEqualMatch(float64(1.12138))},
				out:            []interface{}{666},
				ignoreReceiver: true,
			},
			wantErr: false,
		},
		{
			name: "test receiver condition - not ignoreReceiver",
			args: args{
				fun: receiverMethod0,
				params: []interface{}{NewAnyMatch(), NewEqualMatch(1), NewEqualMatch("except"),
					NewEqualMatch(uint(6)), NewEqualMatch(float64(1.12138))},
				out:            []interface{}{666},
				ignoreReceiver: false,
			},
			wantErr: false,
		},
		{
			name: "test condition - not matcher params",
			args: args{
				fun:            testNotVariadicMSTFunc,
				params:         []interface{}{nil, []int{1, 2, 3}, testStruct{}},
				out:            []interface{}{"except"},
				ignoreReceiver: false,
			},
			wantErr: false,
		},
		{
			name: "test condition - with CustomMatch",
			args: args{
				fun: testNotVariadicISUFFunc,
				params: []interface{}{NewCustomMatch(func(i int, s string, u uint, f float64) bool {
					return true
				})},
				out:            []interface{}{666},
				ignoreReceiver: false,
			},
			wantErr: false,
		},
		{
			name: "test condition - with illegal input params numbers CustomMatch",
			args: args{
				fun: testNotVariadicISUFFunc,
				params: []interface{}{NewCustomMatch(func(i int) bool {
					return true
				})},
				out:            []interface{}{},
				ignoreReceiver: false,
			},
			wantErr: true,
		},
		{
			name: "test condition - with illegal input param type CustomMatch",
			args: args{
				fun: testNotVariadicISUFFunc,
				params: []interface{}{NewCustomMatch(func(i int, s string, u uint, f map[string]interface{}) bool {
					return true
				})},
				out:            []interface{}{},
				ignoreReceiver: false,
			},
			wantErr: true,
		},
		{
			name: "test condition - with illegal input CustomMatch",
			args: args{
				fun: testNotVariadicISUFFunc,
				params: []interface{}{NewCustomMatch(func() bool {
					return true
				})},
				out:            []interface{}{666},
				ignoreReceiver: false,
			},
			wantErr: true,
		},
		{
			name: "test condition - with RetFunc",
			args: args{
				fun: testNotVariadicISUFFunc,
				params: []interface{}{NewEqualMatch(1), NewEqualMatch("except"),
					NewEqualMatch(uint(6)), NewEqualMatch(float64(1.12138))},
				out: []interface{}{RetFunc(func(params ...interface{}) (rets []interface{}, err error) {
					return []interface{}{12138}, nil
				})},
				ignoreReceiver: false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewExcept(tt.args.fun, tt.args.params, tt.args.out, tt.args.ignoreReceiver)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewExcept() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestExceptMatches(t *testing.T) {
	receiverMethod0 := reflect.TypeOf(&testReceiver{}).Method(0).Func.Interface()

	type fields struct {
		fun            interface{}
		params         []interface{}
		out            []interface{}
		ignoreReceiver bool
	}
	type args struct {
		params []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "test condition - match success",
			fields: fields{
				fun: testNotVariadicISUFFunc,
				params: []interface{}{NewEqualMatch(1), NewEqualMatch("except"),
					NewEqualMatch(uint(6)), NewEqualMatch(float64(1.12138))},
				out:            []interface{}{666},
				ignoreReceiver: false,
			},
			args: args{params: []interface{}{1, "except", uint(6), float64(1.12138)}},
			want: true,
		},
		{
			name: "test condition - match fail",
			fields: fields{
				fun: testNotVariadicISUFFunc,
				params: []interface{}{NewEqualMatch(1), NewEqualMatch("except"),
					NewEqualMatch(uint(6)), NewEqualMatch(float64(1.12138))},
				out:            []interface{}{666},
				ignoreReceiver: false,
			},
			args: args{params: []interface{}{1, "fail", uint(6), float64(1.12138)}},
			want: false,
		},
		{
			name: "test receiver condition - match success",
			fields: fields{
				fun: receiverMethod0,
				params: []interface{}{NewEqualMatch(1), NewEqualMatch("except"),
					NewEqualMatch(uint(6)), NewEqualMatch(float64(1.12138))},
				out:            []interface{}{666},
				ignoreReceiver: true,
			},
			args: args{params: []interface{}{1, "except", uint(6), float64(1.12138)}},
			want: true,
		},
		{
			name: "test receiver condition - match fail",
			fields: fields{
				fun: receiverMethod0,
				params: []interface{}{NewEqualMatch(1), NewEqualMatch("except"),
					NewEqualMatch(uint(6)), NewEqualMatch(float64(1.12138))},
				out:            []interface{}{666},
				ignoreReceiver: true,
			},
			args: args{params: []interface{}{1, "fail", uint(6), float64(1.12138)}},
			want: false,
		},
		{
			name: "test condition - CustomMatch success",
			fields: fields{
				fun: testNotVariadicISUFFunc,
				params: []interface{}{NewCustomMatch(func(i int, s string, u uint, f float64) bool {
					return true
				})},
				out:            []interface{}{666},
				ignoreReceiver: false,
			},
			args: args{params: []interface{}{1, "except", uint(6), float64(1.12138)}},
			want: true,
		},
		{
			name: "test variadicFunc condition - match success",
			fields: fields{
				fun:            testVariadicFunc,
				params:         []interface{}{NewAnyMatch(), NewEqualMatch([]int{1, 2, 3})},
				out:            []interface{}{true},
				ignoreReceiver: false,
			},
			args: args{params: []interface{}{"except", 1, 2, 3}},
			want: true,
		},
		{
			name: "test variadicFunc condition - match fail",
			fields: fields{
				fun:            testVariadicFunc,
				params:         []interface{}{NewEqualMatch("except"), NewEqualMatch([]int{1, 2, 3})},
				out:            []interface{}{true},
				ignoreReceiver: false,
			},
			args: args{params: []interface{}{"fail", 1, 2, 3}},
			want: false,
		},
		{
			name: "test variadicFunc condition - missing input",
			fields: fields{
				fun:            testVariadicFunc,
				params:         []interface{}{NewAnyMatch(), NewEqualMatch([]int{1})},
				out:            []interface{}{true},
				ignoreReceiver: false,
			},
			args: args{params: []interface{}{}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := NewExcept(tt.fields.fun, tt.fields.params, tt.fields.out, tt.fields.ignoreReceiver)
			if err != nil {
				t.Errorf(err.Error())
				return
			}
			if got := e.Matches(tt.args.params...); got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExceptsFindMatch(t *testing.T) {
	except1, _ := NewExcept(testNotVariadicISUFFunc,
		[]interface{}{NewEqualMatch(1), NewEqualMatch("except"), NewEqualMatch(uint(6)), NewEqualMatch(float64(1.12138))},
		[]interface{}{6},
		false)
	except2, _ := NewExcept(testNotVariadicISUFFunc,
		[]interface{}{NewAnyMatch(), NewEqualMatch("except"), NewEqualMatch(uint(6)), NewEqualMatch(float64(1.12138))},
		[]interface{}{10},
		false)
	except3, _ := NewExcept(testNotVariadicISUFFunc,
		[]interface{}{NewAnyMatch(), NewAnyMatch(), NewEqualMatch(uint(6)), NewEqualMatch(float64(1.12138))},
		[]interface{}{205},
		false)

	ecs := excepts{except1, except2, except3}

	type args struct {
		params []interface{}
	}
	tests := []struct {
		name    string
		e       excepts
		args    args
		want    *Except
		wantErr bool
	}{
		{
			name:    "test match except1",
			e:       ecs,
			args:    args{params: []interface{}{1, "except", uint(6), float64(1.12138)}},
			want:    except1,
			wantErr: false,
		},
		{
			name:    "test match except2",
			e:       ecs,
			args:    args{params: []interface{}{432423, "except", uint(6), float64(1.12138)}},
			want:    except2,
			wantErr: false,
		},
		{
			name:    "test match except3",
			e:       ecs,
			args:    args{params: []interface{}{1234598, "fail", uint(6), float64(1.12138)}},
			want:    except3,
			wantErr: false,
		},
		{
			name:    "test no match",
			e:       ecs,
			args:    args{params: []interface{}{1234598, "fail", uint(10), float64(1.12138)}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.FindMatch(tt.args.params...)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindMatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindMatch() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExceptCall(t *testing.T) {
	except1, _ := NewExcept(testNotVariadicISUFFunc,
		[]interface{}{NewEqualMatch(1), NewEqualMatch("except"), NewEqualMatch(uint(6)), NewEqualMatch(float64(1.12138))},
		[]interface{}{6},
		false)

	type fields struct {
		e *Except
	}
	type args struct {
		params []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes []interface{}
		wantErr bool
	}{
		{
			name:    "test except call",
			fields:  fields{e: except1},
			args:    args{params: []interface{}{1, "except", uint(6), float64(1.12138)}},
			wantRes: []interface{}{6},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.fields.e.Call(tt.args.params...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Call() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
