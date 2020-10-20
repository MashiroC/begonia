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
	serviceInfoRawSchema = `{
	"namespace": "begonia.entry",
	"type": "record",
	"name": "ServiceInfo",
	"fields": [{
			"name": "service",
			"type": "string"
		},
		{
			"name": "funs",
			"type": {
				"type": "array",
				"items": {
					"type": "record",
					"name": "FunInfo",
					"fields": [{
							"name": "name",
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
				}
			}
		}
	]
}`
)
