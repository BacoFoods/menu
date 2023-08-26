package category

import (
	"github.com/BacoFoods/menu/pkg/product"
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const LogDBRepository string = "pkg/category/db_repository"

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db: db}
}

// Create method for create a new category in database
func (r *DBRepository) Create(category *Category) (*Category, error) {
	if err := r.db.Preload(clause.Associations).Save(category).Error; err != nil {
		shared.LogError("error creating category", LogDBRepository, "Create", err, category)
		return nil, err
	}
	return category, nil
}

// Find method for find categories in database
func (r *DBRepository) Find(filters map[string]string) ([]Category, error) {
	var categories []Category
	if err := r.db.Preload(clause.Associations).Find(&categories, filters).Error; err != nil {
		shared.LogError("error getting categories", LogDBRepository, "Find", err, filters)
		return nil, err
	}
	return categories, nil
}

// Get method for get a category in database
func (r *DBRepository) Get(categoryID string) (*Category, error) {
	var category Category
	if err := r.db.Preload(clause.Associations).First(&category, categoryID).Error; err != nil {
		shared.LogError("error getting category", LogDBRepository, "Get", err, categoryID)
		return nil, err
	}
	return &category, nil
}

// Update method for update a category in database
func (r *DBRepository) Update(category *Category) (*Category, error) {
	var categoryDB Category
	if err := r.db.Preload(clause.Associations).First(&categoryDB, category.ID).Error; err != nil {
		shared.LogError("error getting category", LogDBRepository, "Update", err, category)
		return nil, err
	}
	if err := r.db.Model(&categoryDB).Updates(category).Error; err != nil {
		shared.LogError("error updating category", LogDBRepository, "Update", err, category)
		return nil, err
	}
	return &categoryDB, nil
}

// Delete method for delete a category in database
func (r *DBRepository) Delete(categoryID string) (*Category, error) {
	var category Category
	if err := r.db.First(&category, categoryID).Error; err != nil {
		shared.LogError("error getting category", LogDBRepository, "Delete", err, categoryID)
		return nil, err
	}
	if err := r.db.Delete(&category).Error; err != nil {
		shared.LogError("error deleting category", LogDBRepository, "Delete", err, category)
		return nil, err
	}
	return &category, nil
}

// GetMenusByCategory method for get menus by category in database
func (r *DBRepository) GetMenusByCategory(categoryID string) ([]MenusCategory, error) {
	var menusCategory []MenusCategory
	if err := r.db.Debug().Table("menus").
		Select("menus.id id, menus.name name, menus.enable enable").
		Joins("left join menus_categories mc on menus.id = mc.menu_id").
		Where("mc.category_id = ?", categoryID).
		Find(&menusCategory).Error; err != nil {
		shared.LogError("error getting menus by category", LogDBRepository, "GetMenusByCategory", err, categoryID)
		return nil, err
	}

	return menusCategory, nil
}

// AddProduct method for add product to category in database
func (r *DBRepository) AddProduct(categoryID, productID uint) (*Category, error) {
	var category Category
	if err := r.db.First(&category, categoryID).Error; err != nil {
		shared.LogError("error getting category", LogDBRepository, "AddProduct", err, categoryID)
		return nil, err
	}

	var productDB product.Product
	if err := r.db.First(&productDB, productID).Error; err != nil {
		shared.LogError("error getting product", LogDBRepository, "AddProduct", err, productID)
		return nil, err
	}

	if err := r.db.Model(&category).Association("Products").Append(&productDB); err != nil {
		shared.LogError("error adding products to category", LogDBRepository, "AddProduct", err, category)
		return nil, err
	}

	return &category, nil
}

// RemoveProduct method for remove product to category in database
func (r *DBRepository) RemoveProduct(categoryID, productID uint) (*Category, error) {
	var category Category
	if err := r.db.First(&category, categoryID).Error; err != nil {
		shared.LogError("error getting category", LogDBRepository, "RemoveProduct", err, categoryID)
		return nil, err
	}

	if err := r.db.Model(&category).Association("Products").Delete(product.Product{ID: productID}); err != nil {
		shared.LogError("error removing products to category", LogDBRepository, "RemoveProduct", err, category)
		return nil, err
	}

	return &category, nil
}
