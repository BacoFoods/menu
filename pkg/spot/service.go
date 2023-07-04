package spot

// Service to handle business logic for the spot service
type Service interface {
	Find(map[string]string) ([]Spot, error)
	Get(string) (*Spot, error)
	Create(*Spot) (*Spot, error)
	Update(*Spot) (*Spot, error)
	Delete(string) (*Spot, error)
}

// service is the default implementation of the Service interface for spot.
type service struct {
	repository Repository
}

// NewService creates a new instance of the service for spot, using the provided repository implementation.
func NewService(repository Repository) service {
	return service{repository}
}

// Find returns a list of spot objects filtering by query map.
func (s service) Find(filter map[string]string) ([]Spot, error) {
	return s.repository.Find(filter)
}

// Get returns a single spot object by ID.
func (s service) Get(spotID string) (*Spot, error) {
	return s.repository.Get(spotID)
}

// Create creates a new spot object.
func (s service) Create(spot *Spot) (*Spot, error) {
	return s.repository.Create(spot)
}

// Update updates an existing spot object only the fields that are present in the provided object.
// this method doesn't create new register if the provided id doesn't exist.
func (s service) Update(spot *Spot) (*Spot, error) {
	return s.repository.Update(spot)
}

// Delete deletes an existing spot object.
func (s service) Delete(spotID string) (*Spot, error) {
	return s.repository.Delete(spotID)
}
