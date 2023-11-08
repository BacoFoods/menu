package store

import (
	"fmt"
	channelPkg "github.com/BacoFoods/menu/pkg/channel"
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/BacoFoods/menu/pkg/tables"
)

const (
	LogService = "pkg/store/service"
)

// Service to handle business logic for the store service
type Service interface {
	Find(map[string]string) ([]Store, error)
	Get(string) (*Store, error)
	Create(*Store) (*Store, error)
	Update(*Store) (*Store, error)
	Delete(string) (*Store, error)
	Enable(string) (*Store, error)
	AddChannel(storeID, channelID string) (*Store, error)

	FindZonesByStore(storeID string) ([]tables.Zone, error)
	GetZoneByStore(storeID, zoneID string) (*tables.Zone, error)
}

// service is the default implementation of the Service interface for store.
type service struct {
	repository Repository
	channel    channelPkg.Repository
}

// NewService creates a new instance of the service for store, using the provided repository implementation.
func NewService(repository Repository, channel channelPkg.Repository) service {
	return service{repository, channel}
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
	channels, err := s.channel.Find(map[string]string{"brand_id": fmt.Sprintf("%v", *store.BrandID)})
	if err != nil {
		shared.LogError("error getting channels by brand id", LogService, "Create", err, store.BrandID)
		return nil, fmt.Errorf(ErrorStoreGettingChannels)
	}

	store.Channels = channels

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

// Enable enables an existing store object.
func (s service) Enable(storeID string) (*Store, error) {
	return s.repository.Enable(storeID)
}

// AddChannel adds a channel to a store.
func (s service) AddChannel(storeID, channelID string) (*Store, error) {
	channel, err := s.channel.Get(channelID)
	if err != nil {
		return nil, err
	}

	return s.repository.AddChannel(storeID, channel)
}

// FindZonesByStore returns a list of zones by store id.
func (s service) FindZonesByStore(storeID string) ([]tables.Zone, error) {
	return s.repository.FindZonesByStore(storeID)
}

// GetZoneByStore returns a single zone by store id and zone id.
func (s service) GetZoneByStore(storeID, zoneID string) (*tables.Zone, error) {
	return s.repository.GetZoneByStore(storeID, zoneID)
}
