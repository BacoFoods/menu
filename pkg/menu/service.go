package menu

import (
	overridersEntity "github.com/BacoFoods/menu/pkg/overriders"
	"github.com/BacoFoods/menu/pkg/product"
)

// Service is the interface that provides menu methods, used for dependency injection.
type Service interface {
	Find(map[string]string) ([]Menu, error)
	Get(string) (*Menu, error)
	Create(*Menu) (*Menu, error)
	Update(*Menu) (*Menu, error)
	Delete(string) (*Menu, error)

	GetByPlace(string, string, string) (*Menu, error)
}

// service is the default implementation of the Service interface for menu.
type service struct {
	repository Repository
	overriders overridersEntity.Repository
}

// NewService creates a new instance of the service for menu, using the provided repository implementation.
func NewService(repository Repository, overriders overridersEntity.Repository) service {
	return service{repository, overriders}
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

// GetByPlace returns a single menu object loading overriders by ID.
func (s service) GetByPlace(place, placeID, menuID string) (*Menu, error) {
	menuItems, err := s.repository.GetMenuItems(menuID)
	if err != nil {
		return nil, err
	}

	overriders, err := s.overriders.FindByPlace(place, placeID)
	if err != nil {
		return nil, err
	}

	var itemsByCategories map[uint][]product.Product
	for _, item := range menuItems {
		var prod product.Product
		for _, overrider := range overriders {
			if item.ID == *overrider.ProductID && IsAllowOverride(item, overrider) {
				prod = product.Product{
					ID:          item.ID,
					Name:        item.Name,
					Description: overrider.Description,
					Image:       overrider.Image,
					SKU:         item.SKU,
					Price:       overrider.Price,
					TaxID:       item.TaxID,
					DiscountID:  item.DiscountID,
					Unit:        item.Unit,
				}
			}
		}
		itemsByCategories[*item.CategoryID] = append(itemsByCategories[*item.CategoryID], prod)
	}

	menu, err := s.repository.Get(menuID)
	if err != nil {
		return nil, err
	}

	for _, category := range menu.Categories {
		category.Products = itemsByCategories[category.ID]
	}

	return menu, nil
}
