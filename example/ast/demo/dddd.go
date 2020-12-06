package demo

//type Demo struct {
//}

//func (d *Demo) Echo(str1,str2 string) string {
//	return "Hello, " + str1+str2
//}
//
//func (d *Demo) Add(num int) (res int, err error) {
//	return num, errors.New("test")
//}

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

type EchoService struct {
}

//go:generate cmd.exe
func (*EchoService) Echo(i1 int, i2 int8, i3 int16, i4 int32, i5 int64,
	f1 float32, f2 float64, ok bool, str string,
	s1 []int, s2 []string, s6 []byte, st TestStruct,
	m1 map[string]string, m2 map[string]int, m3 map[string]TestStruct,
) {
	//st TestStruct, st2 *TestStruct
	//m1 map[string]string, m2 map[string]int, m3 map[string]TestStruct
}

func (*EchoService) Echo2() (i1 int, i2 int8, i3 int16, i4 int32, i5 int64,
	f1 float32, f2 float64, ok bool, str string,
	s1 []int, s2 []string, s6 []byte, st TestStruct,
	m1 map[string]string, m2 map[string]int, m3 map[string]TestStruct,
) {
	return
}

func TestEcho() {

}

type TestService struct {
	str string
}
