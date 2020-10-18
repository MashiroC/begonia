// Time : 2020/10/6 2:04
// Author : Kieran

// frame
package frame

import (
	"testing"
)

// frame_request_test.go something

func TestRequest1(t *testing.T) {
	d := []byte("test_-reqid")
	d = append(d, breakByte)
	d = append(d, []byte("test_-service")...)
	d = append(d, breakByte)
	d = append(d, []byte("test_-fun")...)
	d = append(d, breakByte)
	d = append(d, []byte("test_-data")...)
	res, err := unMarshalRequest(d)
	if err!=nil || res==nil {
		t.Fail()
	}
}

func TestRequest2(t *testing.T) {
	// 缺少分隔符
	d := []byte("test_-reqid")
	d = append(d, breakByte)
	d = append(d, []byte("test_-service")...)
	d = append(d, []byte("test_-fun")...)
	d = append(d, []byte("test_-data")...)
	res, err := unMarshalRequest(d)
	if res!=nil || err==nil {
		t.Fail()
	}
}

func TestRequest3(t *testing.T) {
	// 中间参数空
	d := []byte("test_-reqid")
	d = append(d, breakByte)
	d = append(d, []byte("test_-service")...)
	d = append(d, breakByte)
	d = append(d, breakByte)
	d = append(d, []byte("test_-data")...)
	res, err := unMarshalRequest(d)
	if res!=nil || err==nil{
		t.Fail()
	}
}
