package overriders

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
)

const LogDBRepository string = "pkg/overriders/db_repository"

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db: db}
}

// Create method for create a new overrider in database
func (r *DBRepository) Create(overrider *Overriders) (*Overriders, error) {
	if err := r.db.Save(overrider).Error; err != nil {
		shared.LogError("error creating overrider", LogDBRepository, "Create", err, overrider)
		return nil, err
	}
	return overrider, nil
}

// Find method for find overriders in database
func (r *DBRepository) Find(filters map[string]string) ([]Overriders, error) {
	var overriders []Overriders
	if err := r.db.Find(&overriders, filters).Error; err != nil {
		shared.LogError("error getting overriders", LogDBRepository, "Find", err, filters)
		return nil, err
	}
	return overriders, nil
}

// Get method for get an overrider in database
func (r *DBRepository) Get(overriderID string) (*Overriders, error) {
	var overrider Overriders
	if err := r.db.First(&overrider, overriderID).Error; err != nil {
		shared.LogError("error getting overrider", LogDBRepository, "Get", err, overriderID)
		return nil, err
	}
	return &overrider, nil
}

// Update method for update an overrider in database
func (r *DBRepository) Update(overrider *Overriders) (*Overriders, error) {
	var overriderDB Overriders
	if err := r.db.First(&overriderDB, overrider.ID).Error; err != nil {
		shared.LogError("error getting overrider", LogDBRepository, "Update", err, overrider)
		return nil, err
	}
	if err := r.db.Model(&overriderDB).Updates(overrider).Error; err != nil {
		shared.LogError("error updating overrider", LogDBRepository, "Update", err, overrider)
		return nil, err
	}
	return &overriderDB, nil
}

// Delete method for delete an overrider in database
func (r *DBRepository) Delete(overriderID string) (*Overriders, error) {
	var overrider Overriders
	if err := r.db.First(&overrider, overriderID).Error; err != nil {
		shared.LogError("error getting overrider", LogDBRepository, "Delete", err, overriderID)
		return nil, err
	}
	if err := r.db.Delete(&overrider).Error; err != nil {
		shared.LogError("error deleting overrider", LogDBRepository, "Delete", err, overriderID)
		return nil, err
	}
	return &overrider, nil
}

// FindByPlace method for find overriders in database
func (r *DBRepository) FindByPlace(place, placeID string) ([]Overriders, error) {
	var overriders []Overriders
	if err := r.db.Where("place = ? AND place_id = ?", place, placeID).Find(&overriders).Error; err != nil {
		shared.LogError("error getting overriders", LogDBRepository, "FindByPlace", err, place, placeID)
		return nil, err
	}
	return overriders, nil
}
