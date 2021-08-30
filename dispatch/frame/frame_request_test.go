// Time : 2020/10/6 2:04
// Author : Kieran

// frame
package frame

import (
	"fmt"
	"testing"
)

// frame_request_test.go something

func TestRequest1(t *testing.T) {
	d := []byte("test_-reqid")
	d = append(d, breakByte)
	d = append(d, []byte("test_-server")...)
	d = append(d, breakByte)
	d = append(d, []byte("test_-fun")...)
	d = append(d, breakByte)
	d = append(d, []byte("test_-data")...)
	res, err := unMarshalRequest(d)
	fmt.Println(res.ReqID)
	fmt.Println(res.Fun)
	fmt.Println(res.Service)
	fmt.Println(string(res.Params))
	if err != nil || res == nil {
		t.Fail()
	}
}

func TestRequest2(t *testing.T) {
	// 缺少分隔符
	d := []byte("test_-reqid")
	d = append(d, breakByte)
	d = append(d, []byte("test_-server")...)
	d = append(d, []byte("test_-fun")...)
	d = append(d, []byte("test_-data")...)
	_, err := unMarshalRequest(d)
	if err == nil {
		t.Fail()
	}
}

func TestRequest3(t *testing.T) {
	// 中间参数空
	d := []byte("test_-reqid")
	d = append(d, breakByte)
	d = append(d, []byte("test_-server")...)
	d = append(d, breakByte)
	d = append(d, breakByte)
	d = append(d, []byte("test_-data")...)
	_, err := unMarshalRequest(d)
	if  err == nil {
		t.Fail()
	}
}
