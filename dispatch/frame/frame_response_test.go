// Time : 2020/10/6 2:34
// Author : Kieran

// frame
package frame

import (
	"fmt"
	"testing"
)

// frame_response_test.go something
func TestResponseResult(t *testing.T) {
	d := []byte("test_-reqid")
	d = append(d, breakByte)
	d = append(d, breakByte)
	d = append(d, []byte("testResult")...)
	res, err := unMarshalResponse(d)
	fmt.Println(res, err)
}

func TestResponseError(t *testing.T) {
	d := []byte("test_-reqid")
	d = append(d, breakByte)
	d = append(d, []byte("test_-error")...)
	d = append(d, breakByte)
	res, err := unMarshalResponse(d)
	fmt.Println(res, err)
}

func TestResponseFailed(t *testing.T) {
	d := []byte("test_-reqid")
	d = append(d, []byte("test_-error")...)
	d = append(d, breakByte)
	res, err := unMarshalResponse(d)
	fmt.Println(res, err)
}
