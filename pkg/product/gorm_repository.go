package product

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
)

const LogGormRepository string = "pkg/product/gorm_repository"

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

// Create method for create a new product in database
func (r *GormRepository) Create(product *Product) (*Product, error) {
	if err := r.db.Save(product).Error; err != nil {
		shared.LogError("error creating product", LogGormRepository, "Create", err, product)
		return nil, err
	}
	return product, nil
}

// Find method for find products in database
func (r *GormRepository) Find(filters map[string]string) ([]Product, error) {
	var products []Product
	if err := r.db.Find(&products, filters).Error; err != nil {
		shared.LogError("error getting products", LogGormRepository, "Find", err, filters)
		return nil, err
	}
	return products, nil
}

// Get method for get a product in database
func (r *GormRepository) Get(productID string) (*Product, error) {
	var product Product
	if err := r.db.First(&product, productID).Error; err != nil {
		shared.LogError("error getting product", LogGormRepository, "Get", err, productID)
		return nil, err
	}
	return &product, nil
}

// Update method for update a product in database
func (r *GormRepository) Update(product *Product) (*Product, error) {
	var productDB Product
	if err := r.db.First(&productDB, product.ID).Error; err != nil {
		shared.LogError("error getting product", LogGormRepository, "Update", err, product)
		return nil, err
	}
	if err := r.db.Model(&productDB).Updates(product).Error; err != nil {
		shared.LogError("error updating product", LogGormRepository, "Update", err, product)
		return nil, err
	}
	return &productDB, nil
}

// Delete method for delete a product in database
func (r *GormRepository) Delete(productID string) (*Product, error) {
	var product Product
	if err := r.db.First(&product, productID).Error; err != nil {
		shared.LogError("error getting product", LogGormRepository, "Delete", err, productID)
		return nil, err
	}

	if err := r.db.Delete(&product).Error; err != nil {
		shared.LogError("error deleting product", LogGormRepository, "Delete", err, productID)
		return nil, err
	}

	return &product, nil
}
