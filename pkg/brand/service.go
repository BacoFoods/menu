package brand

type Service interface {
	Find(map[string]string) ([]Brand, error)
	Get(string) (*Brand, error)
	Create(*Brand) (*Brand, error)
	Update(*Brand) (*Brand, error)
	Delete(string) (*Brand, error)
}

// service is the default implementation of the Service interface for Brand.
type service struct {
	repository Repository
}

// NewService creates a new instance of the service for Brand, using the provided repository implementation.
func NewService(repository Repository) service {
	return service{repository: repository}
}

// Find returns a list of Brand objects filtering by query map.
func (s service) Find(query map[string]string) ([]Brand, error) {
	return s.repository.Find(query)
}

// Get returns a single Brand object by ID.
func (s service) Get(brandID string) (*Brand, error) {
	return s.repository.Get(brandID)
}

// Create creates a new Brand object.
func (s service) Create(brand *Brand) (*Brand, error) {
	return s.repository.Create(brand)
}

// Update updates an existing Brand object only the fields that are present in the provided object.
// this method doesn't create new register if the provided id doesn't exist.
func (s service) Update(brand *Brand) (*Brand, error) {
	return s.repository.Update(brand)
}

// Delete deletes an existing Brand object.
func (s service) Delete(brandID string) (*Brand, error) {
	return s.repository.Delete(brandID)
}
