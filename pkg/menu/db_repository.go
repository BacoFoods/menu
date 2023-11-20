package menu

import (
	"fmt"
	"strings"

	"github.com/BacoFoods/menu/pkg/category"
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const LogDBRepository string = "pkg/menu/db_repository"

// DBRepository struct for database repository
type DBRepository struct {
	db *gorm.DB
}

// NewDBRepository method for create a new instance of DBRepository
func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db: db}
}

// Create method for create a new menu in database
func (r *DBRepository) Create(menu *Menu) (*Menu, error) {
	if err := r.db.Save(menu).Error; err != nil {
		shared.LogError("error creating menu", LogDBRepository, "Create", err, menu)
		return nil, err
	}
	return menu, nil
}

// Find method for find menus in database
func (r *DBRepository) Find(filters map[string]string) ([]Menu, error) {
	var menus []Menu
	if err := r.db.Preload(clause.Associations).
		Find(&menus, filters).Error; err != nil {
		shared.LogError("error getting menus", LogDBRepository, "Find", err, filters)
		return nil, err
	}
	return menus, nil
}

// Get method for get a menu in database
func (r *DBRepository) Get(menuID string) (*Menu, error) {
	if strings.TrimSpace(menuID) == "" {
		err := fmt.Errorf(ErrorMenuIDEmpty)
		shared.LogWarn("error getting menu", LogDBRepository, "Get", err)
		return nil, err
	}

	var menu Menu
	if err := r.db.Preload(clause.Associations).
		First(&menu, menuID).Error; err != nil {
		shared.LogError("error getting menu", LogDBRepository, "Get", err, menuID)
		return nil, err
	}
	return &menu, nil
}

// Update method for update a menu in database
func (r *DBRepository) Update(menu *Menu) (*Menu, error) {
	var menuDB Menu
	if err := r.db.First(&menuDB, menu.ID).Error; err != nil {
		shared.LogError("error getting menu", LogDBRepository, "Update", err, menu)
		return nil, err
	}
	if err := r.db.Model(&menuDB).Updates(menu).Error; err != nil {
		shared.LogError("error updating menu", LogDBRepository, "Update", err, menu)
		return nil, err
	}
	return &menuDB, nil
}

// Delete method for delete a menu in database
func (r *DBRepository) Delete(menuID string) (*Menu, error) {
	var menu Menu
	if err := r.db.First(&menu, menuID).Error; err != nil {
		shared.LogError("error getting menu", LogDBRepository, "Delete", err, menuID)
		return nil, err
	}
	if err := r.db.Delete(&menu).Error; err != nil {
		shared.LogError("error deleting menu", LogDBRepository, "Delete", err, menu)
		return nil, err
	}
	return &menu, nil
}

// FindByPlace method for find menus by place
func (r *DBRepository) FindByPlace(place, placeID string) ([]Menu, error) {
	var menus []Menu

	// Getting brandID from placeID
	brandID := "0"
	switch place {
	case "brand":
		brandID = placeID

	case "store":
		if err := r.db.Select("brand_id").Table("stores").Where("id = ?", placeID).Scan(&brandID).Error; err != nil {
			shared.LogError("error getting brand_id", LogDBRepository, "FindByPlace", err, place, placeID)
			return nil, err
		}

	case "channel":
		if err := r.db.Select("brand_id").Table("channels").Where("id = ?", placeID).Scan(&brandID).Error; err != nil {
			shared.LogError("error getting brand_id", LogDBRepository, "FindByPlace", err, place, placeID)
			return nil, err
		}

	default:
		return nil, fmt.Errorf(ErrorMenuFindingByPlace)
	}

	// Getting Menu by brandID
	if err := r.db.Preload(clause.Associations).
		Preload("Categories.Products.Modifiers.Products").
		Find(&menus, "brand_id = ?", brandID).Error; err != nil {
		shared.LogError("error getting menus", LogDBRepository, "FindByPlace", err, brandID, place, placeID)
		return nil, err
	}

	return menus, nil
}

// GetMenuItems method for get menu items in database
func (r *DBRepository) GetMenuItems(menuID string) ([]Item, error) {
	var products []Item
	if err := r.db.Table("products").
		Select("mc.menu_id menu_id, cp.category_id category_id, products.* ").
		Joins("left join categories_products cp on products.id = cp.product_id").
		Joins("left join menus_categories mc on cp.category_id = mc.category_id").
		Where("mc.menu_id = ?", menuID).
		Find(&products).Error; err != nil {
		shared.LogError("error getting menu items", LogDBRepository, "GetMenuItems", err, menuID)
		return nil, err
	}
	return products, nil
}

// AddCategory method for add a category to menu
func (r *DBRepository) AddCategory(menuID string, category *category.Category) (*Menu, error) {
	var menu Menu
	if err := r.db.Preload(clause.Associations).First(&menu, menuID).Error; err != nil {
		shared.LogError("error getting menu", LogDBRepository, "AddCategory", err, menuID, *category)
		return nil, err
	}

	if *menu.BrandID != *category.BrandID {
		err := fmt.Errorf(ErrorMenuWrongBrand)
		shared.LogError("error adding category to menu", LogDBRepository, "AddCategory", err, menuID, *category)
		return nil, err
	}

	if err := r.db.Model(&menu).Association("Categories").Append(category); err != nil {
		shared.LogError("error adding category to menu", LogDBRepository, "AddCategory", err, menuID, *category)
		return nil, err
	}

	return &menu, nil
}

// RemoveCategory method for remove a category from menu
func (r *DBRepository) RemoveCategory(menuID string, category *category.Category) (*Menu, error) {
	var menu Menu
	if err := r.db.Preload(clause.Associations).First(&menu, menuID).Error; err != nil {
		shared.LogError("error getting menu", LogDBRepository, "RemoveCategory", err, menuID, *category)
		return nil, err
	}

	if *menu.BrandID != *category.BrandID {
		err := fmt.Errorf(ErrorMenuWrongBrand)
		shared.LogError("error removing category from menu", LogDBRepository, "RemoveCategory", err, menuID, *category)
		return nil, err
	}

	if err := r.db.Model(&menu).Association("Categories").Delete(category); err != nil {
		shared.LogError("error removing category from menu", LogDBRepository, "RemoveCategory", err, menuID, *category)
		return nil, err
	}

	return &menu, nil
}
