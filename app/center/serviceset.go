package center

type serviceSet struct {
}

func newServiceSet() *serviceSet {
	return &serviceSet{}
}

func (s *serviceSet) Get(service string) (connID string, ok bool) {
	return
}

func (s *serviceSet) Add(service string) (err error) {
	return
}
