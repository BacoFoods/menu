package country

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const LogGormRepository string = "pkg/country/repository"

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db}
}

func (r *GormRepository) Create(country *Country) (*Country, error) {
	if err := r.db.Save(country).Error; err != nil {
		shared.LogError("error creating country", LogGormRepository, "Create", err, country)
		return nil, err
	}
	return country, nil
}

func (r *GormRepository) Find(query map[string]string) ([]Country, error) {
	var country []Country
	if err := r.db.Preload(clause.Associations).Find(&country, query).Error; err != nil {
		shared.LogError("error finding country", LogGormRepository, "Find", err)
		return nil, err
	}
	return country, nil
}

func (r *GormRepository) Get(countryID string) (*Country, error) {
	var country Country
	if err := r.db.Preload(clause.Associations).First(&country, countryID).Error; err != nil {
		shared.LogError("error getting country", LogGormRepository, "Find", err, countryID)
		return nil, err
	}
	return &country, nil
}

func (r *GormRepository) Update(country Country) (*Country, error) {
	var countryDB Country
	if err := r.db.First(&countryDB, country.ID).Error; err != nil {
		shared.LogError("error getting country", LogGormRepository, "Update", err, country.ID, country)
		return nil, err
	}

	if err := r.db.Model(&countryDB).Updates(country).Error; err != nil {
		shared.LogError("error updating country", LogGormRepository, "Update", err, country.ID, country, country)
		return nil, err
	}
	return &countryDB, nil
}

func (r *GormRepository) Delete(countryID string) (*Country, error) {
	var countryDB Country
	if err := r.db.First(&countryDB, countryID).Error; err != nil {
		shared.LogError("error getting country", LogGormRepository, "Delete", err, countryID)
		return nil, err
	}

	if err := r.db.Delete(&countryDB, countryID).Error; err != nil {
		shared.LogError("error deleting country", LogGormRepository, "Delete", err, countryDB)
		return nil, err
	}
	return &countryDB, nil
}
