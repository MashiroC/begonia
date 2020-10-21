package core

const (

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