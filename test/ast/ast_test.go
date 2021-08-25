package ast

import (
	"fmt"
	"github.com/MashiroC/begonia"
	"github.com/MashiroC/begonia/app"
	"github.com/MashiroC/begonia/app/center"
	"github.com/MashiroC/begonia/app/option"
	"github.com/MashiroC/begonia/test/ast/call"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	addr := ":12306"
	center.New(option.Addr(addr), option.Mode(app.Ast))

	s := begonia.NewServer(option.Addr(addr))
	s.Register("test", &testRegister{})

	call.Init()
}

func TestNull(t *testing.T) {
	testCases := []struct {
		hasErr bool
	}{
		{
			hasErr: false,
		},
	}

	for idx, c := range testCases {
		t.Run(fmt.Sprintf("case-%d", idx), func(t *testing.T) {
			err := call.Null()

			a := assert.New(t)

			if c.hasErr {
				a.NotNil(err)
			} else {
				a.Nil(err)
			}
		})
	}
}

func TestOnlyInput(t *testing.T) {

	testCases := []struct {
		i      int
		hasErr bool
	}{
		{
			i:      1,
			hasErr: false,
		},
		{
			i:      -1,
			hasErr: false,
		},
	}

	for idx, c := range testCases {
		t.Run(fmt.Sprintf("case-%d", idx), func(t *testing.T) {
			err := call.OnlyInput(c.i)

			a := assert.New(t)

			if c.hasErr {
				a.NotNil(err)
			} else {
				a.Nil(err)
			}
		})
	}

}

func TestOnlyOutput(t *testing.T) {
	testCases := []struct {
		res    interface{}
		hasErr bool
	}{
		{
			res:    49,
			hasErr: false,
		},
	}

	for idx, c := range testCases {
		t.Run(fmt.Sprintf("case-%d", idx), func(t *testing.T) {
			res, err := call.OnlyOutput()

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

func TestBothInAndOut(t *testing.T) {
	testCases := []struct {
		i      int
		res    interface{}
		hasErr bool
	}{
		{
			i:      1,
			res:    50,
			hasErr: false,
		},
		{
			i:      -49,
			res:    0,
			hasErr: false,
		},
	}

	for idx, c := range testCases {
		t.Run(fmt.Sprintf("case-%d", idx), func(t *testing.T) {
			res, err := call.BothInAndOut(c.i)

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

func TestOutWithError(t *testing.T) {
	testCases := []struct {
		i      int
		res    int
		hasErr bool
	}{
		{
			i:      1,
			res:    1,
			hasErr: false,
		},
		{
			i:      -1,
			res:    0,
			hasErr: true,
		},
	}

	for idx, c := range testCases {
		t.Run(fmt.Sprintf("case-%d", idx), func(t *testing.T) {
			var res interface{}

			res, err := call.OutWithError(c.i)

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

func TestOutWithContext(t *testing.T) {
	testCases := []struct {
		i      int
		res    string
		hasErr bool
	}{
		{
			i:      1,
			res:    "true",
			hasErr: false,
		},
		{
			i:      -49,
			res:    "true",
			hasErr: false,
		},
	}

	for idx, c := range testCases {
		t.Run(fmt.Sprintf("case-%d", idx), func(t *testing.T) {

			res, err := call.OutWithContext(c.i)

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

func TestOutWithErrorAndContext(t *testing.T) {
	testCases := []struct {
		i      int
		res    int
		hasErr bool
	}{
		{
			i:      1,
			res:    0,
			hasErr: false,
		},
		{
			i:      -49,
			res:    0,
			hasErr: true,
		},
	}

	for idx, c := range testCases {
		t.Run(fmt.Sprintf("case-%d", idx), func(t *testing.T) {
			res, err := call.OutWithErrorAndContext(c.i)

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
