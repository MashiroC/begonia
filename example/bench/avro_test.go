// Time : 2020/9/27 17:35
// Author : Kieran

// bench
package bench

import (
	"bytes"
	"fmt"
	"github.com/MashiroC/begonia/dispatch/frame"
	"github.com/MashiroC/begonia/tool/berr"
	"github.com/hamba/avro"
	"log"
	"strconv"
	"testing"
	"time"
)

// avro_test.go something

//func BenchmarkLinkedinEncode(b *testing.B) {
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		linkedinEncode()
//	}
//}
//
//func BenchmarkLinkedinDecode(b *testing.B){
//	b.ResetTimer()
//	for i:=0;i<b.N;i++{
//		linkedinDecode()
//	}
//}

//func BenchmarkHambaEncode(b *testing.B) {
//	hambaEncode()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		hambaEncode()
//	}
//}

//func BenchmarkHambaDecode(b *testing.B){
//	hambaDecode()
//	b.ResetTimer()
//	for i:=0;i<b.N;i++{
//		hambaDecode()
//	}
//}

func TestLinkedinEncode(t *testing.T) {
	linkedinEncode()
}

func TestLinkedinDecode(T *testing.T) {
	linkedinDecode()
}

func TestHambaEncode(t *testing.T) {
	hambaEncode()
}

func TestHambaDecode(t *testing.T) {
	hambaDecode()
}

func Test(t *testing.T) {
	sIn := `{"type":"string","name":"serviceName"}`
	c, _ := avro.Parse(sIn)
	res, _ := avro.Marshal(c, "hello")
	fmt.Println(res)

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
	b, err := avro.Marshal(o, tt)
	if err != nil {
		panic(err)
	}

	fmt.Println(b)

	var ttt []test
	err = avro.Unmarshal(o, b, &ttt)
	if err != nil {
		panic(err)
	}

	fmt.Println(ttt)
	//var f =  func(v interface{}){
	//	fmt.Println(reflect2.TypeOf(v).Kind())
	//}
	//f(map[string]interface{}{})

}

type rFun struct {
	Name      string `avro:"name"`
	InSchema  string `avro:"inSchema"`
	OutSchema string `avro:"outSchema"`
}

func TestByte(t *testing.T) {
	//tmp:=make([]bool,256)
	//str:="qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM_-1234567890"
	//for i:=0;i<len(str);i++{
	//	tmp[str[i]]=true
	//}
	//
	//for i:=0;i<256;i++{
	//	if !tmp[i]{
	//		fmt.Printf("%d ",i)
	//	}
	//}
	//str:=string([]byte{0x00})
	//fmt.Printf(str)

	//fmt.Printf("%b\n", 0xfa)
	//fmt.Printf("%b\n", 15)
	//fmt.Printf("%b\n", 10)
	//fmt.Printf("%b\n", 15<<4)

	typCode := 1 // 0 ~ 1

	dispatchCode := 4 // 0 ~ 7

	version := 8 // 0 ~ 15

	opcode := ((typCode<<3)|dispatchCode)<<4 | version
	fmt.Printf("opcode: %08b %d\n", opcode, opcode)

	version = opcode & 0b00001111
	fmt.Printf("versionCode:%04b %d\n", version, version)

	dispatchCode = opcode >> 4 & 0b0111
	fmt.Printf("dispatchCode:%03b %d\n", dispatchCode, dispatchCode)
	//
	typCode = opcode >> 7
	fmt.Printf("typCode:%01b %d\n", typCode, typCode)

	//l:=500
	//fmt.Printf("%b\n",l)
	//
	//l2:=l
	//for l2>255{
	//	l2=l2>>1
	//}
	//fmt.Printf("%b\n",l2)
}

func TestPPT(t *testing.T) {

	schema := avro.MustParse(`
{
	"namespace": "example.avro",
	"type": "record",
	"name": "User",
	"fields": [
		 {"name": "name", "type": "string"},
		 {"name": "age",  "type": "int"}
	]
}`)

	res, err := avro.Marshal(schema, People{
		Name: "kieran",
		Age:  18,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("0x%x\n", res)

	type RemoteFun struct {
		Name      string `avro:"name"`
		InSchema  string `avro:"inSchema"`
		OutSchema string `avro:"outSchema"`
	}

}

type People struct {
	Name string `avro:"name"`
	Age  int    `avro:"age"`
}

type Service int

func (*Service) AddAge(p People) People {
	p.Age += 1
	return p
}

func TestSelect(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan int)
	close(ch1)
	close(ch2)
	select {
	case res, ok := <-ch1:
		fmt.Println("ch1", res, ok)
	case res, ok := <-ch2:
		fmt.Println("ch2", res, ok)
	}
	fmt.Println("exit")
}

func TestResp(t *testing.T) {
	typCode := 1
	dispatchCode := frame.CtrlDefaultCode // 0 ~ 7

	version := frame.ProtocolVersion // 0 ~ 15

	res := ((typCode<<3)|dispatchCode)<<4 | version
	fmt.Println(res)
}

type FunInfo struct {
	Fun       string `avro:"fun"`
	Mode      string `avro:"mode"`
	InSchema  string `avro:"inSchema"`
	OutSchema string `avro:"outSchema"`
}

type ServiceInfo struct {
	Service string    `avro:"service"`
	Funs    []FunInfo `avro:"funs"`
}

func TestAvroStruct(t *testing.T) {
	rawSchema := `[{
		"type": "record",
		
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
		"namespace": "github.com/MashiroC/begonia.entry",
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
					"items": "string"
				}
			}
		]
	}
]`

	schema := avro.MustParse(rawSchema)

	obj := ServiceInfo{
		Service: "test",
		//Funs:    []F{"asd","zxc"},
	}

	//obj := map[string]interface{}{
	//	"service":"test",
	//	"funs":[]map[string]interface{}{{
	//		"fun":"test1",
	//		"mode":"avro",
	//		"inSchema":"asdasd",
	//		"outSchema":"asdasd",
	//	}},
	//}

	b, err := avro.Marshal(schema, &obj)
	if err != nil {
		panic(err)
	}
	fmt.Println(b)
}

func TestAvroSt(t *testing.T) {
	//rawSchema:=
	//s:=avro.MustParse(rawSchema)
	//fmt.Println(s)

	//b, err := avro.Marshal(s, ServiceInfoCall{
	//	Service: "test",
	//	Funs: []FunInfo{{
	//		Fun:       "asd",
	//		Mode:      "zxc",
	//		InSchema:  "zxcv",
	//		OutSchema: "zxcv",
	//	}},
	//})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(b)
}

func TestForMap(t *testing.T) {
	m := make(map[string]*People)

	m["in1"] = &People{
		Name: "shiina",
		Age:  18,
	}

	var i int64 = 1
	ok := true
	for ok {
		v, ok := m["in"+strconv.FormatInt(i, 10)]
		fmt.Println(v, ok)
		i++
		time.Sleep(5 * time.Second)
	}
}

func TestInt64(t *testing.T) {
	rawSchema := `
{
			"namespace":"github.com/MashiroC/begonia.func.Test",
			"type":"record",
			"name":"In",
			"fields":[
				{"name":"Test","type":"int"}
			]
		}`
	schema:=avro.MustParse(rawSchema)
	obj:= struct {
		Test uint
	}{Test: 123456}

	res,err:=avro.Marshal(schema,obj)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func TestBErr(t *testing.T) {
	berr.New("error system","test","you are in a black hole")

	berr.NewAuto("auto","you are in a black hole")
}

//func BenchmarkBerr(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		berr.New("error system","test","you are in a black hole")
//
//	}
//}
//
//func BenchmarkBErrAuto(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		berr.NewAuto("auto","you are in a black hole")
//	}
//}

func BenchmarkLog1(b *testing.B) {
	log.SetFlags(log.Ldate|log.Llongfile)
	tmp:=bytes.NewBuffer([]byte{})
	log.SetOutput(tmp)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Println("hello")
	}
}

func BenchmarkLog2(b *testing.B) {
	log.SetFlags(log.Ldate)
	tmp:=bytes.NewBuffer([]byte{})
	log.SetOutput(tmp)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Println("hello")
	}
}