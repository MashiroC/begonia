package coding

import (
	"fmt"
	"testing"
)

func TestEmpty(t *testing.T) {
	fmt.Println(empty.Encode(map[string]interface{}{}))
	fmt.Println(empty.Decode([]byte{1}))

	var m map[string]interface{}

	err := empty.DecodeIn([]byte{1},&m)
	fmt.Println(err)
	fmt.Println(m)
}