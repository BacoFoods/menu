package product

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const LogDBRepository string = "pkg/product/db_repository"

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db: db}
}

// Create method for create a new product in database
func (r *DBRepository) Create(product *Product) (*Product, error) {
	if err := r.db.Save(product).Error; err != nil {
		shared.LogError("error creating product", LogDBRepository, "Create", err, product)
		return nil, err
	}
	return product, nil
}

// Find method for find products in database
func (r *DBRepository) Find(filters map[string]string) ([]Product, error) {
	var products []Product
	if err := r.db.Preload(clause.Associations).
		Preload("Modifiers.Products").
		Find(&products, filters).Error; err != nil {
		shared.LogError("error getting products", LogDBRepository, "Find", err, filters)
		return nil, err
	}
	return products, nil
}

// Get method for get a product in database
func (r *DBRepository) Get(productID string) (*Product, error) {
	var product Product
	if err := r.db.Preload(clause.Associations).
		Preload("Modifiers.Products").
		First(&product, productID).Error; err != nil {
		shared.LogError("error getting product", LogDBRepository, "Get", err, productID)
		return nil, err
	}
	return &product, nil
}

// GetByIDs method for get products by ids in database
func (r *DBRepository) GetByIDs(productIDs []string) ([]Product, error) {
	var products []Product
	if err := r.db.Find(&products, productIDs).Error; err != nil {
		shared.LogError("error getting products", LogDBRepository, "GetByIDs", err, productIDs)
		return nil, err
	}
	return products, nil
}

// Update method for update a product in database
func (r *DBRepository) Update(product *Product) (*Product, error) {
	var productDB Product
	if err := r.db.First(&productDB, product.ID).Error; err != nil {
		shared.LogError("error getting product", LogDBRepository, "Update", err, product)
		return nil, err
	}
	if err := r.db.Model(&productDB).Updates(product).Error; err != nil {
		shared.LogError("error updating product", LogDBRepository, "Update", err, product)
		return nil, err
	}
	return &productDB, nil
}

// Delete method for delete a product in database
func (r *DBRepository) Delete(productID string) (*Product, error) {
	var product Product
	if err := r.db.First(&product, productID).Error; err != nil {
		shared.LogError("error getting product", LogDBRepository, "Delete", err, productID)
		return nil, err
	}

	if err := r.db.Delete(&product).Error; err != nil {
		shared.LogError("error deleting product", LogDBRepository, "Delete", err, productID)
		return nil, err
	}

	return &product, nil
}

// AddModifier method for add a modifier in product
func (r *DBRepository) AddModifier(product *Product, modifier *Modifier) (*Product, error) {
	if err := r.db.Model(&product).Association("Modifiers").Append(modifier); err != nil {
		shared.LogError("error adding modifier to product", LogDBRepository, "AddModifier", err, product, modifier)
		return nil, err
	}

	return product, nil
}

// RemoveModifier method for remove a modifier in product
func (r *DBRepository) RemoveModifier(product *Product, modifier *Modifier) (*Product, error) {
	if err := r.db.Model(&product).Association("Modifiers").Delete(modifier); err != nil {
		shared.LogError("error removing modifier to product", LogDBRepository, "RemoveModifier", err, product, modifier)
		return nil, err
	}

	return product, nil
}

// ModifierCreate method for create a new modifier in database
func (r *DBRepository) ModifierCreate(modifier *Modifier) (*Modifier, error) {
	if err := r.db.Save(modifier).Error; err != nil {
		shared.LogError("error creating modifier", LogDBRepository, "CreateModifier", err, modifier)
		return nil, err
	}

	return modifier, nil
}

// ModifierGet method for get a modifier in database
func (r *DBRepository) ModifierGet(modifierID string) (*Modifier, error) {
	var modifier Modifier
	if err := r.db.Preload(clause.Associations).First(&modifier, modifierID).Error; err != nil {
		shared.LogError("error getting modifier", LogDBRepository, "GetModifier", err, modifierID)
		return nil, err
	}

	return &modifier, nil
}

// ModifierFind method for find modifiers in database
func (r *DBRepository) ModifierFind(filters map[string]string) ([]Modifier, error) {
	var modifiers []Modifier
	if err := r.db.Preload(clause.Associations).Find(&modifiers, filters).Error; err != nil {
		shared.LogError("error getting modifiers", LogDBRepository, "FindModifier", err, filters)
		return nil, err
	}

	return modifiers, nil
}

// ModifierAddProduct method for add a product in modifier
func (r *DBRepository) ModifierAddProduct(product *Product, modifier *Modifier) (*Modifier, error) {
	if err := r.db.Model(&modifier).Association("Products").Append(product); err != nil {
		shared.LogError("error adding product to modifier", LogDBRepository, "AddProduct", err, modifier, product)
		return nil, err
	}

	return modifier, nil
}

// ModifierRemoveProduct method for remove a product in modifier
func (r *DBRepository) ModifierRemoveProduct(product *Product, modifier *Modifier) (*Modifier, error) {
	if err := r.db.Model(&modifier).Association("Products").Delete(product); err != nil {
		shared.LogError("error removing product to modifier", LogDBRepository, "RemoveProduct", err, modifier, product)
		return nil, err
	}

	return modifier, nil
}
