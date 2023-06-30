package taxes

type service struct {
	repository Repository
}

func NewService(repository Repository) service {
	return service{repository}
}

func (s service) Create(tax *Tax) (*Tax, error) {
	return s.repository.Create(tax)
}

func (s service) Find(query map[string]string) ([]Tax, error) {
	return s.repository.Find(query)
}

func (s service) Get(taxID string) (*Tax, error) {
	return s.repository.Get(taxID)
}

func (s service) Update(tax Tax) (*Tax, error) {
	return s.repository.Update(tax)
}

func (s service) Delete(taxID string) (*Tax, error) {
	return s.repository.Delete(taxID)
}

type Service interface {
	Create(tax *Tax) (*Tax, error)
	Find(map[string]string) ([]Tax, error)
	Get(TaxID string) (*Tax, error)
	Update(tax Tax) (*Tax, error)
	Delete(TaxID string) (*Tax, error)
}
