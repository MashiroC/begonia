package coding

import (
	"errors"
	"fmt"
)

var (
	empty       = &emptyCoder{}
	emptyReturn = []byte{1}
)

type emptyCoder struct{}

func (e *emptyCoder) Encode(data interface{}) (b []byte, err error) {
	if data != nil {
		m, ok := data.(map[string]interface{})
		if ok && len(m) != 0 {
			err = fmt.Errorf("input length need 0 but %d", len(m))
		}
	}

	b = emptyReturn
	return
}

func (e *emptyCoder) Decode(bytes []byte) (data interface{}, err error) {
	if len(bytes) == 1 && bytes[0] == 1 {
		return map[string]interface{}{}, nil
	}
	return nil, errors.New("unknow decode error")
}

func (e *emptyCoder) DecodeIn(bytes []byte, i interface{}) error {
	if len(bytes) == 1 && bytes[0] == 1 {
		ptr,ok := i.(*map[string]interface{})
		if ok {
			*ptr = map[string]interface{}{}
		}
		return nil
	}
	return errors.New("unknow decode error")
}
