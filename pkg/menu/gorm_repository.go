package menu

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
)

const LogGormRepository string = "pkg/menu/repository"

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

// Create method for create a new menu in database
func (r *GormRepository) Create(menu *Menu) (*Menu, error) {
	if err := r.db.Save(menu).Error; err != nil {
		shared.LogError("error creating menu", LogGormRepository, "Create", err, menu)
		return nil, err
	}
	return menu, nil
}

// Find method for find menus in database
func (r *GormRepository) Find(filters map[string]string) ([]Menu, error) {
	var menus []Menu
	if err := r.db.Find(&menus, filters).Error; err != nil {
		shared.LogError("error getting menus", LogGormRepository, "Find", err, filters)
		return nil, err
	}
	return menus, nil
}

// Get method for get a menu in database
func (r *GormRepository) Get(menuID string) (*Menu, error) {
	var menu Menu
	if err := r.db.First(&menu, menuID).Error; err != nil {
		shared.LogError("error getting menu", LogGormRepository, "Get", err, menuID)
		return nil, err
	}
	return &menu, nil
}

// Update method for update a menu in database
func (r *GormRepository) Update(menu *Menu) (*Menu, error) {
	var menuDB Menu
	if err := r.db.First(&menuDB, menu.ID).Error; err != nil {
		shared.LogError("error getting menu", LogGormRepository, "Update", err, menu)
		return nil, err
	}
	if err := r.db.Model(&menuDB).Updates(menu).Error; err != nil {
		shared.LogError("error updating menu", LogGormRepository, "Update", err, menu)
		return nil, err
	}
	return &menuDB, nil
}

// Delete method for delete a menu in database
func (r *GormRepository) Delete(menuID string) (*Menu, error) {
	var menu Menu
	if err := r.db.First(&menu, menuID).Error; err != nil {
		shared.LogError("error getting menu", LogGormRepository, "Delete", err, menuID)
		return nil, err
	}
	if err := r.db.Delete(&menu).Error; err != nil {
		shared.LogError("error deleting menu", LogGormRepository, "Delete", err, menu)
		return nil, err
	}
	return &menu, nil
}
