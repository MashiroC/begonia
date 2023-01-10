package mock

import (
	"reflect"
	"strings"
	"testing"
)

func TestAnyMatch(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test int",
			args: args{i: 1},
		},
		{
			name: "test string",
			args: args{i: "AnyMatch"},
		},
		{
			name: "test uint",
			args: args{i: uint(0)},
		},
		{
			name: "test map",
			args: args{i: map[string]interface{}{}},
		},
		{
			name: "test float64",
			args: args{i: 1.12345232},
		},
		{
			name: "test array",
			args: args{i: [...]interface{}{"AnyMatch is good", 12138, 6.6666, uintptr(123)}},
		},
		{
			name: "test uintptr",
			args: args{i: uintptr(159753)},
		},
		{
			name: "test interface",
			args: args{i: interface{}("aAnyMatch is good")},
		},
		{
			name: "test slice",
			args: args{i: []interface{}{"AnyMatch is good", 12138, 6.6666, uintptr(123)}},
		},
		{
			name: "test struct",
			args: args{i: struct {
				Name string
				Age  int
			}{
				Name: "sky",
				Age:  18,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAnyMatch()
			if got := a.Match(tt.args.i); got != true {
				t.Errorf("Match() = %v, want %v", got, true)
			}
		})
	}
}

func TestCustomMatch(t *testing.T) {
	t.Run("test CustomMatch with illegal output param fun", func(t *testing.T) {
		defer func() {
			if re := recover(); (re != nil) != true {
				t.Errorf("RegisterMock() want panic = %v but got = %v", true, re != nil)
			}
		}()

		NewCustomMatch(func() string {
			return ""
		})
	})

	type people struct {
		name string
		age  int
	}

	notVariadicISUFFunc := func(i int, s string, u uint, f float64) bool {
		return i > 0 && strings.HasPrefix(s, "custom") && u < 1000 && f > 0.0001
	}
	notVariadicMSTFunc := func(m map[string]string, s []int, p people) bool {
		name, exist := m["people"]
		return exist && p.name == name && len(s) > 0 && s[0] >= p.age
	}
	variadicFunc := func(s string, ints ...int) bool {
		return len(ints) > 0
	}

	type fields struct {
		f interface{}
	}
	type args struct {
		x interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "test notVariadicISUFFunc with types int,string,uint,float64",
			fields: fields{f: notVariadicISUFFunc},
			args:   args{x: []interface{}{1, "customMatch", uint(666), 3.1592614}},
			want:   true,
		},
		{
			name:   "test notVariadicISUFFunc with types int,string,uint,float64 fail",
			fields: fields{f: notVariadicISUFFunc},
			args:   args{x: []interface{}{0, "custom match", uint(999), 0.00159}},
			want:   false,
		},
		{
			name:   "test notVariadicMSTFunc with types map,slice,struct",
			fields: fields{f: notVariadicMSTFunc},
			args: args{x: []interface{}{
				map[string]string{
					"people": "sky",
				},
				[]int{21},
				people{
					name: "sky",
					age:  18,
				}}},
			want: true,
		},
		{
			name:   "test notVariadicMSTFunc with types map,slice,struct fail",
			fields: fields{f: notVariadicMSTFunc},
			args: args{x: []interface{}{
				map[string]string{
					"people": "sky",
				},
				[]int{},
				people{
					name: "blue",
					age:  1,
				}}},
			want: false,
		},
		{
			name:   "test notVariadicMSTFunc with types map,slice,struct missing input",
			fields: fields{f: notVariadicMSTFunc},
			args:   args{x: []interface{}{}},
			want:   false,
		},
		{
			name:   "test variadicFunc",
			fields: fields{f: variadicFunc},
			args:   args{x: []interface{}{"CustomMatch", 1, 2, 3}},
			want:   true,
		},
		{
			name:   "test variadicFunc fail",
			fields: fields{f: variadicFunc},
			args:   args{x: []interface{}{"CustomMatch"}},
			want:   false,
		},
		{
			name:   "test variadicFunc missing input",
			fields: fields{f: variadicFunc},
			args:   args{x: []interface{}{}},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewCustomMatch(tt.fields.f)
			if got := m.Match(tt.args.x); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEqualMatch(t *testing.T) {
	type fields struct {
		Value interface{}
	}
	type args struct {
		i interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "test int",
			fields: fields{Value: 1},
			args:   args{i: 1},
			want:   true,
		},
		{
			name:   "test int fail",
			fields: fields{Value: 213},
			args:   args{i: 1},
			want:   false,
		},
		{
			name:   "test string",
			fields: fields{Value: "Equal Match"},
			args:   args{i: "Equal Match"},
			want:   true,
		},
		{
			name:   "test string fail",
			fields: fields{Value: "Any Match"},
			args:   args{i: "Equal Match"},
			want:   false,
		},
		{
			name:   "test uint",
			fields: fields{Value: uint(0)},
			args:   args{i: uint(0)},
			want:   true,
		},
		{
			name:   "test uint fail",
			fields: fields{Value: uint(4)},
			args:   args{i: uint(0)},
			want:   false,
		},
		{
			name: "test map",
			fields: fields{Value: map[string]interface{}{
				"name": "sky",
				"age":  21,
			}},
			args: args{i: map[string]interface{}{
				"name": "sky",
				"age":  21,
			}},
			want: true,
		},
		{
			name: "test map fail",
			fields: fields{Value: map[string]interface{}{
				"name": "sky",
				"age":  38,
			}},
			args: args{i: map[string]interface{}{
				"name": "sky",
				"age":  21,
			}},
			want: false,
		},
		{
			name:   "test float64",
			fields: fields{Value: 1.12345232},
			args:   args{i: 1.12345232},
			want:   true,
		},
		{
			name:   "test float64 fail",
			fields: fields{Value: 23.12345232},
			args:   args{i: 1.12345232},
			want:   false,
		},
		{
			name:   "test array",
			fields: fields{Value: [...]interface{}{"Equal Match is good", 12138, 6.6666, uintptr(123)}},
			args:   args{i: [...]interface{}{"Equal Match is good", 12138, 6.6666, uintptr(123)}},
			want:   true,
		},
		{
			name:   "test array fail",
			fields: fields{Value: [...]interface{}{"Any Match is good", 12138, 6.6666, uintptr(123)}},
			args:   args{i: [...]interface{}{"Equal Match is good", 21315, 8.8, uintptr(666)}},
			want:   false,
		},
		{
			name:   "test uintptr",
			fields: fields{Value: uintptr(159753)},
			args:   args{i: uintptr(159753)},
			want:   true,
		},
		{
			name:   "test uintptr fail",
			fields: fields{Value: uintptr(159753)},
			args:   args{i: uintptr(89)},
			want:   false,
		},
		{
			name:   "test interface",
			fields: fields{Value: interface{}("Equal Match is good")},
			args:   args{i: interface{}("Equal Match is good")},
			want:   true,
		},
		{
			name:   "test interface fail",
			fields: fields{Value: interface{}(85)},
			args:   args{i: interface{}("Equal Match is good")},
			want:   false,
		},
		{
			name:   "test slice",
			fields: fields{Value: []interface{}{"Equal Match is good", 12138, 6.6666, uintptr(123)}},
			args:   args{i: []interface{}{"Equal Match is good", 12138, 6.6666, uintptr(123)}},
			want:   true,
		},
		{
			name:   "test slice fail",
			fields: fields{Value: []interface{}{"Any Match is good", 12138, 6.6666, uintptr(123)}},
			args:   args{i: []interface{}{"Equal Match is good", 21315, 8.8, uintptr(666)}},
			want:   false,
		},
		{
			name: "test struct",
			fields: fields{Value: struct {
				name string
				age  int
			}{
				name: "sky",
				age:  18,
			}},
			args: args{i: struct {
				name string
				age  int
			}{
				name: "sky",
				age:  18,
			}},
			want: true,
		},
		{
			name: "test struct fail",
			fields: fields{Value: struct {
				name string
				age  int
				addr string
			}{
				name: "sky",
				age:  18,
				addr: "your heart",
			}},
			args: args{i: struct {
				name string
				age  int
			}{
				name: "sky",
				age:  18,
			}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEqualMatch(tt.fields.Value)
			if got := e.Match(tt.args.i); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFuncMatch(t *testing.T) {
	stringPrefixMatch := NewFuncMatch(func(i interface{}) bool {
		s, _ := i.(string)
		return strings.HasPrefix(s, "FuncMatch")
	})
	intEnoughMatch := NewFuncMatch(func(i interface{}) bool {
		x, _ := i.(int)
		return x > 100
	})
	sliceLenMatch := NewFuncMatch(func(i interface{}) bool {
		return reflect.ValueOf(i).Len() == 3
	})

	type args struct {
		x interface{}
	}
	tests := []struct {
		name string
		f    FuncMatch
		args args
		want bool
	}{
		{
			name: "test with stringPrefixMatch",
			f:    stringPrefixMatch,
			args: args{x: "FuncMatch is excellent"},
			want: true,
		},
		{
			name: "test with stringPrefixMatch fail",
			f:    stringPrefixMatch,
			args: args{x: "Func Match is bad"},
			want: false,
		},
		{
			name: "test with stringPrefixMatch",
			f:    intEnoughMatch,
			args: args{x: 125},
			want: true,
		},
		{
			name: "test with stringPrefixMatch fail",
			f:    intEnoughMatch,
			args: args{x: 3},
			want: false,
		},
		{
			name: "test with sliceLenMatch",
			f:    sliceLenMatch,
			args: args{x: []int{1, 2, 3}},
			want: true,
		},
		{
			name: "test with sliceLenMatch fail",
			f:    sliceLenMatch,
			args: args{x: []string{"sliceMatch"}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Match(tt.args.x); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNilMatch(t *testing.T) {
	type people struct{}
	var (
		nilMap    map[string]interface{}
		nilSlice  []interface{}
		nilStruct *people
	)

	type args struct {
		i interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test int",
			args: args{i: 1},
			want: false,
		},
		{
			name: "test string",
			args: args{i: "Equal Match"},
			want: false,
		},
		{
			name: "test uint",
			args: args{i: uint(0)},
			want: false,
		},
		{
			name: "test map",
			args: args{i: map[string]interface{}{
				"name": "sky",
				"age":  21,
			}},
			want: false,
		},
		{
			name: "test nil map",
			args: args{i: nilMap},
			want: true,
		},
		{
			name: "test float64",
			args: args{i: 1.12345232},
			want: false,
		},
		{
			name: "test array",
			args: args{i: [...]interface{}{"equal match is good", 12138, 6.6666, uintptr(123)}},
			want: false,
		},
		{
			name: "test uintptr",
			args: args{i: uintptr(159753)},
			want: false,
		},
		{
			name: "test interface",
			args: args{i: interface{}("equal match is good")},
			want: false,
		},
		{
			name: "test slice",
			args: args{i: []interface{}{}},
			want: false,
		},
		{
			name: "test nil slice",
			args: args{i: nilSlice},
			want: true,
		},
		{
			name: "test ptr struct",
			args: args{i: &people{}},
			want: false,
		},
		{
			name: "test nil ptr struct",
			args: args{i: nilStruct},
			want: true,
		},
		{
			name: "test nil",
			args: args{i: nil},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewNilMatch()
			if got := n.Match(tt.args.i); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotMatcher(t *testing.T) {
	type fields struct {
		M Matcher
	}
	type args struct {
		x interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "test AnyMatch",
			fields: fields{M: NewAnyMatch()},
			args:   args{},
			want:   false,
		},
		{
			name:   "test EqualMatch",
			fields: fields{M: NewEqualMatch(12)},
			args:   args{x: 12},
			want:   false,
		},
		{
			name:   "test NilMatch",
			fields: fields{M: NewNilMatch()},
			args: args{x: map[string]string{
				"name": "sky",
			}},
			want: true,
		},
		{
			name: "test CustomMatch",
			fields: fields{M: NewCustomMatch(func() bool {
				return false
			})},
			args: args{[]interface{}{}},
			want: true,
		},
		{
			name: "test FuncMatch",
			fields: fields{M: NewFuncMatch(func(i interface{}) bool {
				return i == nil
			})},
			args: args{x: nil},
			want: false,
		},
		{
			name:   "test NotMatch",
			fields: fields{M: NewNotMatch(NewAnyMatch())},
			args:   args{},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewNotMatch(tt.fields.M)
			if got := n.Match(tt.args.x); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAndMatcher(t *testing.T) {
	successAnyMatcher := NewAnyMatch()
	successCustomMatcher := NewCustomMatch(func(i int) bool {
		return true
	})
	failNotMatcher := NewNotMatch(NewAnyMatch())
	failFuncMatcher := NewFuncMatch(func(i interface{}) bool {
		return false
	})
	neutralNilMatcher := NewNilMatch()
	neutralEqualMatcher := NewEqualMatch("AndMatcher")

	type fields struct {
		Matchers []Matcher
	}
	type args struct {
		x interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "test all match",
			fields: fields{Matchers: []Matcher{successAnyMatcher, successCustomMatcher}},
			args:   args{x: []interface{}{555}},
			want:   true,
		},
		{
			name:   "test some match",
			fields: fields{Matchers: []Matcher{successAnyMatcher, failNotMatcher, failFuncMatcher}},
			args:   args{},
			want:   false,
		},
		{
			name:   "test all not match",
			fields: fields{Matchers: []Matcher{neutralNilMatcher, neutralEqualMatcher}},
			args:   args{x: "fail"},
			want:   false,
		},
		{
			name:   "test nil match",
			fields: fields{Matchers: nil},
			args:   args{x: "success"},
			want:   true,
		},
		{
			name:   "test empty match",
			fields: fields{Matchers: nil},
			args:   args{x: "success"},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAndMatch(tt.fields.Matchers...)
			if got := a.Match(tt.args.x); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrMatcher(t *testing.T) {
	successAnyMatcher := NewAnyMatch()
	successCustomMatcher := NewCustomMatch(func(i int) bool {
		return true
	})
	failNotMatcher := NewNotMatch(NewAnyMatch())
	failFuncMatcher := NewFuncMatch(func(i interface{}) bool {
		return false
	})
	neutralNilMatcher := NewNilMatch()
	neutralEqualMatcher := NewEqualMatch("OrMatcher")

	type fields struct {
		Matchers []Matcher
	}
	type args struct {
		x interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "test all match",
			fields: fields{Matchers: []Matcher{successAnyMatcher, successCustomMatcher}},
			args:   args{x: []interface{}{555}},
			want:   true,
		},
		{
			name:   "test some match",
			fields: fields{Matchers: []Matcher{successAnyMatcher, failNotMatcher, failFuncMatcher}},
			args:   args{},
			want:   true,
		},
		{
			name:   "test all not match",
			fields: fields{Matchers: []Matcher{neutralNilMatcher, neutralEqualMatcher}},
			args:   args{x: "fail"},
			want:   false,
		},
		{
			name:   "test nil match",
			fields: fields{Matchers: nil},
			args:   args{x: "fail"},
			want:   false,
		},
		{
			name:   "test empty match",
			fields: fields{Matchers: nil},
			args:   args{x: "fail"},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := NewOrMatch(tt.fields.Matchers...)
			if got := o.Match(tt.args.x); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
