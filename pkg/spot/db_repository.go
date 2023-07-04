package spot

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
)

const LogDBRepository string = "pkg/spot/db_repository"

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db: db}
}

// Create method for create a new spot in database
func (r *DBRepository) Create(spot *Spot) (*Spot, error) {
	if err := r.db.Save(spot).Error; err != nil {
		shared.LogError("error creating spot", LogDBRepository, "Create", err, spot)
		return nil, err
	}
	return spot, nil
}

// Find method for find spots in database
func (r *DBRepository) Find(filters map[string]string) ([]Spot, error) {
	var spots []Spot
	if err := r.db.Find(&spots, filters).Error; err != nil {
		shared.LogError("error getting spots", LogDBRepository, "Find", err, filters)
		return nil, err
	}
	return spots, nil
}

// Get method for get a spot in database
func (r *DBRepository) Get(spotID string) (*Spot, error) {
	var spot Spot
	if err := r.db.First(&spot, spotID).Error; err != nil {
		shared.LogError("error getting spot", LogDBRepository, "Get", err, spotID)
		return nil, err
	}
	return &spot, nil
}

// Update method for update a spot in database
func (r *DBRepository) Update(spot *Spot) (*Spot, error) {
	var spotDB Spot
	if err := r.db.First(&spotDB, spot.ID).Error; err != nil {
		shared.LogError("error getting spot", LogDBRepository, "Update", err, spot)
		return nil, err
	}
	if err := r.db.Model(&spotDB).Updates(spot).Error; err != nil {
		shared.LogError("error updating spot", LogDBRepository, "Update", err, spot)
		return nil, err
	}
	return &spotDB, nil
}

// Delete method for delete a spot in database
func (r *DBRepository) Delete(spotID string) (*Spot, error) {
	var spot Spot
	if err := r.db.First(&spot, spotID).Error; err != nil {
		shared.LogError("error getting spot", LogDBRepository, "Delete", err, spotID)
		return nil, err
	}

	if err := r.db.Delete(&spot).Error; err != nil {
		shared.LogError("error deleting spot", LogDBRepository, "Delete", err, spot)
		return nil, err
	}

	return &spot, nil
}
