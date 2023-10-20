package product

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/availability"
	channels "github.com/BacoFoods/menu/pkg/channel"
	"github.com/BacoFoods/menu/pkg/shared"
)

type Service interface {
	Find(map[string]string) ([]Product, error)
	Get(productID string) (*Product, error)
	Create(*Product) (*Product, error)
	Update(*Product) (*Product, error)
	Delete(string) (*Product, error)
	AddModifier(productID, modifierID string) (*Product, error)
	RemoveModifier(productID, modifierID string) (*Product, error)
	GetOverriders(productID, field string) ([]OverriderDTO, error)
	UpdateAllOverriders(productID, field string, value any) error
	GetCategory(productID string) ([]CategoryDTO, error)

	ModifierFind(map[string]string) ([]Modifier, error)
	ModifierCreate(*Modifier) (*Modifier, error)
	ModifierAddProduct(productID, modifierID string) (*Modifier, error)
	ModifierRemoveProduct(productID, modifierID string) (*Modifier, error)
	ModifierUpdate(*Modifier) (*Modifier, error)

	OverriderFind(map[string]string) ([]Overrider, error)
	OverriderGet(string) (*Overrider, error)
	OverriderCreate(*Overrider) (*Overrider, error)
	OverriderUpdate(*Overrider) (*Overrider, error)
	OverriderDelete(string) (*Overrider, error)
}

type service struct {
	repository Repository
	channel    channels.Repository
}

func NewService(repository Repository, channel channels.Repository) service {
	return service{repository, channel}
}

// Product

func (s service) Find(filter map[string]string) ([]Product, error) {
	return s.repository.Find(filter)
}

func (s service) Get(productID string) (*Product, error) {
	return s.repository.Get(productID)
}

func (s service) Create(product *Product) (*Product, error) {
	product, err := s.repository.Create(product)
	if err != nil {
		shared.LogError("error creating product", LogDomain, "Create", err)
		return nil, err
	}

	chanelList, err := s.channel.Find(map[string]string{
		"brand_id": fmt.Sprintf("%d", *product.BrandID),
	})
	if err != nil {
		shared.LogError("error finding channels", LogDomain, "Create", err)
		return nil, err
	}

	defaultOverriders := make([]Overrider, 0)
	for _, c := range chanelList {
		productID := product.ID
		channelID := c.ID
		defaultOverriders = append(defaultOverriders, Overrider{
			ProductID:   &productID,
			Place:       string(availability.PlaceChannel),
			PlaceID:     &channelID,
			Name:        product.Name,
			Description: product.Description,
			Image:       product.Image,
			Price:       product.Price,
			Enable:      false,
		})
	}

	if err := s.repository.OverriderCreateAll(defaultOverriders); err != nil {
		shared.LogError("error creating default overriders", LogDomain, "Create", err)
		return nil, err
	}

	return product, nil
}

func (s service) Update(product *Product) (*Product, error) {
	return s.repository.Update(product)
}

func (s service) Delete(productID string) (*Product, error) {
	return s.repository.Delete(productID)
}

func (s service) AddModifier(productID, modifierID string) (*Product, error) {
	product, err := s.repository.Get(productID)
	if err != nil {
		return nil, err
	}

	modifier, err := s.repository.ModifierGet(modifierID)
	if err != nil {
		return nil, err
	}

	return s.repository.AddModifier(product, modifier)
}

func (s service) RemoveModifier(productID, modifierID string) (*Product, error) {
	product, err := s.repository.Get(productID)
	if err != nil {
		return nil, err
	}

	modifier, err := s.repository.ModifierGet(modifierID)
	if err != nil {
		return nil, err
	}

	return s.repository.RemoveModifier(product, modifier)
}

func (s service) GetOverriders(productID, field string) ([]OverriderDTO, error) {
	return s.repository.GetOverriders(productID, field)
}

func (s service) GetCategory(productID string) ([]CategoryDTO, error) {
	return s.repository.GetCategory(productID)
}

// Modifier

func (s service) ModifierFind(filter map[string]string) ([]Modifier, error) {
	return s.repository.ModifierFind(filter)
}

func (s service) ModifierCreate(modifier *Modifier) (*Modifier, error) {
	return s.repository.ModifierCreate(modifier)
}

func (s service) ModifierAddProduct(productID, modifierID string) (*Modifier, error) {
	product, err := s.repository.Get(productID)
	if err != nil {
		return nil, err
	}

	modifier, err := s.repository.ModifierGet(modifierID)
	if err != nil {
		return nil, err
	}

	return s.repository.ModifierAddProduct(product, modifier)
}

func (s service) ModifierRemoveProduct(productID, modifierID string) (*Modifier, error) {
	product, err := s.repository.Get(productID)
	if err != nil {
		return nil, err
	}

	modifier, err := s.repository.ModifierGet(modifierID)
	if err != nil {
		return nil, err
	}

	return s.repository.ModifierRemoveProduct(product, modifier)
}

func (s service) UpdateAllOverriders(productID, field string, value any) error {
	overridersIDs, err := s.repository.GetOverriderIDs(productID)
	if err != nil {
		return err
	}

	return s.repository.UpdateOverriders(overridersIDs, field, value)
}

func (s service) ModifierUpdate(modifier *Modifier) (*Modifier, error) {
	return s.repository.ModifierUpdate(modifier)
}

// Overrider

func (s service) OverriderFind(filter map[string]string) ([]Overrider, error) {
	return s.repository.OverriderFind(filter)
}

func (s service) OverriderGet(overriderID string) (*Overrider, error) {
	return s.repository.OverriderGet(overriderID)
}

func (s service) OverriderCreate(overrider *Overrider) (*Overrider, error) {
	return s.repository.OverriderCreate(overrider)
}

func (s service) OverriderUpdate(overrider *Overrider) (*Overrider, error) {
	return s.repository.OverriderUpdate(overrider)
}

func (s service) OverriderDelete(overriderID string) (*Overrider, error) {
	return s.repository.OverriderDelete(overriderID)
}
