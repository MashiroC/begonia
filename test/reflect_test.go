package test

import (
	"context"
	"errors"
	"fmt"
	"github.com/MashiroC/begonia"
	"github.com/MashiroC/begonia/app"
	"github.com/MashiroC/begonia/app/center"
	"github.com/MashiroC/begonia/app/client"
	"github.com/MashiroC/begonia/app/option"
	"github.com/stretchr/testify/assert"
	"testing"
)

var service client.Service

type testRegister struct{}

func init() {
	var err error
	addr := ":12306"
	center.New(option.Addr(addr), option.Mode(app.Ast))

	s := begonia.NewServer(option.Addr(addr))
	s.Register("test", &testRegister{})

	client := begonia.NewClient(option.Addr(addr))
	service, err = client.Service("test")
	if err != nil {
		panic(err)
	}
}

func (*testRegister) Null() {
	return
}

func TestNull(t *testing.T) {
	testCases := []struct {
		i      interface{}
		res    interface{}
		hasErr bool
	}{
		{
			i:      nil,
			res:    true,
			hasErr: false,
		},
		{
			i:      "asd",
			res:    nil,
			hasErr: true,
		},
		{
			i:      0,
			res:    nil,
			hasErr: true,
		},
	}
	fun, err := service.FuncSync("Null")
	if err != nil {
		panic(err)
	}

	for idx, c := range testCases {
		t.Run(fmt.Sprintf("case-%d", idx), func(t *testing.T) {
			var res interface{}

			if c.i != nil {
				res, err = fun(c.i)
			} else {
				res, err = fun()
			}

			a := assert.New(t)

			a.Equal(res, c.res)

			if c.hasErr {
				a.NotNil(err)
			} else {
				a.Nil(err)
			}
		})
	}
}

func (*testRegister) OnlyInput(i int) {

}

func TestOnlyInput(t *testing.T) {

	testCases := []struct {
		i      interface{}
		res    interface{}
		hasErr bool
	}{
		{
			i:      1,
			res:    true,
			hasErr: false,
		},
		{
			i:      -1,
			res:    true,
			hasErr: false,
		},
		{
			i:      "asd",
			res:    nil,
			hasErr: true,
		},
		{
			i:      nil,
			res:    nil,
			hasErr: true,
		},
	}
	fun, err := service.FuncSync("OnlyInput")

	for idx, c := range testCases {
		t.Run(fmt.Sprintf("case-%d", idx), func(t *testing.T) {
			if err != nil {
				panic(err)
			}
			res, err := fun(c.i)

			a := assert.New(t)

			a.Equal(res, c.res)

			if c.hasErr {
				a.NotNil(err)
			} else {
				a.Nil(err)
			}
		})
	}

}

func (*testRegister) OnlyOutput() (i int) {
	return 49
}

func TestOnlyOutput(t *testing.T) {
	testCases := []struct {
		i      interface{}
		res    interface{}
		hasErr bool
	}{
		{
			i:      nil,
			res:    49,
			hasErr: false,
		},
		{
			i:      "asd",
			res:    nil,
			hasErr: true,
		},
		{
			i:      0,
			res:    nil,
			hasErr: true,
		},
	}
	fun, err := service.FuncSync("OnlyOutput")
	if err != nil {
		panic(err)
	}

	for idx, c := range testCases {
		t.Run(fmt.Sprintf("case-%d", idx), func(t *testing.T) {
			var res interface{}

			if c.i != nil {
				res, err = fun(c.i)
			} else {
				res, err = fun()
			}

			a := assert.New(t)

			a.Equal(res, c.res)

			if c.hasErr {
				a.NotNil(err)
			} else {
				a.Nil(err)
			}
		})
	}
}

func (*testRegister) BothInAndOut(i int) (j int) {
	return i + 49
}

func TestBothInAndOut(t *testing.T) {
	testCases := []struct {
		i      interface{}
		res    interface{}
		hasErr bool
	}{
		{
			i:      1,
			res:    50,
			hasErr: false,
		},
		{
			i:      nil,
			res:    nil,
			hasErr: true,
		},
		{
			i:      "asd",
			res:    nil,
			hasErr: true,
		},
		{
			i:      -49,
			res:    0,
			hasErr: false,
		},
	}
	fun, err := service.FuncSync("BothInAndOut")
	if err != nil {
		panic(err)
	}

	for idx, c := range testCases {
		t.Run(fmt.Sprintf("case-%d", idx), func(t *testing.T) {
			var res interface{}

			res, err = fun(c.i)

			a := assert.New(t)

			a.Equal(res, c.res)

			if c.hasErr {
				a.NotNil(err)
			} else {
				a.Nil(err)
			}
		})
	}
}

func (*testRegister) OutWithError(i int) (j int, err error) {
	if i < 0 {
		return 0, errors.New("error")
	}
	return i, nil
}

func TestOutWithError(t *testing.T) {
	testCases := []struct {
		i      interface{}
		res    interface{}
		hasErr bool
	}{
		{
			i:      1,
			res:    1,
			hasErr: false,
		},
		{
			i:      nil,
			res:    nil,
			hasErr: true,
		},
		{
			i:      "asd",
			res:    nil,
			hasErr: true,
		},
		{
			i:      -1,
			res:    nil,
			hasErr: true,
		},
	}
	fun, err := service.FuncSync("OutWithError")
	if err != nil {
		panic(err)
	}

	for idx, c := range testCases {
		t.Run(fmt.Sprintf("case-%d", idx), func(t *testing.T) {
			var res interface{}

			res, err = fun(c.i)

			a := assert.New(t)

			a.Equal(res, c.res)

			if c.hasErr {
				a.NotNil(err)
			} else {
				a.Nil(err)
			}
		})
	}
}

func (*testRegister) OutWithContext(ctx context.Context, i int) (s string, err error) {
	v:=ctx.Value("info").(map[string]string)
	if v==nil{
			err=errors.New("ctx err")
	}
	return "true", nil
}

func TestOutWithContext(t *testing.T) {
	testCases := []struct {
		i      interface{}
		res    interface{}
		hasErr bool
	}{
		{
			i:      1,
			res:    "true",
			hasErr: false,
		},
		{
			i:      nil,
			res:    nil,
			hasErr: true,
		},
		{
			i:      "asd",
			res:    nil,
			hasErr: true,
		},
		{
			i:      -49,
			res:    "true",
			hasErr: false,
		},
	}
	fun, err := service.FuncSync("OutWithContext")
	if err != nil {
		panic(err)
	}

	for idx, c := range testCases {
		t.Run(fmt.Sprintf("case-%d", idx), func(t *testing.T) {
			var res interface{}

			res, err = fun(c.i)

			a := assert.New(t)

			a.Equal(res, c.res)

			if c.hasErr {
				a.NotNil(err)
			} else {
				a.Nil(err)
			}
		})
	}
}



func (*testRegister) OutWithErrorAndContext(ctx context.Context, i int) (j int, err error) {
	if i < 0 {
		return 0, errors.New("error")
	}
	return 0, nil
}

func TestOutWithErrorAndContext(t *testing.T) {
	testCases := []struct {
		i      interface{}
		res    interface{}
		hasErr bool
	}{
		{
			i:      1,
			res:    0,
			hasErr: false,
		},
		{
			i:      nil,
			res:    nil,
			hasErr: true,
		},
		{
			i:      "asd",
			res:    nil,
			hasErr: true,
		},
		{
			i:      -49,
			res:    nil,
			hasErr: true,
		},
	}
	fun, err := service.FuncSync("OutWithErrorAndContext")
	if err != nil {
		panic(err)
	}

	for idx, c := range testCases {
		t.Run(fmt.Sprintf("case-%d", idx), func(t *testing.T) {
			var res interface{}

			res, err = fun(c.i)

			a := assert.New(t)

			a.Equal(res, c.res)

			if c.hasErr {
				a.NotNil(err)
			} else {
				a.Nil(err)
			}
		})
	}
}