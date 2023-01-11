package main

import (
	"context"
	"fmt"
	"github.com/MashiroC/begonia"
	"github.com/MashiroC/begonia/app/option"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

const (
	mode = "bgacenter"
)

var (
	count int64
	flag  bool
	l     sync.Mutex
)

func PprofWeb() {
	err := http.ListenAndServe(":9909", nil)
	if err != nil {
		panic(err)
	}
}

func QPS() {
	if !flag {
		l.Lock()
		if !flag {
			flag = true
			count=1
			go func() {
				step:=1
				for i:=0;i<6;i++{
					time.Sleep(10*time.Second)
					fmt.Printf("%d seconds QPS: %d\n",step*10,count/int64(step*10))
					step++
				}
				os.Exit(0)
				//time.Sleep(1*time.Second)
				//fmt.Println(count)
				//count=0
				//flag=false
			}()
		} else {
			atomic.AddInt64(&count, 1)
		}
		l.Unlock()
	} else {
		atomic.AddInt64(&count, 1)
	}
}

//go:generate begonia -c -s -r .
func main() {
	go PprofWeb()

	count = 0
	flag = false

	s := begonia.NewServer(option.Addr(":12306"),option.P2P())

	echoService := &EchoService{}
	testService := TestService(0)

	s.Register("Echo", echoService)
	s.Register("Test", &testService)

	s.Wait()
}

type EchoService struct {
}

func (*EchoService) SayHello(name string) (string, error) {
	QPS()
	return "Hello " + name, nil
}

func (*EchoService) SayHelloWithContext(ctx context.Context, name string) string {
	v := ctx.Value("info")
	info := v.(map[string]string)
	fmt.Println(info)
	fmt.Println(info["reqID"])
	fmt.Println(info["connID"])
	return "Hello ctx " + name
}

func (*EchoService) Add(i1, i2 int) (res int, err error) {
	res = i1 + i2
	return
}

func (*EchoService) Mod(i1, i2 int) (res1 int, res2 int) {
	res1 = i1 / i2
	res2 = i1 % i2
	return
}

func (*EchoService) NULL() {}

func (EchoService) Hello() {

}

type TestStruct struct {
	I1 int `json:"i_1"`
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
	B1 []byte
	B2 []uint8
}

type TestService int

func (*TestService) Echo(i1, it int, i2 int8, i3 int16, i4 int32, i5 int64,
	f1 float32, f2 float64, ok bool, str string,
	s1 []int, s2 []string, s6 []byte, st TestStruct,
	m1 map[string]string, m2 map[string]int, m3 map[string]TestStruct,
) string {
	QPS()
	//fmt.Println(i1, i2, i3, i4, i5)
	//fmt.Println(f1, f2, ok, str)
	//fmt.Println(s1, s2, s6, st)
	//fmt.Println(m1, m2, m3)
	return "ok"
}

func (*TestService) Echo2() (i1 int, i2 int8, i3 int16, i4 int32, i5 int64,
	f1 float32, f2 float64, ok bool, str string,
	s1 []int, s2 []string, s6 []byte, st TestStruct,
	m1 map[string]string, m2 map[string]int, m3 map[string]TestStruct,
) {
	//QPS()
	i1 = 1
	i2 = 2
	i3 = 3
	i4 = 4
	i5 = 5
	f1 = 6.0
	f2 = 7.0
	ok = true
	str = "test"
	s1 = []int{1, 2, 3}
	s2 = []string{"a", "s", "d"}
	s6 = []byte{4, 5, 6}
	st = TestStruct{
		I1:  9,
		I2:  8,
		I3:  7,
		I4:  6,
		I5:  5,
		Str: "shiina",
		S1:  []int{4, 5, 6},
		S2:  []string{"z", "x", "c"},
		TestStruct2: TestStruct2{
			B1: []byte{7, 8, 9},
			B2: []uint8{5, 6, 7},
		},
		Test3: TestStruct2{
			B1: []byte{7, 8, 9},
			B2: []uint8{5, 6, 7},
		},
		Map1: map[string]string{"test": "kieran"},
		Map2: map[string][]int{"hello": {3, 4, 5}},
	}

	m1 = map[string]string{"hello": "kieran"}
	m2 = map[string]int{"welcome": 1}
	m3 = map[string]TestStruct{"shiina": st}

	return
}
