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

func (r *DBRepository) EnableEntity(entity Entity, place Place, entityID, placeID uint, enable bool) error {
	var availabilityDB Availability

	if err := r.db.FirstOrInit(&availabilityDB, Availability{
		Entity:   string(entity),
		EntityID: &entityID,
		Place:    string(place),
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

func (r *DBRepository) FindEntityByPlace(entity Entity, place Place, placeID string) ([]Availability, error) {
	var availabilities []Availability

	if err := r.db.Where("entity = ? AND place = ? AND place_id = ?", entity, place, placeID).Find(&availabilities).Error; err != nil {
		shared.LogError("error getting availabilities", LogDBRepository, "FindEntityByPlace", err, entity, place, placeID)
		return nil, err
	}

	return availabilities, nil
}

func (r *DBRepository) FindPlacesByEntity(entity Entity, entityID uint, place string) ([]Availability, error) {
	var availabilities []Availability

	if err := r.db.Where("entity = ? AND entity_id = ? AND place = ?", entity, entityID, place).Find(&availabilities).Error; err != nil {
		shared.LogError("error getting availabilities", LogDBRepository, "FindPlacesByEntity", err, entity, entityID)
		return nil, err
	}

	return availabilities, nil
}

func (r *DBRepository) Get(entity Entity, place Place, entityID, placeID uint) (Availability, error) {
	var availability Availability

	if err := r.db.Where("entity = ? AND entity_id = ? AND place = ? AND place_id = ?", entity, entityID, place, placeID).First(&availability).Error; err != nil {
		shared.LogError("error getting availability", LogDBRepository, "Get", err, entity, entityID, place, placeID)
		return Availability{}, err
	}

	return availability, nil
}

func (r *DBRepository) Find(entity Entity, place Place, entityID uint) ([]Availability, error) {
	var availabilities []Availability

	if err := r.db.Where("entity = ? AND entity_id = ? AND place = ?", entity, entityID, place).Find(&availabilities).Error; err != nil {
		shared.LogError("error getting availabilities", LogDBRepository, "Find", err, entity, entityID, place)
		return nil, err
	}

	return availabilities, nil
}
