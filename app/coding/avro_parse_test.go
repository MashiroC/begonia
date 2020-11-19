package coding

import (
	"fmt"
	"github.com/hamba/avro"
	"github.com/mitchellh/mapstructure"
	"github.com/modern-go/reflect2"
	"reflect"
	"testing"
)

type TestStruct struct {
	I1 int
	I2 int8
	I3 int16
	I4 int32
	I5 int64

	Str string
	S1  []int
	S2  []string

	TestStruct2
	Test3 TestStruct2

	Map1 map[string]string
	Map2 map[string][]int
}

type TestStruct2 struct {
	b1 []byte
	b2 []uint8
}

type EchoService int

func (*EchoService) Echo(i1 int, i2 int8, i3 int16, i4 int32, i5 int64,
	f1 float32, f2 float64, ok bool, str string,
	s1 []int, s2 []string, s3 []*string, s6 []byte, st TestStruct, stp *TestStruct,
	m1 map[string]string, m2 map[string]int, m3 map[string]TestStruct,
) {
	//st TestStruct, st2 *TestStruct
	//m1 map[string]string, m2 map[string]int, m3 map[string]TestStruct
}

func (*EchoService) Echo2() (i1 int, i2 int8, i3 int16, i4 int32, i5 int64,
	f1 float32, f2 float64, ok bool, str string,
	s1 []int, s2 []string, s3 []*string, s6 []byte, st TestStruct, stp *TestStruct,
	m1 map[string]string, m2 map[string]int, m3 map[string]TestStruct,
) {
	return
}

func TestParse(t *testing.T) {
	type Input struct {
		F1  int                   `avro:"f1"`
		F2  int8                  `avro:"f2"`
		F3  int16                 `avro:"f3"`
		F4  int32                 `avro:"f4"`
		F5  int64                 `avro:"f5"`
		F6  float32               `avro:"f6"`
		F7  float64               `avro:"f7"`
		F8  bool                  `avro:"f8"`
		F9  string                `avro:"f9"`
		F10 []int                 `avro:"f10"`
		F11 []string              `avro:"f11"`
		F12 []*string             `avro:"f12"`
		F13 []byte                `avro:"f13"`
		F14 TestStruct            `avro:"f14"`
		F15 *TestStruct           `avro:"f15"`
		F16 map[string]string     `avro:"f16"`
		F17 map[string]int        `avro:"f17"`
		F18 map[string]TestStruct `avro:"f18"`
	}

	obj := Input{
		F1:  1,
		F2:  2,
		F3:  3,
		F4:  4,
		F5:  5,
		F6:  6,
		F7:  7,
		F8:  true,
		F9:  "test",
		F10: []int{1, 2, 3},
		F11: []string{"a", "s", "d"},
		F12: []*string{},
		F13: []byte{4, 5, 6},
		F14: TestStruct{
			I1:  9,
			I2:  8,
			I3:  7,
			I4:  6,
			I5:  5,
			Str: "shiina",
			S1:  []int{4, 5, 6},
			S2:  []string{"z", "x", "c"},
			TestStruct2: TestStruct2{
				b1: []byte{7, 8, 9},
				b2: []uint8{5, 6, 7},
			},
			Test3: TestStruct2{
				b1: []byte{7, 8, 9},
				b2: []uint8{5, 6, 7},
			},
			Map1: map[string]string{"test": "kieran"},
			Map2: map[string][]int{"hello": {3, 4, 5}},
		},
		F15: &TestStruct{
			I1:  9,
			I2:  8,
			I3:  7,
			I4:  6,
			I5:  5,
			Str: "shiina",
			S1:  []int{4, 5, 6},
			S2:  []string{"z", "x", "c"},
			TestStruct2: TestStruct2{
				b1: []byte{7, 8, 9},
				b2: []uint8{5, 6, 7},
			},
			Test3: TestStruct2{
				b1: []byte{7, 8, 9},
				b2: []uint8{5, 6, 7},
			},
			Map1: map[string]string{"test": "kieran"},
			Map2: map[string][]int{"hello": {3, 4, 5}},
		},
		F16: map[string]string{"hello": "kieran"},
		F17: map[string]int{"welcome": 1},
		F18: map[string]TestStruct{},
	}

	s := EchoService(1)
	e := &s
	typ := reflect.TypeOf(e)
	m := typ.Method(0)

	rawSchema := InSchema(m)
	schema := avro.MustParse(rawSchema)
	res, err := avro.Marshal(schema, obj)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(res), res)

	m2 := typ.Method(1)
	rawSchema2, _ := OutSchema(m2)
	schema2 := avro.MustParse(rawSchema2)
	res2, err := avro.Marshal(schema2, obj)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(res2), res2)

	var obj2 Input
	err = avro.Unmarshal(schema2, res2, &obj2)
	if err != nil {
		panic(err)
	}
	fmt.Println(obj2)
}

func TestAvroSliceParse(t *testing.T) {
	var arr []interface{}
	arr = []interface{}{1, 2, 3}
	typ := reflect.TypeOf(arr)
	//childTyp := reflect.TypeOf(arr[0])
	sTyp := reflect.TypeOf([]int{})
	fmt.Println(typ, sTyp)
	slice := reflect.MakeSlice(sTyp, 0, 2)

	fmt.Println(slice.Type())
	for i := 0; i < len(arr); i++ {
		v := arr[i]
		rv := reflect.ValueOf(v)
		slice = reflect.Append(slice, rv)
		fmt.Println(rv.Type())
	}
	fmt.Println(slice)
}

func TestReflectPtr(t *testing.T) {
	type People struct {
		Name string
		Age  int
	}

	pTyp := reflect.TypeOf(People{})

	var m interface{}
	m = map[string]interface{}{"Name": "asd", "Age": 123}

	in := reflect.New(pTyp)
	people:=in.Interface()
	err := mapstructure.Decode(m, &people)
	if err != nil {
		panic(err)
	}

	fmt.Println(reflect2.Type2(pTyp).Indirect(people))
	//fmt.Println(reflect.ValueOf(people).Elem().Interface())
	fmt.Println(people)
	fmt.Println(in.Interface())
	fmt.Println(in.Elem().Interface())

	//v:=reflect.ValueOf(p)

	//fmt.Println(v)
}
