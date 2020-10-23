package coding

import "errors"

// SuccessCoder 对于一些只需要bool值的类型，直接使用successCoder payload只有error或bool
type SuccessCoder struct {
}

// Encode success encode
func (s *SuccessCoder) Encode(data interface{}) ([]byte, error) {
	if res, ok := data.(bool); ok && res {
		return []byte{1}, nil
	}
	return []byte{0}, nil
}

// Decode Success decode
func (s *SuccessCoder) Decode(bytes []byte) (data interface{}, err error) {
	if len(bytes) != 1 {
		data = false
		err = errors.New("resp byte len decode error")
	} else {
		if bytes[0] == 1 {
			data = true
		} else if bytes[0] == 0 {
			data = false
		} else {
			err = errors.New("resp byte decode error")
		}
	}
	return
}

// DecodeIn success DecodeIN
func (s *SuccessCoder) DecodeIn(bytes []byte, i interface{}) (err error) {
	if len(bytes) != 1 {
		i = false
		err = errors.New("resp byte len decode error")
	} else {
		if bytes[0] == 1 {
			i = true
		} else if bytes[0] == 0 {
			i = false
		} else {
			err = errors.New("resp byte decode error")
		}
	}
	return
}
