package discount

type service struct {
	repository RepositoryI
}

func NewService(repository RepositoryI) service {
	return service{repository}
}

func (s service) Create(discount *Discount) (*Discount, error) {
	return s.repository.Create(discount)
}

func (s service) Find(filter map[string]string) ([]Discount, error) {
	return s.repository.Find(filter)
}

func (s service) Get(discountID string) (*Discount, error) {
	return s.repository.Get(discountID)
}

func (s service) Update(discount Discount) (*Discount, error) {
	return s.repository.Update(discount)
}

func (s service) Delete(discountID string) (*Discount, error) {
	return s.repository.Delete(discountID)
}

type Service interface {
	Create(discount *Discount) (*Discount, error)
	Find(map[string]string) ([]Discount, error)
	Get(DiscountID string) (*Discount, error)
	Update(discount Discount) (*Discount, error)
	Delete(DiscountID string) (*Discount, error)
}
