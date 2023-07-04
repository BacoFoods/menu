package menu

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
)

const LogDBRepository string = "pkg/menu/db_repository"

type DBRepository struct {
	db *gorm.DB
}

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
	if err := r.db.Find(&menus, filters).Error; err != nil {
		shared.LogError("error getting menus", LogDBRepository, "Find", err, filters)
		return nil, err
	}
	return menus, nil
}

// Get method for get a menu in database
func (r *DBRepository) Get(menuID string) (*Menu, error) {
	var menu Menu
	if err := r.db.First(&menu, menuID).Error; err != nil {
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
