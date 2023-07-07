package availability

import (
	"strconv"
)

type Service interface {
	EnableEntity(string, string, string, string, bool) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) service {
	return service{repository}
}

func (s service) EnableEntity(entity, entityID, place, placeID string, enable bool) error {
	entityId, err := strconv.ParseUint(entityID, 10, 64)
	if err != nil {
		return err
	}

	placeId, err := strconv.ParseUint(placeID, 10, 64)
	if err != nil {
		return err
	}

	uintEntityID := uint(entityId)
	uintPlaceID := uint(placeId)

	if err := s.repository.EnableEntity(entity, place, uintEntityID, uintPlaceID, enable); err != nil {
		return err
	}

	return nil
}
