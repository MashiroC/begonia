// Code generated by Begonia. DO NOT EDIT.
// versions:
// 	Begonia v1.0.2
// source: register.go
// begonia server1 file

package register

import (
	"context"
	"errors"
	"github.com/MashiroC/begonia/app/coding"
)

var (
	_CoreRegisterFuncList []FunInfo

	_CoreRegisterRegisterInSchema = `
{
			"namespace":"begonia.func.Register",
			"type":"record",
			"name":"In",
			"fields":[
				{"name":"F1","type":{
				"type": "record",
				"name": "Service",
				"fields":[{"name":"Name","type":"string"}
,{"name":"Mode","type":"string"}
,{"name":"Funs","type":{
				"type": "array",
				"items": {
				"type": "record",
				"name": "FunInfo",
				"fields":[{"name":"Name","type":"string"}
,{"name":"InSchema","type":"string"}
,{"name":"OutSchema","type":"string"}

				]
			}
			}}

				]
			},"alias":"si"}

			]
		}`

	_CoreRegisterRegisterOutSchema = `EMPTY_AVRO_SCHEMA`
	_CoreRegisterRegisterInCoder   coding.Coder

	_CoreRegisterRegisterOutCoder coding.Coder

	_CoreRegisterServiceInfoInSchema = `
{
			"namespace":"begonia.func.ServiceInfo",
			"type":"record",
			"name":"In",
			"fields":[
				{"name":"F1","type":"string","alias":"serviceName"}

			]
		}`
	_CoreRegisterServiceInfoOutSchema = `
{
			"namespace":"begonia.func.ServiceInfo",
			"type":"record",
			"name":"Out",
			"fields":[
				{"name":"F1","type":{
				"type": "record",
				"name": "Service",
				"fields":[{"name":"Name","type":"string"}
,{"name":"Mode","type":"string"}
,{"name":"Funs","type":{
				"type": "array",
				"items": {
				"type": "record",
				"name": "FunInfo",
				"fields":[{"name":"Name","type":"string"}
,{"name":"InSchema","type":"string"}
,{"name":"OutSchema","type":"string"}

				]
			}
			}}

				]
			},"alias":"si"}

			]
		}`
	_CoreRegisterServiceInfoInCoder  coding.Coder
	_CoreRegisterServiceInfoOutCoder coding.Coder
)

type _CoreRegisterRegisterIn struct {
	F1 Service
}

type _CoreRegisterRegisterOut struct {
}

type _CoreRegisterServiceInfoIn struct {
	F1 string
}

type _CoreRegisterServiceInfoOut struct {
	F1 Service
}

func init() {
	var err error

	_CoreRegisterRegisterInCoder, err = coding.NewAvro(_CoreRegisterRegisterInSchema)
	if err != nil {
		panic(err)
	}
	_CoreRegisterRegisterOutCoder, err = coding.NewAvro(_CoreRegisterRegisterOutSchema)
	if err != nil {
		panic(err)
	}

	_CoreRegisterServiceInfoInCoder, err = coding.NewAvro(_CoreRegisterServiceInfoInSchema)
	if err != nil {
		panic(err)
	}
	_CoreRegisterServiceInfoOutCoder, err = coding.NewAvro(_CoreRegisterServiceInfoOutSchema)
	if err != nil {
		panic(err)
	}

	_CoreRegisterFuncList = []FunInfo{

		{
			Name:      "Register",
			InSchema:  _CoreRegisterRegisterInSchema,
			OutSchema: _CoreRegisterRegisterOutSchema},

		{
			Name:      "ServiceInfo",
			InSchema:  _CoreRegisterServiceInfoInSchema,
			OutSchema: _CoreRegisterServiceInfoOutSchema},
	}
}

func (r *CoreRegister) Do(ctx context.Context, fun string, param []byte) (result []byte, err error) {
	switch fun {

	case "Register":
		var in _CoreRegisterRegisterIn
		err = _CoreRegisterRegisterInCoder.DecodeIn(param, &in)
		if err != nil {
			panic(err)
		}

		err := r.Register(
			ctx,
			in.F1,
		)
		if err != nil {
			return nil, err
		}
		var out _CoreRegisterRegisterOut

		res, err := _CoreRegisterRegisterOutCoder.Encode(out)
		if err != nil {
			panic(err)
		}
		return res, nil

	case "ServiceInfo":
		var in _CoreRegisterServiceInfoIn
		err = _CoreRegisterServiceInfoInCoder.DecodeIn(param, &in)
		if err != nil {
			panic(err)
		}

		res1, err := r.ServiceInfo(

			in.F1,
		)
		if err != nil {
			return nil, err
		}
		var out _CoreRegisterServiceInfoOut
		out.F1 = res1

		res, err := _CoreRegisterServiceInfoOutCoder.Encode(out)
		if err != nil {
			panic(err)
		}
		return res, nil

	default:
		err = errors.New("rpc call error: fun not exist")
	}
	return
}

func (r *CoreRegister) FuncList() []FunInfo {
	return _CoreRegisterFuncList
}
