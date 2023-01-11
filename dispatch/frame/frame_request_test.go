// Time : 2020/10/6 2:04
// Author : Kieran

// frame
package frame

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// frame_request_test.go something

func TestHeader(t *testing.T) {
	header := strings.Join([]string{"key1", "value1", "key2", "value2"}, string(rune(headerBreakByte)))

	h,err:=unMarshalHeader([]byte(header))
	if err != nil {
	    panic(err)
	}
	fmt.Println(h)
}

func TestRequest1(t *testing.T) {
	header := strings.Join([]string{"key", "value"}, string(rune(headerBreakByte)))
	reqID := "test_-reqid"
	service := "test_-server1"
	fun := "test_-fun"
	data := "test_-data"

	target := []string{header, reqID, service, fun, data}

	d := []byte{}

	for i, b := range target {
		d = append(d, b...)
		if i != len(target)-1 {
			d = append(d, breakByte)
		}
	}
	req, err := unMarshalRequest(d)
	a := assert.New(t)
	a.Nil(err)
	a.Equal(req.ReqID, reqID)
	a.Equal(req.Service, service)
	a.Equal(req.Fun, fun)
	a.Equal(string(req.Params), data)
	a.Len(req.Header, 1)

	//fmt.Println(req)
}

func TestRequest2(t *testing.T) {
	// 缺少分隔符
	d := []byte("test_-reqid")
	d = append(d, breakByte)
	d = append(d, []byte("test_-server1")...)
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
	d = append(d, []byte("test_-server1")...)
	d = append(d, breakByte)
	d = append(d, breakByte)
	d = append(d, []byte("test_-data")...)
	_, err := unMarshalRequest(d)
	if err == nil {
		t.Fail()
	}
}
