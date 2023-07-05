package overriders

type Service interface {
	Find(map[string]string) ([]Overriders, error)
	Get(string) (*Overriders, error)
	Create(*Overriders) (*Overriders, error)
	Update(*Overriders) (*Overriders, error)
	Delete(string) (*Overriders, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) service {
	return service{repository: repository}
}

func (s service) Find(filter map[string]string) ([]Overriders, error) {
	return s.repository.Find(filter)
}

func (s service) Get(overridersID string) (*Overriders, error) {
	return s.repository.Get(overridersID)
}

func (s service) Create(overriders *Overriders) (*Overriders, error) {
	return s.repository.Create(overriders)
}

func (s service) Update(overriders *Overriders) (*Overriders, error) {
	return s.repository.Update(overriders)
}

func (s service) Delete(overridersID string) (*Overriders, error) {
	return s.repository.Delete(overridersID)
}
