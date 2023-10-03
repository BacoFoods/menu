package payment

type Service interface {
	Get(paymentID string) (*Payment, error)
	Find(filter map[string]any) ([]Payment, error)
	Create(payment *Payment) (*Payment, error)
	Update(payment *Payment) (*Payment, error)
	Delete(paymentID string) (*Payment, error)

	FindPaymentMethods(filter map[string]any) ([]PaymentMethod, error)
	GetPaymentMethod(string) (*PaymentMethod, error)
	CreatePaymentMethod(*PaymentMethod) (*PaymentMethod, error)
	UpdatePaymentMethod(*PaymentMethod) (*PaymentMethod, error)
	DeletePaymentMethod(string) (*PaymentMethod, error)
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

func (s service) FindPaymentMethods(filter map[string]any) ([]PaymentMethod, error) {
	return s.repository.FindPaymentMethods(filter)
}

func (s service) GetPaymentMethod(paymentMethodID string) (*PaymentMethod, error) {
	return s.repository.GetPaymentMethod(paymentMethodID)
}

func (s service) CreatePaymentMethod(paymentMethod *PaymentMethod) (*PaymentMethod, error) {
	return s.repository.CreatePaymentMethod(paymentMethod)
}

func (s service) UpdatePaymentMethod(paymentMethod *PaymentMethod) (*PaymentMethod, error) {
	return s.repository.UpdatePaymentMethod(paymentMethod)
}

func (s service) DeletePaymentMethod(paymentMethodID string) (*PaymentMethod, error) {
	return s.repository.DeletePaymentMethod(paymentMethodID)
}
