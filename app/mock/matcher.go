package mock

import (
	"reflect"
)

type Matcher interface {
	Match(i interface{}) bool
}

type FuncMatch func(interface{}) bool

func (f FuncMatch) Match(x interface{}) bool {
	return f(x)
}

// NewFuncMatch 返回一个匹配器，匹配规则自定义
func NewFuncMatch(fun func(interface{}) bool) FuncMatch {
	return fun
}

type CustomMatch struct {
	f reflect.Value
}

func (m *CustomMatch) Match(x interface{}) bool {
	params := x.([]interface{})

	numIn := m.f.Type().NumIn()
	inVals := make([]reflect.Value, 0, numIn)
	if !m.f.Type().IsVariadic() {
		if len(params) != numIn {
			return false
		}
	} else {
		if len(params) < numIn-1 {
			return false
		}
	}

	for i := range params {
		inVals = append(inVals, reflect.ValueOf(params[i]))
	}

	calls := m.f.Call(inVals)
	return calls[0].Bool()
}

// NewCustomMatch 返回一个匹配器，采用传入的所有参数进行自定义匹配
func NewCustomMatch(fun interface{}) *CustomMatch {
	vf := reflect.ValueOf(fun)
	if vf.Type().NumOut() != 1 || vf.Type().Out(0).Kind() != reflect.Bool {
		panic("illegal fun param, the fun out params must and only be `bool`")
	}

	return &CustomMatch{f: vf}
}

type AnyMatch struct{}

func (a *AnyMatch) Match(i interface{}) bool {
	return true
}

// NewAnyMatch 返回一个始终匹配的匹配器
func NewAnyMatch() *AnyMatch {
	return &AnyMatch{}
}

type NilMatch struct{}

func (n *NilMatch) Match(i interface{}) bool {
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

// NewNilMatch 返回一个匹配器，如果接收到的值为 nil，则返回true
func NewNilMatch() *NilMatch {
	return &NilMatch{}
}

type EqualMatch struct {
	Value interface{}
}

func (e *EqualMatch) Match(i interface{}) bool {
	return reflect.DeepEqual(e.Value, i)
}

// NewEqualMatch 返回一个匹配器，如果返回值反射深度相等，则返回true
func NewEqualMatch(i interface{}) *EqualMatch {
	return &EqualMatch{Value: i}
}

type NotMatcher struct {
	M Matcher
}

func (n *NotMatcher) Match(x interface{}) bool {
	return !n.M.Match(x)
}

// NewNotMatch 反转其给定子匹配器的结果
func NewNotMatch(m Matcher) *NotMatcher {
	return &NotMatcher{M: m}
}

type AndMatcher struct {
	Matchers []Matcher
}

func (a *AndMatcher) Match(x interface{}) bool {
	for _, matcher := range a.Matchers {
		if !matcher.Match(x) {
			return false
		}
	}

	return true
}

// NewAndMatch 所有子匹配器都符合时，才认为匹配
// 若子匹配器为空，也认为匹配
func NewAndMatch(matchers ...Matcher) *AndMatcher {
	return &AndMatcher{Matchers: matchers}
}

type OrMatcher struct {
	Matchers []Matcher
}

func (o *OrMatcher) Match(x interface{}) bool {
	for _, matcher := range o.Matchers {
		if matcher.Match(x) {
			return true
		}
	}

	return false
}

// NewOrMatch 所有子匹配器有一个符合时，就认为匹配
// 若子匹配器为空，认为不匹配
func NewOrMatch(matchers ...Matcher) *OrMatcher {
	return &OrMatcher{Matchers: matchers}
}
