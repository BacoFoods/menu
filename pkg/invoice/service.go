package invoice

const LogService = "pkg/invoice/service"

type Service interface {
	Get(string) (*Invoice, error)
}

type service struct {
	repository Repository
	invoice Repository
}

func NewService(repository Repository, invoice Repository) service {
	return service{repository, invoice}
}

// Get returns a single Invoice object by ID.
func (s service) Get(invoiceID string) (*Invoice, error) {
	return s.repository.Get(invoiceID)
}