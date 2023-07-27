package surcharge

type Service interface {
	Find(query map[string]string) ([]Surcharge, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) service {
	return service{repository}
}

func (s service) Find(query map[string]string) ([]Surcharge, error) {
	return s.repository.Find(query)
}
