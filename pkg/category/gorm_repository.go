package category

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
)

const LogGormRepository string = "pkg/category/repository"

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

// Create method for create a new category in database
func (r *GormRepository) Create(category *Category) (*Category, error) {
	if err := r.db.Save(category).Error; err != nil {
		shared.LogError("error creating category", LogGormRepository, "Create", err, category)
		return nil, err
	}
	return category, nil
}

// Find method for find categories in database
func (r *GormRepository) Find(filters map[string]string) ([]Category, error) {
	var categories []Category
	if err := r.db.Find(&categories, filters).Error; err != nil {
		shared.LogError("error getting categories", LogGormRepository, "Find", err, filters)
		return nil, err
	}
	return categories, nil
}

// Get method for get a category in database
func (r *GormRepository) Get(categoryID string) (*Category, error) {
	var category Category
	if err := r.db.First(&category, categoryID).Error; err != nil {
		shared.LogError("error getting category", LogGormRepository, "Get", err, categoryID)
		return nil, err
	}
	return &category, nil
}

// Update method for update a category in database
func (r *GormRepository) Update(category *Category) (*Category, error) {
	var categoryDB Category
	if err := r.db.First(&categoryDB, category.ID).Error; err != nil {
		shared.LogError("error getting category", LogGormRepository, "Update", err, category)
		return nil, err
	}
	if err := r.db.Model(&categoryDB).Updates(category).Error; err != nil {
		shared.LogError("error updating category", LogGormRepository, "Update", err, category)
		return nil, err
	}
	return &categoryDB, nil
}

// Delete method for delete a category in database
func (r *GormRepository) Delete(categoryID string) (*Category, error) {
	var category Category
	if err := r.db.First(&category, categoryID).Error; err != nil {
		shared.LogError("error getting category", LogGormRepository, "Delete", err, categoryID)
		return nil, err
	}
	if err := r.db.Delete(&category).Error; err != nil {
		shared.LogError("error deleting category", LogGormRepository, "Delete", err, category)
		return nil, err
	}
	return &category, nil
}
