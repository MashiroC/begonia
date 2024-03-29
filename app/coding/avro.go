package coding

import (
	"fmt"
	"github.com/hamba/avro"
)

const (
	EmptyAvroSchema = "EMPTY_AVRO_SCHEMA"
)

// NewAvro 使用avro模式创建一个新的coder
func NewAvro(rawSchema string) (c Coder, err error) {
	if rawSchema == EmptyAvroSchema {
		c = empty
		return
	}

	schema, err := avro.Parse(rawSchema)
	if err != nil {
		return
	}

	c = &avroCoder{Schema: schema}

	return
}

// avroCoder Avro模式的coder
type avroCoder struct {
	Schema avro.Schema
}

func (c *avroCoder) Encode(data interface{}) ([]byte, error) {
	return avro.Marshal(c.Schema, data)
}

func (c *avroCoder) Decode(bytes []byte) (data interface{}, err error) {
	data = make(map[string]interface{})
	err = avro.Unmarshal(c.Schema, bytes, &data)
	return
}

func (c *avroCoder) DecodeIn(bytes []byte, i interface{}) (err error) {
	err = avro.Unmarshal(c.Schema, bytes, &i)
	return
}

// ToAvroObj 将参数转化为适用于avro的结构
func ToAvroObj(params []interface{}) interface{} {
	out := make(map[string]interface{})
	for i := 0; i < len(params); i++ {
		out["F"+fmt.Sprintf("%d", i+1)] = params[i]
	}
	return out
}
