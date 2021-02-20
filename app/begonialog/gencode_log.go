package begonialog

import (
	"context"
	"errors"
	cRegister "github.com/MashiroC/begonia/core/register"
	"github.com/MashiroC/begonia/app"
	"github.com/MashiroC/begonia/app/coding"
)


var (
	_LFuncList []cRegister.FunInfo

	
	_LGetAllLogInSchema  = `
{
			"namespace":"begonia.func.GetAllLog",
			"type":"record",
			"name":"In",
			"fields":[
				
			]
		}`
	_LGetAllLogOutSchema = `
{
			"namespace":"begonia.func.GetAllLog",
			"type":"record",
			"name":"Out",
			"fields":[
				{"name":"F1","type":"bytes"}

			]
		}`
	_LGetAllLogInCoder   coding.Coder
	_LGetAllLogOutCoder  coding.Coder

)




 	
type _LGetAllLogIn struct {
	}

type _LGetAllLogOut struct {
	
			F1 []byte  
		 
}


func init() {
	app.ServiceAppMode = app.Ast

	var err error
 	
	_LGetAllLogInCoder, err = coding.NewAvro(_LGetAllLogInSchema)
	if err != nil {
		panic(err)
	}
	_LGetAllLogOutCoder, err = coding.NewAvro(_LGetAllLogOutSchema)
	if err != nil {
		panic(err)
	}


	_LFuncList = []cRegister.FunInfo{
		 	
			{
				Name: "GetAllLog", 
				InSchema: _LGetAllLogInSchema, 
				OutSchema: _LGetAllLogOutSchema }, 
		
	}
}

func (d *L) Do(ctx context.Context, fun string, param []byte) (result []byte, err error) {
	switch fun { 
 	
	case "GetAllLog":
		var in _LGetAllLogIn
		err = _LGetAllLogInCoder.DecodeIn(param, &in)
		if err != nil {
			panic(err)
		}
	
		res1, err := d.GetAllLog(
					
					
					)
		if err!=nil{
			return nil,err
		}
		var out _LGetAllLogOut
		 out.F1 = res1 
		
		
		res, err := _LGetAllLogOutCoder.Encode(out)
		if err != nil {
			panic(err)
		}
		return res, nil
	

	default:
		err = errors.New("rpc call error: fun not exist")
	}
	return
}

func (d *L) FuncList() []cRegister.FunInfo {
	return _LFuncList
}