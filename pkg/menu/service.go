package menu

import (
	availabilityPkg "github.com/BacoFoods/menu/pkg/availability"
	overridersPkg "github.com/BacoFoods/menu/pkg/overriders"
)

// Service is the interface that provides menu methods, used for dependency injection.
type Service interface {
	Find(map[string]string) ([]Menu, error)
	Get(string) (*Menu, error)
	Create(*Menu) (*Menu, error)
	Update(*Menu) (*Menu, error)
	Delete(string) (*Menu, error)

	FindByPlace(string, string) ([]Menu, error)
	GetByPlace(string, string, string) (*Menu, error)
}

// service is the default implementation of the Service interface for menu.
type service struct {
	repository   Repository
	overriders   overridersPkg.Repository
	availability availabilityPkg.Repository
}

// NewService creates a new instance of the service for menu, using the provided repository implementation.
func NewService(repository Repository,
	overriders overridersPkg.Repository,
	availability availabilityPkg.Repository) service {
	return service{repository, overriders, availability}
}

// Find returns a list of menu objects filtering by query map.
func (s service) Find(filter map[string]string) ([]Menu, error) {
	return s.repository.Find(filter)
}

// Get returns a single menu object by ID.
func (s service) Get(menuID string) (*Menu, error) {
	return s.repository.Get(menuID)
}

// Create creates a new menu object.
func (s service) Create(menu *Menu) (*Menu, error) {
	return s.repository.Create(menu)
}

// Update updates an existing menu object only the fields that are present in the provided object.
// this method doesn't create new register if the provided id doesn't exist.
func (s service) Update(menu *Menu) (*Menu, error) {
	return s.repository.Update(menu)
}

// Delete deletes an existing menu object.
func (s service) Delete(menuID string) (*Menu, error) {
	return s.repository.Delete(menuID)
}

// FindByPlace returns a list of menu objects filtering by place and placeID.
func (s service) FindByPlace(place, placeID string) ([]Menu, error) {
	menus, err := s.repository.FindByPlace(place, placeID)
	if err != nil {
		return []Menu{}, err
	}

	if len(menus) == 0 {
		return []Menu{}, nil
	}

	availabilities, err := s.availability.FindEntityByPlace(availabilityPkg.EntityMenu, place, placeID)
	if err != nil {
		return []Menu{}, err
	}

	if len(availabilities) == 0 {
		return menus, nil
	}

	var menuList []Menu
	for _, menu := range menus {
		for _, availability := range availabilities {
			if menu.ID == *availability.EntityID {
				menu.Enable = availability.Enable
			}
			menuList = append(menuList, menu)
		}
	}

	return menuList, nil
}

// GetByPlace returns a single menu object loading overriders by ID.
func (s service) GetByPlace(place, placeID, menuID string) (*Menu, error) {
	menuItems, err := s.repository.GetMenuItems(menuID)
	if err != nil {
		return nil, err
	}

	if len(menuItems) == 0 {
		return nil, nil
	}

	overriders, err := s.overriders.FindByPlace(place, placeID)
	if err != nil {
		return nil, err
	}

	productsByCategory := OverrideProducts(menuItems, overriders)

	menu, err := s.repository.Get(menuID)
	if err != nil {
		return nil, err
	}

	for _, category := range menu.Categories {
		category.Products = productsByCategory[category.ID]
	}

	return menu, nil
}
