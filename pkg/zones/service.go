package zones

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/shared"
)

const (
	LogService = "pkg/zone/service"
)

type Service interface {
	Find(filters map[string]any) ([]Zone, error)
	Get(zoneID string) (*Zone, error)
	Create(zone *Zone) (*Zone, error)
	Update(zonID string, zone *Zone) (*Zone, error)
	Delete(zoneID string) error
	AddTablesToZone(zoneID string, tables []uint) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) service {
	return service{repository}
}

func (s service) Find(filters map[string]any) ([]Zone, error) {
	return s.repository.Find(filters)
}

func (s service) Get(zoneID string) (*Zone, error) {
	return s.repository.GetZone(zoneID)
}

func (s service) Create(zone *Zone) (*Zone, error) {
	return s.repository.Create(zone)
}

func (s service) Update(zoneID string, zone *Zone) (*Zone, error) {
	return s.repository.Update(zoneID, zone)
}

func (s service) Delete(zoneID string) error {
	return s.repository.Delete(zoneID)
}

func (s service) AddTablesToZone(zoneID string, tables []uint) error {
	zone, err := s.repository.GetZone(zoneID)
	if err != nil {
		return err
	}

	if zone == nil {
		err := fmt.Errorf(ErrorZoneNotFound)
		shared.LogError("error finding zone", LogService, "AddTablesToZone", err, zoneID)
		return err
	}

	return s.repository.AddTablesToZone(zoneID, tables)
}
