package currency

type service struct {
	repository Repository
}

func NewService(repository Repository) service {
	return service{repository}
}

func (s service) Create(currency *Currency) (*Currency, error) {
	return s.repository.Create(currency)
}

func (s service) Find(query map[string]string) ([]Currency, error) {
	return s.repository.Find(query)
}

func (s service) Get(currencyID string) (*Currency, error) {
	return s.repository.Get(currencyID)
}

func (s service) Update(currency Currency) (*Currency, error) {
	return s.repository.Update(currency)
}

func (s service) Delete(currencyID string) (*Currency, error) {
	return s.repository.Delete(currencyID)
}

type Service interface {
	Create(*Currency) (*Currency, error)
	Find(map[string]string) ([]Currency, error)
	Get(string) (*Currency, error)
	Update(Currency) (*Currency, error)
	Delete(string) (*Currency, error)
}
