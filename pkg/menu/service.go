package menu

import (
	availabilityPkg "github.com/BacoFoods/menu/pkg/availability"
	overridersPkg "github.com/BacoFoods/menu/pkg/overriders"
	"github.com/BacoFoods/menu/pkg/shared"
	storePkg "github.com/BacoFoods/menu/pkg/store"
	"strconv"
)

const (
	LogService string = "pkg/menu/service"
)

// Service is the interface that provides menu methods, used for dependency injection.
type Service interface {
	Find(map[string]string) ([]Menu, error)
	Get(string) (*Menu, error)
	Create(name string, brandID uint, place string, stores []uint) (*Menu, error)
	Update(*Menu) (*Menu, error)
	Delete(string) (*Menu, error)

	FindByPlace(string, string) ([]Menu, error)
	GetByPlace(string, string, string) (*Menu, error)

	UpdateAvailability(menuID, place string, placeIDs map[uint]bool) (*Menu, error)

	FindChannels(menuID, storeID string) ([]any, error)
}

// service is the default implementation of the Service interface for menu.
type service struct {
	repository   Repository
	overriders   overridersPkg.Repository
	availability availabilityPkg.Repository
	store        storePkg.Repository
}

// NewService creates a new instance of the service for menu, using the provided repository implementation.
func NewService(repository Repository,
	overriders overridersPkg.Repository,
	availability availabilityPkg.Repository,
	store storePkg.Repository) service {
	return service{repository, overriders, availability, store}
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
func (s service) Create(name string, brandID uint, placeName string, placeIds []uint) (*Menu, error) {
	menu := &Menu{
		Name:    name,
		BrandID: &brandID,
		Enable:  true,
	}

	menu, err := s.repository.Create(menu)
	if err != nil {
		return nil, err
	}

	place, err := availabilityPkg.GetPlace(placeName)
	if err != nil {
		return nil, err
	}

	for _, placeId := range placeIds {
		if err := s.availability.EnableEntity(
			availabilityPkg.EntityMenu,
			place,
			menu.ID,
			placeId,
			true,
		); err != nil {
			shared.LogError("error enabling menu", LogService, "Create", err, name, place, brandID, placeId)
			return nil, err
		}
	}

	return menu, nil
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

	placeName, err := availabilityPkg.GetPlace(place)
	if err != nil {
		return []Menu{}, err
	}

	availabilities, err := s.availability.FindEntityByPlace(availabilityPkg.EntityMenu, placeName, placeID)
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

// UpdateAvailability updates the availability of a menu.
func (s service) UpdateAvailability(menuID string, placeName string, placeIDs map[uint]bool) (*Menu, error) {
	place, err := availabilityPkg.GetPlace(placeName)
	if err != nil {
		return nil, err
	}

	menu, err := s.repository.Get(menuID)
	if err != nil {
		return nil, err
	}

	availabilities, err := s.availability.FindPlacesByEntity(availabilityPkg.EntityMenu, menu.ID, placeName)
	if err != nil {
		return nil, err
	}

	for _, availability := range availabilities {
		enable, ok := placeIDs[*availability.PlaceID]
		if ok {
			availability.Enable = enable
			if err := s.availability.EnableEntity(availabilityPkg.EntityMenu, place, menu.ID, *availability.PlaceID, enable); err != nil {
				return nil, err
			}
		}
	}

	return menu, nil
}

// FindChannels returns a list of channels available for a menu.
func (s service) FindChannels(menuID, storeID string) ([]any, error) {
	store, err := s.store.Get(storeID)
	if err != nil {
		return nil, err
	}

	menuId, err := strconv.ParseUint(menuID, 10, 64)
	if err != nil {
		return nil, err
	}

	availabilities, err := s.availability.FindPlacesByEntity(availabilityPkg.EntityMenu, uint(menuId), string(availabilityPkg.PlaceChannel))
	if err != nil {
		return nil, err
	}

	availabilityByPlace := availabilityPkg.MapAvailabilityByPlace(availabilities)

	menuChannels := make([]any, 0)
	for _, channel := range store.Channels {
		availability, ok := availabilityByPlace[channel.ID]
		if ok {
			menuChannels = append(menuChannels, struct {
				ID     uint   `json:"id"`
				Enable bool   `json:"enable"`
				Name   string `json:"name"`
			}{
				ID:     channel.ID,
				Enable: availability.Enable,
				Name:   channel.Name,
			})
		} else {
			menuChannels = append(menuChannels, struct {
				ID     uint   `json:"id"`
				Enable bool   `json:"enable"`
				Name   string `json:"name"`
			}{
				ID:     channel.ID,
				Enable: false,
				Name:   channel.Name,
			})
		}
	}

	return menuChannels, nil
}
