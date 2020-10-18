// Time : 2020/9/26 19:47
// Author : Kieran

// coding
package coding

import (
	"begonia2/opcode"
	"fmt"
	"github.com/hamba/avro"
	"github.com/linkedin/goavro/v2"
)

// avro.go something

func NewAvro(rawSchema string) (c Coder, err error) {
	schema, err := avro.Parse(rawSchema)
	if err != nil {
		return
	}

	c = &AvroCoder{schema: schema}

	return
}

type AvroCoder struct {
	schema avro.Schema
}

func (c *AvroCoder) Encode(data interface{}) ([]byte, error) {
	return avro.Marshal(c.schema, data)
}

func (c *AvroCoder) Decode(bytes []byte) (data interface{}, err error) {
	data = make(map[string]interface{})
	err = avro.Unmarshal(c.schema, bytes, &data)
	return
}

func (c *AvroCoder) DecodeIn(bytes []byte, i interface{}) (err error) {
	err = avro.Unmarshal(c.schema, bytes, &i)
	return
}

func init() {

	schemaMap := make(map[uint8]*goavro.Codec)

	sign(schemaMap)

	signInfo(schemaMap)

	reqCodec, err := goavro.NewCodec(`
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
}`)
	if err != nil {
		panic("codec error!")
	}
	schemaMap[opcode.Request] = reqCodec

	respCodec, err := goavro.NewCodec(`
{
	"namespace": "begonia.entry",
	"type": "record",
	"name": "Response",
	"fields": [{
			"name": "reqId",
			"type": "string"
		},
		{
			"name": "respErr",
			"type": ["string","null"]
		},
		{
			"name": "result",
			"type": "bytes"
		}
	]
}`)
	if err != nil {
		panic("codec error!")
	}
	schemaMap[opcode.Response] = respCodec

	//AvroCoder = &rAvroCoder{
	//	schemaMap: schemaMap,
	//}
}

type rAvroCoder struct {
	schemaMap map[uint8]*goavro.Codec
}

func (c *rAvroCoder) Decode(opcode uint8, data []byte) (m map[string]interface{}, err error) {
	codec, ok := c.schemaMap[opcode]
	if !ok {
		err = fmt.Errorf("opcode [%d] not in avro schema map", opcode)
		return
	}
	res, _, err := codec.NativeFromBinary(data)
	if err != nil {
		return
	}
	m = res.(map[string]interface{})
	return
}

func (c *rAvroCoder) Encode(opcode uint8, data map[string]interface{}) (b []byte, err error) {
	codec, ok := c.schemaMap[opcode]
	if !ok {
		err = fmt.Errorf("opcode [%d] not in avro schema map", opcode)
		return
	}

	b, err = codec.BinaryFromNative(nil, data)
	return
}
