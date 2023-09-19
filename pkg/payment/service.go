package payment

type Service interface {
	Get(paymentID string) (*Payment, error)
	Find(filter map[string]any) ([]Payment, error)
	Create(payment *Payment) (*Payment, error)
	Update(payment *Payment) (*Payment, error)
	Delete(paymentID string) (*Payment, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) service {
	return service{repository: repository}
}

func (s service) Get(paymentID string) (*Payment, error) {
	return s.repository.Get(paymentID)
}

func (s service) Find(filter map[string]any) ([]Payment, error) {
	return s.repository.Find(filter)
}

func (s service) Create(payment *Payment) (*Payment, error) {
	return s.repository.Create(payment)
}

func (s service) Update(payment *Payment) (*Payment, error) {
	return s.repository.Update(payment)
}

func (s service) Delete(paymentID string) (*Payment, error) {
	return s.repository.Delete(paymentID)
}
