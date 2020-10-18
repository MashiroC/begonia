// Time : 2020/8/3 13:37
// Author : MashiroC

// opcode
package coding

import (
	"begonia2/opcode"
	"github.com/linkedin/goavro/v2"
)

// avro_sign_info.go something

func signInfo(schemaMap map[uint8]*goavro.Codec) {
	signInfoReqCodec, err := goavro.NewCodec(`
{
    "namespace":"begonia.entry",
    "type":"record",
    "name":"SignInfoReq",
    "fields":[
        {
            "name":"reqId",
            "type":"string"
        },
        {
            "name":"service",
            "type":"string"
        }
    ]
}`)
	if err != nil {
		panic("codec error!")
	}
	schemaMap[opcode.SignInfoReq] = signInfoReqCodec

	signInfoRespCodec,err:=goavro.NewCodec(`
{
    "namespace":"begonia.entry",
    "type":"record",
    "name":"SignInfoResp",
    "fields":[
        {
            "name":"reqId",
            "type":"string"
        },
        {
            "name":"service",
            "type":"string"
        },
        {
            "name":"funs",
            "type":"array",
            "items":{
                "type":"record",
                "name":"funInfo",
                "fields":[
                    {
                        "name":"fun",
                        "type":"string"
                    },
                    {
                        "name":"inSchema",
                        "type":"string"
                    },
                    {
                        "name":"outSchema",
                        "type":"string"
                    }
                ]
            }
        }
    ]
}`)

	if err != nil {
		panic("codec error!")
	}
	schemaMap[opcode.SignInfoResp] = signInfoRespCodec
}
