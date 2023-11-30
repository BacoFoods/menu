package product

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

const LogDBRepository string = "pkg/product/db_repository"

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db: db}
}

// Product

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
	if strings.TrimSpace(productID) == "" {
		err := fmt.Errorf(ErrorProductIDEmpty)
		shared.LogWarn("error getting product", LogDBRepository, "Get", err)
		return nil, err
	}

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
	if err := r.db.Where("id in ?", productIDs).Preload(clause.Associations).Find(&products).Error; err != nil {
		shared.LogError("error getting products", LogDBRepository, "GetByIDs", err, productIDs)
		return nil, err
	}
	return products, nil
}

// GetAsMapByIDs method for get products as map by ids in database
func (r *DBRepository) GetAsMapByIDs(productIDs []string) (map[string]Product, error) {
	var products []Product
	if err := r.db.Preload(clause.Associations).
		Preload("Modifiers.Products").
		Find(&products, productIDs).Error; err != nil {
		shared.LogError("error getting products", LogDBRepository, "GetAsMapByIDs", err, productIDs)
		return nil, err
	}

	productsMap := make(map[string]Product)
	for _, product := range products {
		productsMap[fmt.Sprintf("%d", product.ID)] = product
	}

	return productsMap, nil
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

// GetOverriders method for get overriders in product
func (r *DBRepository) GetOverriders(productID, field string) ([]OverriderDTO, error) {
	var overriders []OverriderDTO

	if err := r.db.Table("overriders as o").
		Select(fmt.Sprintf("o.id as id, o.product_id as product_id, c.name as place_name, c.id as place_id, o.%s as field_value", field)).
		Joins("left join channels c on o.place_id = c.id").
		Where("o.product_id = ?", productID).
		Scan(&overriders).Error; err != nil {
		shared.LogError("error getting overriders", LogDBRepository, "GetOverriders", err)
		return nil, err
	}

	return overriders, nil
}

// GetOverriderIDs method for get overrider ids in product
func (r *DBRepository) GetOverriderIDs(productID string) ([]uint, error) {
	var overriderIDs []uint

	if err := r.db.Table("overriders").
		Distinct().
		Pluck("id", &overriderIDs).
		Where("product_id = ?", productID).Error; err != nil {
		shared.LogError("error getting overriders ids", LogDBRepository, "GetOverriderIDs", err)
		return nil, err
	}

	return overriderIDs, nil
}

// UpdateOverriders method for update overriders in product
func (r *DBRepository) UpdateOverriders(ids []uint, field string, value any) error {
	if err := r.db.Table("overriders").
		Where("id in (?)", ids).
		Update(field, value).Error; err != nil {
		shared.LogError("error updating overriders", LogDBRepository, "UpdateOverriders", err)
		return err
	}
	return nil
}

// GetCategory method for get categories by product id
func (r *DBRepository) GetCategory(productID string) ([]CategoryDTO, error) {
	var categories []CategoryDTO
	if err := r.db.Table("categories as c").
		Select("c.id as id, c.name as name").
		Joins("left join categories_products pc on c.id = pc.category_id").
		Where("pc.product_id = ?", productID).
		Scan(&categories).Error; err != nil {
		shared.LogError("error getting categories", LogDBRepository, "GetCategory", err)
		return nil, err
	}

	return categories, nil
}

// Modifier

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

// ModifierUpdate method for update a modifier in database
func (r *DBRepository) ModifierUpdate(modifier *Modifier) (*Modifier, error) {
	var modifierDB Modifier
	if err := r.db.First(&modifierDB, modifier.ID).Error; err != nil {
		shared.LogError("error getting modifier", LogDBRepository, "UpdateModifier", err, *modifier)
		return nil, err
	}

	// To avoid zero values update trouble
	modifierMap := make(map[string]any)
	modifierMap["id"] = modifier.ID
	modifierMap["name"] = modifier.Name
	modifierMap["description"] = modifier.Description
	modifierMap["image"] = modifier.Image
	modifierMap["apply_price"] = modifier.ApplyPrice
	modifierMap["category"] = modifier.Category
	modifierMap["brand_id"] = modifier.BrandID

	if err := r.db.Model(&modifierDB).Updates(modifierMap).Error; err != nil {
		shared.LogError("error updating modifier", LogDBRepository, "UpdateModifier", err, *modifier)
		return nil, err
	}

	return &modifierDB, nil
}

// Overrider

// OverriderCreate method for create a new overrider in database
func (r *DBRepository) OverriderCreate(overrider *Overrider) (*Overrider, error) {
	if err := r.db.Save(overrider).Error; err != nil {
		shared.LogError("error creating overrider", LogDBRepository, "Create", err, overrider)
		return nil, err
	}

	return overrider, nil
}

// OverriderCreateAll method for create a new overrider in database
func (r *DBRepository) OverriderCreateAll(overrider []Overrider) error {
	if err := r.db.Create(&overrider).Error; err != nil {
		shared.LogError("error creating overriders in batch", LogDBRepository, "CreateAll", err, overrider)
		return err
	}

	return nil
}

// OverriderFind method for find overriders in database
func (r *DBRepository) OverriderFind(filters map[string]string) ([]Overrider, error) {
	var overrider []Overrider
	if err := r.db.Find(&overrider, filters).Error; err != nil {
		shared.LogError("error getting overriders", LogDBRepository, "Find", err, filters)
		return nil, err
	}

	return overrider, nil
}

// OverriderGet method for get an overrider in database
func (r *DBRepository) OverriderGet(overriderID string) (*Overrider, error) {
	var overrider Overrider
	if err := r.db.First(&overrider, overriderID).Error; err != nil {
		shared.LogError("error getting overrider", LogDBRepository, "Get", err, overriderID)
		return nil, err
	}

	return &overrider, nil
}

// OverriderUpdate method for update an overrider in database
func (r *DBRepository) OverriderUpdate(overrider *Overrider) (*Overrider, error) {
	var overriderDB Overrider
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

// OverriderDelete method for delete an overrider in database
func (r *DBRepository) OverriderDelete(overriderID string) (*Overrider, error) {
	var overrider Overrider

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

// OverriderFindByPlace method for find overriders in database
func (r *DBRepository) OverriderFindByPlace(place, placeID string) ([]Overrider, error) {
	var overrider []Overrider

	if err := r.db.Where("place = ? AND place_id = ?", place, placeID).Find(&overrider).Error; err != nil {
		shared.LogError("error getting overrider", LogDBRepository, "FindByPlace", err, place, placeID)
		return nil, err
	}

	return overrider, nil
}
