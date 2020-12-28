// Code generated by Begonia. DO NOT EDIT.
// versions:
// 	Begonia v1.0.2
// source: example\server\server.go
// begonia server file

package main

import (
	"context"
	"errors"
	"github.com/MashiroC/begonia/app"
	"github.com/MashiroC/begonia/app/coding"
	cRegister "github.com/MashiroC/begonia/core/register"
)

var (
	_EchoServiceFuncList []cRegister.FunInfo

	_EchoServiceSayHelloInSchema = `
{
			"namespace":"begonia.func.SayHello",
			"type":"record",
			"name":"In",
			"fields":[
				{"name":"F1","type":"string","alias":"name"}

			]
		}`
	_EchoServiceSayHelloOutSchema = `
{
			"namespace":"begonia.func.SayHello",
			"type":"record",
			"name":"Out",
			"fields":[
				{"name":"F1","type":"string"}

			]
		}`
	_EchoServiceSayHelloInCoder  coding.Coder
	_EchoServiceSayHelloOutCoder coding.Coder

	_EchoServiceSayHelloWithContextInSchema = `
{
			"namespace":"begonia.func.SayHelloWithContext",
			"type":"record",
			"name":"In",
			"fields":[
				{"name":"F1","type":"string","alias":"name"}

			]
		}`
	_EchoServiceSayHelloWithContextOutSchema = `
{
			"namespace":"begonia.func.SayHelloWithContext",
			"type":"record",
			"name":"Out",
			"fields":[
				{"name":"F1","type":"string"}

			]
		}`
	_EchoServiceSayHelloWithContextInCoder  coding.Coder
	_EchoServiceSayHelloWithContextOutCoder coding.Coder

	_EchoServiceAddInSchema = `
{
			"namespace":"begonia.func.Add",
			"type":"record",
			"name":"In",
			"fields":[
				{"name":"F1","type":"int","alias":"i1"}
,{"name":"F2","type":"int","alias":"i2"}

			]
		}`
	_EchoServiceAddOutSchema = `
{
			"namespace":"begonia.func.Add",
			"type":"record",
			"name":"Out",
			"fields":[
				{"name":"F1","type":"int","alias":"res"}

			]
		}`
	_EchoServiceAddInCoder  coding.Coder
	_EchoServiceAddOutCoder coding.Coder

	_EchoServiceModInSchema = `
{
			"namespace":"begonia.func.Mod",
			"type":"record",
			"name":"In",
			"fields":[
				{"name":"F1","type":"int","alias":"i1"}
,{"name":"F2","type":"int","alias":"i2"}

			]
		}`
	_EchoServiceModOutSchema = `
{
			"namespace":"begonia.func.Mod",
			"type":"record",
			"name":"Out",
			"fields":[
				{"name":"F1","type":"int","alias":"res1"}
,{"name":"F2","type":"int","alias":"res2"}

			]
		}`
	_EchoServiceModInCoder  coding.Coder
	_EchoServiceModOutCoder coding.Coder

	_EchoServiceNULLInSchema = `
{
			"namespace":"begonia.func.NULL",
			"type":"record",
			"name":"In",
			"fields":[
				
			]
		}`
	_EchoServiceNULLOutSchema = `
{
			"namespace":"begonia.func.NULL",
			"type":"record",
			"name":"Out",
			"fields":[
				
			]
		}`
	_EchoServiceNULLInCoder  coding.Coder
	_EchoServiceNULLOutCoder coding.Coder
)

type _EchoServiceSayHelloIn struct {
	F1 string
}

type _EchoServiceSayHelloOut struct {
	F1 string
}

type _EchoServiceSayHelloWithContextIn struct {
	F1 string
}

type _EchoServiceSayHelloWithContextOut struct {
	F1 string
}

type _EchoServiceAddIn struct {
	F1 int
	F2 int
}

type _EchoServiceAddOut struct {
	F1 int
}

type _EchoServiceModIn struct {
	F1 int
	F2 int
}

type _EchoServiceModOut struct {
	F1 int
	F2 int
}

type _EchoServiceNULLIn struct {
}

type _EchoServiceNULLOut struct {
}

func init() {
	app.ServiceAppMode = app.Ast

	var err error

	_EchoServiceSayHelloInCoder, err = coding.NewAvro(_EchoServiceSayHelloInSchema)
	if err != nil {
		panic(err)
	}
	_EchoServiceSayHelloOutCoder, err = coding.NewAvro(_EchoServiceSayHelloOutSchema)
	if err != nil {
		panic(err)
	}

	_EchoServiceSayHelloWithContextInCoder, err = coding.NewAvro(_EchoServiceSayHelloWithContextInSchema)
	if err != nil {
		panic(err)
	}
	_EchoServiceSayHelloWithContextOutCoder, err = coding.NewAvro(_EchoServiceSayHelloWithContextOutSchema)
	if err != nil {
		panic(err)
	}

	_EchoServiceAddInCoder, err = coding.NewAvro(_EchoServiceAddInSchema)
	if err != nil {
		panic(err)
	}
	_EchoServiceAddOutCoder, err = coding.NewAvro(_EchoServiceAddOutSchema)
	if err != nil {
		panic(err)
	}

	_EchoServiceModInCoder, err = coding.NewAvro(_EchoServiceModInSchema)
	if err != nil {
		panic(err)
	}
	_EchoServiceModOutCoder, err = coding.NewAvro(_EchoServiceModOutSchema)
	if err != nil {
		panic(err)
	}

	_EchoServiceNULLInCoder, err = coding.NewAvro(_EchoServiceNULLInSchema)
	if err != nil {
		panic(err)
	}
	_EchoServiceNULLOutCoder, err = coding.NewAvro(_EchoServiceNULLOutSchema)
	if err != nil {
		panic(err)
	}

	_EchoServiceFuncList = []cRegister.FunInfo{

		{
			Name:      "SayHello",
			InSchema:  _EchoServiceSayHelloInSchema,
			OutSchema: _EchoServiceSayHelloOutSchema},

		{
			Name:      "SayHelloWithContext",
			InSchema:  _EchoServiceSayHelloWithContextInSchema,
			OutSchema: _EchoServiceSayHelloWithContextOutSchema},

		{
			Name:      "Add",
			InSchema:  _EchoServiceAddInSchema,
			OutSchema: _EchoServiceAddOutSchema},

		{
			Name:      "Mod",
			InSchema:  _EchoServiceModInSchema,
			OutSchema: _EchoServiceModOutSchema},

		{
			Name:      "NULL",
			InSchema:  _EchoServiceNULLInSchema,
			OutSchema: _EchoServiceNULLOutSchema},
	}
}

func (d *EchoService) Do(ctx context.Context, fun string, param []byte) (result []byte, err error) {
	switch fun {

	case "SayHello":
		var in _EchoServiceSayHelloIn
		err = _EchoServiceSayHelloInCoder.DecodeIn(param, &in)
		if err != nil {
			panic(err)
		}

		res1 := d.SayHello(

			in.F1,
		)
		if err != nil {
			return nil, err
		}
		var out _EchoServiceSayHelloOut
		out.F1 = res1

		res, err := _EchoServiceSayHelloOutCoder.Encode(out)
		if err != nil {
			panic(err)
		}
		return res, nil

	case "SayHelloWithContext":
		var in _EchoServiceSayHelloWithContextIn
		err = _EchoServiceSayHelloWithContextInCoder.DecodeIn(param, &in)
		if err != nil {
			panic(err)
		}

		res1 := d.SayHelloWithContext(
			ctx,
			in.F1,
		)
		if err != nil {
			return nil, err
		}
		var out _EchoServiceSayHelloWithContextOut
		out.F1 = res1

		res, err := _EchoServiceSayHelloWithContextOutCoder.Encode(out)
		if err != nil {
			panic(err)
		}
		return res, nil

	case "Add":
		var in _EchoServiceAddIn
		err = _EchoServiceAddInCoder.DecodeIn(param, &in)
		if err != nil {
			panic(err)
		}

		res1, err := d.Add(

			in.F1, in.F2,
		)
		if err != nil {
			return nil, err
		}
		var out _EchoServiceAddOut
		out.F1 = res1

		res, err := _EchoServiceAddOutCoder.Encode(out)
		if err != nil {
			panic(err)
		}
		return res, nil

	case "Mod":
		var in _EchoServiceModIn
		err = _EchoServiceModInCoder.DecodeIn(param, &in)
		if err != nil {
			panic(err)
		}

		res1, res2 := d.Mod(

			in.F1, in.F2,
		)
		if err != nil {
			return nil, err
		}
		var out _EchoServiceModOut
		out.F1 = res1
		out.F2 = res2

		res, err := _EchoServiceModOutCoder.Encode(out)
		if err != nil {
			panic(err)
		}
		return res, nil

	case "NULL":
		var in _EchoServiceNULLIn
		err = _EchoServiceNULLInCoder.DecodeIn(param, &in)
		if err != nil {
			panic(err)
		}

		d.NULL()
		return []byte{1}, nil

	default:
		err = errors.New("rpc call error: fun not exist")
	}
	return
}

func (d *EchoService) FuncList() []cRegister.FunInfo {
	return _EchoServiceFuncList
}
