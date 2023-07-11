package availability

import (
	"fmt"
	channelPkg "github.com/BacoFoods/menu/pkg/channel"
	storePkg "github.com/BacoFoods/menu/pkg/store"
)

type Service interface {
	EnableEntity(entity Entity, place Place, entityID, placeID uint, enable bool) error
	FindEntities() []Entity
	FindPlaces() []Place
	Get(entity Entity, place Place, entityID, placeID uint) (any, error)
	Find(entity Entity, place Place, entityID uint) ([]any, error)
}

type service struct {
	repository Repository
	store      storePkg.Repository
	channel    channelPkg.Repository
}

func NewService(repository Repository, store storePkg.Repository, channel channelPkg.Repository) service {
	return service{repository, store, channel}
}

func (s service) EnableEntity(entity Entity, place Place, entityID, placeID uint, enable bool) error {
	return s.repository.EnableEntity(entity, place, entityID, placeID, enable)
}

func (s service) FindEntities() []Entity {
	return []Entity{EntityMenu, EntityCategory}
}

func (s service) FindPlaces() []Place {
	return []Place{PlaceStore, PlaceChannel}
}

func (s service) Get(entity Entity, place Place, entityID, placeID uint) (any, error) {
	switch place {
	case PlaceStore:
		store, err := s.store.Get(fmt.Sprintf("%v", placeID))
		if err != nil {
			return nil, err
		}
		return store, nil
	default:
		return nil, fmt.Errorf("place %s not supported", place)
	}
}

func (s service) Find(entity Entity, place Place, entityID uint) ([]any, error) {
	availabilities, err := s.repository.Find(entity, place, entityID)
	if err != nil {
		return nil, err
	}

	placeIDs := make([]string, 0)
	availabilityStates := make(map[uint]bool)
	for _, availability := range availabilities {
		availabilityStates[*availability.PlaceID] = availability.Enable
		placeIDs = append(placeIDs, fmt.Sprintf("%v", *availability.PlaceID))
	}

	var values []any
	switch place {
	case PlaceStore:
		stores, err := s.store.FindByIDs(placeIDs)
		if err != nil {
			return nil, err
		}

		for _, element := range stores {
			element.Enabled = availabilityStates[element.ID]
			values = append(values, element)
		}

		return values, nil
	case PlaceChannel:
		channels, err := s.channel.FindByIDs(placeIDs)
		if err != nil {
			return nil, err
		}

		for _, element := range channels {
			element.Enabled = availabilityStates[element.ID]
			values = append(values, element)
		}

		return values, nil
	default:
		return nil, fmt.Errorf("place %s not supported", place)
	}
}
