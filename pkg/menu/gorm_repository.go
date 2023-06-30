package menu

import "gorm.io/gorm"

func (b Menu) JoinTable(db gorm.DB) error {
	err := db.SetupJoinTable(&Menu{}, "Categories", &MenusCategories{})
	if err != nil {
		return err
	}
	return nil
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db}
}

func (r *GormRepository) Create(menu *Menu) (*Menu, error) {
	if err := r.db.Create(menu).Error; err != nil {
		return nil, err
	}
	return menu, nil
}

func (r *GormRepository) Find() ([]Menu, error) {
	var menus []Menu
	if err := r.db.Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *GormRepository) Get(id string) (*Menu, error) {
	var menu Menu
	if err := r.db.First(&menu, id).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *GormRepository) Update(menu *Menu) (*Menu, error) {
	if err := r.db.Save(menu).Error; err != nil {
		return nil, err
	}
	return menu, nil
}

func (r *GormRepository) Delete(id string) (*Menu, error) {
	var menu Menu
	if err := r.db.Delete(&menu, id).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}
