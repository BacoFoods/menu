package shift

type Service interface {
	OpenShift() error
	CloseShift() error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) service {
	return service{repository}
}

func (s service) OpenShift() error {
	return s.repository.OpenShift()
}

func (s service) CloseShift() error {
	return s.repository.CloseShift()
}
