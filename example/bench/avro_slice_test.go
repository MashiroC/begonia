// Time : 2020/9/28 21:14
// Author : Kieran

// bench
package bench

import (
	"github.com/hamba/avro"
	"testing"
)

// avro_slice_test.go something

func BenchmarkMake(b *testing.B) {

	sOut := `{"type":"array","name":"test","items":{
	"type":"record",
    "name":"test1",
    "fields":[
	{"type":"string","name":"fun"},
    {"type":"string","name":"test1"}
]
}}`

	type test struct {
		Fun   string `avro:"fun"`
		Test1 string `avro:"test1"`
	}

	var tt []test
	tt = []test{{Fun: "s1", Test1: "t1"}, {Fun: "s2", Test1: "t2"}}

	o, err := avro.Parse(sOut)
	if err != nil {
		panic(err)
	}
	bt, err := avro.Marshal(o, tt)
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ttt := make([]test, 0)
		err = avro.Unmarshal(o, bt, &ttt)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkNoMake(b *testing.B) {
	sOut := `{"type":"array","name":"test","items":{
	"type":"record",
    "name":"test1",
    "fields":[
	{"type":"string","name":"fun"},
    {"type":"string","name":"test1"}
]
}}`

	type test struct {
		Fun   string `avro:"fun"`
		Test1 string `avro:"test1"`
	}

	var tt []test
	tt = []test{{Fun: "s1", Test1: "t1"}, {Fun: "s2", Test1: "t2"}}

	o, err := avro.Parse(sOut)
	if err != nil {
		panic(err)
	}
	bt, err := avro.Marshal(o, tt)
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var ttt []test
		err = avro.Unmarshal(o, bt, &ttt)
		if err != nil {
			panic(err)
		}
	}
}
