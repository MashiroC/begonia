// Time : 2020/10/6 17:59
// Author : Kieran

// app
package app

// coreservice_schema.go something

const (
	signInfoRawSchema = `{
    "namespace":"begonia.entry",
    "type":"record",
    "name":"SignInfoReq",
    "fields":[
        {
            "name":"service",
            "type":"string"
        }
    ]
}`
	signInfoResultRawSchema = `[{
		"type": "record",
		"namespace": "begonia.entry",
		"name": "FunInfo",
		"fields": [{
				"name": "fun",
				"type": "string"
			},
			{
				"name": "mode",
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
	},
	{
		"namespace": "begonia.entry",
		"type": "record",
		"name": "SignInfoResp",
		"fields": [{
				"name": "service",
				"type": "string"
			},
			{
				"name": "address",
				"type": {
					"type": "array",
					"items": "begonia.entry.FunInfo"
				}
			}
		]
	}
]`
)
