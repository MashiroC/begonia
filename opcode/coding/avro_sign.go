// Time : 2020/8/3 13:36
// Author : MashiroC

// opcode
package coding

import (
	"begonia2/opcode"
	"github.com/linkedin/goavro/v2"
)

// avro_sign.go something

func sign(schemaMap map[uint8]*goavro.Codec) {
	signCodec, err := goavro.NewCodec(`
{
	"namespace": "begonia.entry",
	"type": "record",
	"name": "Sign",
	"fields": [{
			"name": "service",
			"type": "string"
		},
		{
			"name": "check",
			"type": "bytes"
		},
		{
			"name": "funs",
			"type": "array",
			"items": {
				"type": "record",
				"name": "fun",
				"fields": [{
						"name": "fun",
						"type": "string"
					},
					{
						"name": "inSchema",
						"type": "string"
					},
					{
						"name": "outSchema",
						"type": "string"
					}
				]
			}
		}
	]
}`)
	if err != nil {
		panic("codec error!")
	}
	schemaMap[opcode.Sign] = signCodec
}