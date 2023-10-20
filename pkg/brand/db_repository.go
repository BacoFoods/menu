package brand

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
)

const LogDBRepository = "pkg/brand/db_repository"

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db: db}
}

// Create method for create a new brand in database
func (r *DBRepository) Create(brand *Brand) (*Brand, error) {
	if err := r.db.Save(brand).Error; err != nil {
		shared.LogError("error creating brand", LogDBRepository, "Create", err, brand)
		return nil, err
	}

	return brand, nil
}

// Find method for find brands in database
func (r *DBRepository) Find(filters map[string]string) ([]Brand, error) {
	var brands []Brand

	if err := r.db.Find(&brands, filters).Error; err != nil {
		shared.LogError("error getting brands", LogDBRepository, "Find", err, filters)
		return nil, err
	}

	return brands, nil
}

// Get method for get a brand in database
func (r *DBRepository) Get(brandID string) (*Brand, error) {
	if brandID == "" {
		shared.LogWarn("error getting brand", LogDBRepository, "Get", shared.ErrorIDEmpty)
		return nil, shared.ErrorIDEmpty
	}

	var brand Brand

	if err := r.db.First(&brand, brandID).Error; err != nil {
		shared.LogError("error getting brand", LogDBRepository, "Get", err, brandID)
		return nil, err
	}

	return &brand, nil
}

// Update method for update a brand in database
func (r *DBRepository) Update(brand *Brand) (*Brand, error) {
	var brandDB Brand

	if err := r.db.First(&brandDB, brand.ID).Error; err != nil {
		shared.LogError("error getting brand", LogDBRepository, "Update", err, brand)
		return nil, err
	}

	if err := r.db.Model(&brandDB).Updates(brand).Error; err != nil {
		shared.LogError("error updating brand", LogDBRepository, "Update", err, brand)
		return nil, err
	}

	return &brandDB, nil
}

// Delete method for delete a brand in database
func (r *DBRepository) Delete(brandID string) (*Brand, error) {
	var brand Brand
	if err := r.db.First(&brand, brandID).Error; err != nil {
		shared.LogError("error getting brand", LogDBRepository, "Delete", err, brandID)
		return nil, err
	}

	if err := r.db.Delete(&brand).Error; err != nil {
		shared.LogError("error deleting brand", LogDBRepository, "Delete", err, brand)
		return nil, err
	}

	return &brand, nil
}
