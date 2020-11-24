// Time : 2020/9/27 19:13
// Author : Kieran

// bench
package bench

import (
	"github.com/hamba/avro"
	"github.com/linkedin/goavro/v2"
)

// benchi.go something

var (
	signInfoReqCodec   *goavro.Codec
	ReqCodec           *goavro.Codec
	signInfoParamCodec *goavro.Codec

	ReqSchema avro.Schema

	signInfo = map[string]interface{}{
		"reqId":   "test",
		"service": "logService",
	}

	signParam = map[string]interface{}{
		"service": "logService",
	}

	req = map[string]interface{}{
		"reqId":   "test",
		"service": "CoreService",
		"fun":     "ServiceInfoCall",
		"params":  []byte{1, 2, 3},
	}

	reqNative = TestReq{
		ReqId:   "test",
		Service: "logService",
		Fun:     "ServiceInfoCall",
		Params:  []byte{1, 2, 3},
	}

	reqBin = []byte{8, 116, 101, 115, 116, 22, 67, 111, 114, 101, 83, 101, 114, 118, 105, 99, 101, 16, 83, 105, 103, 110, 73, 110, 102, 111, 6, 1, 2, 3}
)

func init() {
	rawSchema := `
{
	"namespace": "github.com/MashiroC/begonia.entry",
	"type": "record",
	"name": "Request",
	"fields": [{
			"name": "reqId",
			"type": "string"
		},
		{
			"name": "service",
			"type": "string"
		},
		{
			"name": "fun",
			"type": "string"
		},
		{
			"name": "params",
			"type": "bytes"
		}
	]
}`
	var err error

	ReqCodec, err = goavro.NewCodec(rawSchema)
	if err != nil {
		panic(err)
	}

	ReqSchema = avro.MustParse(rawSchema)
}

func linkedinEncode() {
	//res, err := ReqCodec.BinaryFromNative(nil, req)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(res)
	ReqCodec.BinaryFromNative(nil, req)

}

func linkedinDecode() {
	ReqCodec.NativeFromBinary(reqBin)

	//res, _, err := ReqCodec.NativeFromBinary(reqBin)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(res)
}

//func linkedinDecode(){
//
//}

type TestReq struct {
	ReqId   string `avro:"reqId"`
	Service string `avro:"service"`
	Fun     string `avro:"fun"`
	Params  []byte `avro:"params"`
}

func hambaEncode() {
	//res, err := avro.Marshal(ReqSchema, req)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Print("[]byte{")
	//for i := 0; i < len(res); i++ {
	//	fmt.Print(res[i], ",")
	//}
	//fmt.Println(res)

	/*
		1. 1 string 1int < 4
		2. 单结构体
		3. other

	*/

	//var reqIn interface{}= reqNative
	avro.Marshal(ReqSchema, req)
}

func hambaDecode() {
	//var res TestReq
	var res map[string]interface{}
	res = make(map[string]interface{})
	avro.Unmarshal(ReqSchema, reqBin, res)

	//err := avro.Unmarshal(ReqSchema, reqBin, &res)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(res)
}
