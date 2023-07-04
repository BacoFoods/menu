package channel

// Service is the interface that provides channel methods.
type Service interface {
	Find(map[string]string) ([]Channel, error)
	Get(string) (*Channel, error)
	Create(*Channel) (*Channel, error)
	Update(*Channel) (*Channel, error)
	Delete(string) (*Channel, error)
}

// service is the default implementation of the Service interface for channel.
type service struct {
	repository Repository
}

// NewService creates a new instance of the service for channel, using the provided repository implementation.
func NewService(repository Repository) service {
	return service{repository}
}

// Find returns a list of channel objects filtering by query map.
func (s service) Find(filter map[string]string) ([]Channel, error) {
	return s.repository.Find(filter)
}

// Get returns a single channel object by ID.
func (s service) Get(channelID string) (*Channel, error) {
	return s.repository.Get(channelID)
}

// Create creates a new channel object.
func (s service) Create(channel *Channel) (*Channel, error) {
	return s.repository.Create(channel)
}

// Update updates an existing channel object only the fields that are present in the provided object.
// this method doesn't create new register if the provided id doesn't exist.
func (s service) Update(channel *Channel) (*Channel, error) {
	return s.repository.Update(channel)
}

// Delete deletes an existing channel object.
func (s service) Delete(channelID string) (*Channel, error) {
	return s.repository.Delete(channelID)
}
