package assets_test

import (
	"errors"

	"github.com/BacoFoods/menu/pkg/assets"
)

type MockAssetRepository struct {
	Assets []assets.Asset
}

func (m *MockAssetRepository) Create(asset *assets.Asset) error {
	m.Assets = append(m.Assets, *asset)
	asset.ID = uint(len(m.Assets))
	return nil
}

func (m *MockAssetRepository) FindByPlaca(placa string) (*assets.Asset, error) {
	for _, a := range m.Assets {
		if a.Placa == placa {
			return &a, nil
		}
	}

	return nil, nil
}

func (m *MockAssetRepository) FindAll(limit int, offset int) ([]assets.Asset, error) {
	return m.Assets, nil
}

func (m *MockAssetRepository) Update(asset assets.Asset) error {
	for i, a := range m.Assets {
		if a.ID == asset.ID {
			m.Assets[i] = asset
			return nil
		}
	}
	return errors.New(assets.ErrorAssetNotFound)
}
