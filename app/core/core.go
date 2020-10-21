package core

const (
	ServiceName = "CORE"
)

type SubService struct {
}

func (s *SubService) Invoke(fun string, param []byte) (result []byte, err error) {
	switch fun {
	case "Register":
		var si ServiceInfo
		err = serviceInfoCoder.DecodeIn(param, &si)
		if err != nil {
			panic(err)
		}

		err = s.Register(si)
		if err != nil {
			return
		}
		result, err = successCoder.Encode(true)
		return
	}

	result = []byte{1, 2, 3}
	return
}
