package cashier

type Service interface {
	Open() error
	Close() error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) service {
	return service{repository}
}

func (s service) Open() error {
	return s.repository.Open()
}

func (s service) Close() error {
	return s.repository.Close()
}
