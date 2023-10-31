package assets

import "gorm.io/gorm"

type AssetRepository struct {
	DB *gorm.DB
}

// NewAssetRepository initializes a new asset repository
func NewAssetRepository(db *gorm.DB) *AssetRepository {
	return &AssetRepository{DB: db}
}

// FindByPlaca finds an asset by its Placa value
func (repo *AssetRepository) FindByPlaca(placa string) (*Asset, error) {
	var asset Asset
	err := repo.DB.Where("placa = ?", placa).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

// FindAll fetches all assets with pagination
func (repo *AssetRepository) FindAll(limit int, offset int) ([]Asset, error) {
	var assets []Asset
	err := repo.DB.Limit(limit).Offset(offset).Find(&assets).Error
	if err != nil {
		return nil, err
	}
	return assets, nil
}

// CreateMany creates multiple assets
func (repo *AssetRepository) CreateMany(assets []Asset) error {
	return repo.DB.Create(&assets).Error
}

// CreateMany creates multiple assets
func (repo *AssetRepository) Create(assets *Asset) error {
	return repo.DB.Create(assets).Error
}

// Update updates an asset based on its ID
func (repo *AssetRepository) Update(asset Asset) error {
	return repo.DB.Save(&asset).Error
}
