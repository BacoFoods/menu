package availability

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
)

const LogDBRepository string = "pkg/availability/db_repository"

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db}
}

func (r *DBRepository) EnableEntity(entity, place string, entityID, placeID uint, enable bool) error {
	var availabilityDB Availability

	if err := r.db.FirstOrInit(&availabilityDB, Availability{
		Entity:   entity,
		EntityID: &entityID,
		Place:    place,
		PlaceID:  &placeID,
	}).Error; err != nil {
		shared.LogError("error getting availability", LogDBRepository, "Update", err, entity, entityID, place, placeID)
		return err
	}

	availabilityDB.Enable = enable
	if err := r.db.Save(&availabilityDB).Error; err != nil {
		shared.LogError("error updating availability", LogDBRepository, "Update", err, availabilityDB)
		return err
	}

	return nil
}

func (r *DBRepository) FindEntityByPlace(entity Entity, place, placeID string) ([]Availability, error) {
	var availabilities []Availability

	if err := r.db.Where("entity = ? AND place = ? AND place_id = ?", entity, place, placeID).Find(&availabilities).Error; err != nil {
		shared.LogError("error getting availabilities", LogDBRepository, "FindEntityByPlace", err, entity, place, placeID)
		return nil, err
	}

	return availabilities, nil
}
