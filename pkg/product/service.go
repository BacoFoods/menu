package product

type Service interface {
	Find(map[string]string) ([]Product, error)
	Get(productID []string) ([]Product, error)
	Create(*Product) (*Product, error)
	Update(*Product) (*Product, error)
	Delete(string) (*Product, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) service {
	return service{repository}
}

func (s service) Find(filter map[string]string) ([]Product, error) {
	return s.repository.Find(filter)
}

func (s service) Get(productID []string) ([]Product, error) {
	return s.repository.Get(productID)
}

func (s service) Create(product *Product) (*Product, error) {
	return s.repository.Create(product)
}

func (s service) Update(product *Product) (*Product, error) {
	return s.repository.Update(product)
}

func (s service) Delete(productID string) (*Product, error) {
	return s.repository.Delete(productID)
}
