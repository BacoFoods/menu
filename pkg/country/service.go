package country

type service struct {
	repository Repository
}

func NewService(repository Repository) service {
	return service{repository}
}

func (s service) Create(country *Country) (*Country, error) {
	return s.repository.Create(country)
}

func (s service) Find(query map[string]string) ([]Country, error) {
	return s.repository.Find(query)
}

func (s service) Get(countryID string) (*Country, error) {
	return s.repository.Get(countryID)
}

func (s service) Update(country Country) (*Country, error) {
	return s.repository.Update(country)
}

func (s service) Delete(countryID string) (*Country, error) {
	return s.repository.Delete(countryID)
}

type Service interface {
	Create(*Country) (*Country, error)
	Find(map[string]string) ([]Country, error)
	Get(string) (*Country, error)
	Update(Country) (*Country, error)
	Delete(string) (*Country, error)
}
