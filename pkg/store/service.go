package store

// Service to handle business logic for the store service
type Service interface {
	Find(map[string]string) ([]Store, error)
	Get(string) (*Store, error)
	Create(*Store) (*Store, error)
	Update(*Store) (*Store, error)
	Delete(string) (*Store, error)
}

// service is the default implementation of the Service interface for store.
type service struct {
	repository Repository
}

// NewService creates a new instance of the service for store, using the provided repository implementation.
func NewService(repository Repository) service {
	return service{repository}
}

// Find returns a list of store objects filtering by query map.
func (s service) Find(filter map[string]string) ([]Store, error) {
	return s.repository.Find(filter)
}

// Get returns a single store object by ID.
func (s service) Get(storeID string) (*Store, error) {
	return s.repository.Get(storeID)
}

// Create creates a new store object.
func (s service) Create(store *Store) (*Store, error) {
	return s.repository.Create(store)
}

// Update updates an existing store object only the fields that are present in the provided object.
// this method doesn't create new register if the provided id doesn't exist.
func (s service) Update(store *Store) (*Store, error) {
	return s.repository.Update(store)
}

// Delete deletes an existing store object.
func (s service) Delete(storeID string) (*Store, error) {
	return s.repository.Delete(storeID)
}
