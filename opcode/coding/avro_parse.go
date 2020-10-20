// Time : 2020/10/20 15:58
// Author : Kieran

// coding
package coding

import "reflect"

// avro_parse.go something

func toAvroSchemaField(t reflect.Type) string {
	return t.String()
}