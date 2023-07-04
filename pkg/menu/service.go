package menu

// Service is the interface that provides menu methods, used for dependency injection.
type Service interface {
	Find(map[string]string) ([]Menu, error)
	Get(string) (*Menu, error)
	Create(*Menu) (*Menu, error)
	Update(*Menu) (*Menu, error)
	Delete(string) (*Menu, error)
}

// service is the default implementation of the Service interface for menu.
type service struct {
	repository Repository
}

// NewService creates a new instance of the service for menu, using the provided repository implementation.
func NewService(repository Repository) service {
	return service{repository}
}

// Find returns a list of menu objects filtering by query map.
func (s service) Find(filter map[string]string) ([]Menu, error) {
	return s.repository.Find(filter)
}

// Get returns a single menu object by ID.
func (s service) Get(menuID string) (*Menu, error) {
	return s.repository.Get(menuID)
}

// Create creates a new menu object.
func (s service) Create(menu *Menu) (*Menu, error) {
	return s.repository.Create(menu)
}

// Update updates an existing menu object only the fields that are present in the provided object.
// this method doesn't create new register if the provided id doesn't exist.
func (s service) Update(menu *Menu) (*Menu, error) {
	return s.repository.Update(menu)
}

// Delete deletes an existing menu object.
func (s service) Delete(menuID string) (*Menu, error) {
	return s.repository.Delete(menuID)
}
