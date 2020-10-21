// Time : 2020/9/26 19:47
// Author : Kieran

// coding
package coding

import (
	"github.com/hamba/avro"
)

// avro.go something

func NewAvro(rawSchema string) (c Coder, err error) {
	schema, err := avro.Parse(rawSchema)
	if err != nil {
		return
	}


	c = &AvroCoder{Schema: schema}

	return
}

type AvroCoder struct {
	Schema avro.Schema
}

func (c *AvroCoder) Encode(data interface{}) ([]byte, error) {
	return avro.Marshal(c.Schema, data)
}

func (c *AvroCoder) Decode(bytes []byte) (data interface{}, err error) {
	data = make(map[string]interface{})
	err = avro.Unmarshal(c.Schema, bytes, &data)
	return
}

func (c *AvroCoder) DecodeIn(bytes []byte, i interface{}) (err error) {
	err = avro.Unmarshal(c.Schema, bytes, &i)
	return
}

func init() {

//	schemaMap := make(map[uint8]*goavro.Codec)
//
//	sign(schemaMap)
//
//	signInfo(schemaMap)
//
//	reqCodec, err := goavro.NewCodec(`
//{
//	"namespace": "begonia.entry",
//	"type": "record",
//	"name": "Request",
//	"fields": [{
//			"name": "reqId",
//			"type": "string"
//		},
//		{
//			"name": "service",
//			"type": "string"
//		},
//		{
//			"name": "fun",
//			"type": "string"
//		},
//		{
//			"name": "params",
//			"type": "bytes"
//		}
//	]
//}`)
//	if err != nil {
//		panic("codec error!")
//	}
//	schemaMap[opcode.Request] = reqCodec
//
//	respCodec, err := goavro.NewCodec(`
//{
//	"namespace": "begonia.entry",
//	"type": "record",
//	"name": "Response",
//	"fields": [{
//			"name": "reqId",
//			"type": "string"
//		},
//		{
//			"name": "respErr",
//			"type": ["string","null"]
//		},
//		{
//			"name": "result",
//			"type": "bytes"
//		}
//	]
//}`)
//	if err != nil {
//		panic("codec error!")
//	}
//	schemaMap[opcode.Response] = respCodec

	//AvroCoder = &rAvroCoder{
	//	schemaMap: schemaMap,
	//}
}