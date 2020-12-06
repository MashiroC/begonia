// Time : 2020/9/27 19:13
// Author : Kieran

// bench
package bench

import (
	"bytes"
	"github.com/actgardner/gogen-avro/v7/soe"
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

	reqAst = TestReq{
		ReqId:   "test",
		Service: "logService",
		Fun:     "ServiceInfoCall",
		Params:  []byte{1, 2, 3},
	}

	reqBin = []byte{8, 116, 101, 115, 116, 22, 67, 111, 114, 101, 83, 101, 114, 118, 105, 99, 101, 16, 83, 105, 103, 110, 73, 110, 102, 111, 6, 1, 2, 3}
	reqBi1 = []byte{8, 116, 101, 115, 116, 20, 108, 111, 103, 83, 101, 114, 118, 105, 99, 101, 30, 83, 101, 114, 118, 105, 99, 101, 73, 110, 102, 111, 67, 97, 108, 108, 6, 1, 2, 3}
)

func init() {
	rawSchema := `
{
	"namespace": "begonia.entry",
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
	//b, err := avro.Marshal(ReqSchema, reqNative)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(b)
}

func hambaDecode() {
	var res TestReq
	//var res map[string]interface{}
	//res = make(map[string]interface{})
	//avro.Unmarshal(ReqSchema, reqBin, res)

	err := avro.Unmarshal(ReqSchema, reqBin, &res)
	if err != nil {
		panic(err)
	}
	//fmt.Println(res)
}

func AstEncode() {
	buf := bytes.NewBuffer(make([]byte, 0, 10))
	err := reqAst.Serialize(buf)
	if err != nil {
		panic(err)
	}
	//res := buf.Bytes()
	buf.Bytes()
	//tmp := []byte("{")
	//for i := 0; i < len(res); i++ {
	//	tmp = append(tmp, []byte(qconv.I2S(int(res[i])))...)
	//	tmp = append(tmp, byte(','))
	//}
	//tmp = tmp[:len(tmp)-1]
	//tmp = append(tmp, byte('}'))
	//fmt.Println(len(res))
	//fmt.Println(string(tmp))

}

func AstDecode() {
	buf := bytes.NewBuffer(reqBi1)
	//res, err := DeserializeTestReq(buf)
	//fmt.Println(res, err)
	//res, err = DeserializeTestReq(buf)
	//fmt.Println(res,err)
	soe.ReadHeader(buf)
	DeserializeTestReq(buf)
	//_, err := DeserializeTestReq(buf)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(res)
}

func hamba() {
	req := TestReq{
		ReqId:   "test",
		Service: "logService",
		Fun:     "ServiceInfoCall",
		Params:  []byte{1, 2, 3},
	}
	b, err := avro.Marshal(ReqSchema, req)
	if err != nil {
		panic(err)
	}
	//fmt.Println(b)

	var req2 TestReq
	err = avro.Unmarshal(ReqSchema, b, &req2)
	if err != nil {
		panic(err)
	}
	//fmt.Println(req2)
}

func ast() {
	req := TestReq{
		ReqId:   "test",
		Service: "logService",
		Fun:     "ServiceInfoCall",
		Params:  []byte{1, 2, 3},
	}
	buf := &bytes.Buffer{}
	err := req.Serialize(buf)
	if err != nil {
		panic(err)
	}
	//fmt.Println(buf.Bytes())
	//res, err := DeserializeTestReq(buf)
	//if err != nil {
	//	panic(err)
	//}
	DeserializeTestReq(buf)

	//fmt.Println(res)
}
