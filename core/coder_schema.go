package core

// coder_schema.go 硬编码的函数的原始avro schema

const (

	// 获取服务的调用
	serviceInfoCallRawSchema = `{
    "namespace":"begonia.entry",
    "type":"record",
    "name":"ServiceInfoCall",
    "fields":[
        {
            "name":"service",
            "type":"string"
        }
    ]
}`

	// 获取服务的结果
	serviceInfoRawSchema = `{
	"namespace": "begonia.entry",
	"type": "record",
	"name": "ServiceInfoCall",
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
