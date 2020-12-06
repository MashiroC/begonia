package service

type astServiceStore struct {
	v map[string]astDo
}

func (s *astServiceStore) get(service string) (do astDo, err error) {
	panic("impl")
}

func (s *astServiceStore) store(service string, fun astDo) error {
	panic("impl")
}
