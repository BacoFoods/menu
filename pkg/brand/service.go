package brand

import (
	channels "github.com/BacoFoods/menu/pkg/channel"
	"github.com/BacoFoods/menu/pkg/shared"
)

const LogService = "pkg/brand/service.go"

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
	channel    channels.Repository
}

// NewService creates a new instance of the service for Brand, using the provided repository implementation.
func NewService(repository Repository, channel channels.Repository) service {
	return service{repository, channel}
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
	newBrand, err := s.repository.Create(brand)
	if err != nil {
		return nil, err
	}

	// Create default channels for brand
	if _, err := s.channel.Create(&channels.Channel{
		Name:      "Salon",
		ShortName: "salon",
		BrandID:   &newBrand.ID,
		Enabled:   false,
	}); err != nil {
		shared.LogWarn("error creating default channel Salon for brand", LogService, "Create", err, newBrand.ID)
	}

	if _, err := s.channel.Create(&channels.Channel{
		Name:      "Pick Up",
		ShortName: "pick up",
		BrandID:   &newBrand.ID,
		Enabled:   false,
	}); err != nil {
		shared.LogWarn("error creating default channel PickUp for brand", LogService, "Create", err, newBrand.ID)
	}

	if _, err := s.channel.Create(&channels.Channel{
		Name:      "Domicilio Propio",
		ShortName: "domicilio",
		BrandID:   &newBrand.ID,
		Enabled:   false,
	}); err != nil {
		shared.LogWarn("error creating default channel Domicilio for brand", LogService, "Create", err, newBrand.ID)
	}

	return newBrand, nil
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
