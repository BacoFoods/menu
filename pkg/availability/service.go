package availability

import (
	"fmt"
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
}

func NewService(repository Repository, store storePkg.Repository) service {
	return service{repository, store}
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
	for _, availability := range availabilities {
		placeIDs = append(placeIDs, fmt.Sprintf("%v", *availability.PlaceID))
	}

	var values []any
	switch place {
	case PlaceStore:
		stores, err := s.store.FindByStores(placeIDs)
		if err != nil {
			return nil, err
		}

		for _, element := range stores {
			values = append(values, element)
		}

		return values, nil

	default:
		return nil, fmt.Errorf("place %s not supported", place)

	}
}
