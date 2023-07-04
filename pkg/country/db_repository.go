package country

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const LogDBRepository string = "pkg/country/db_repository"

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db}
}

func (r *DBRepository) Create(country *Country) (*Country, error) {
	if err := r.db.Save(country).Error; err != nil {
		shared.LogError("error creating country", LogDBRepository, "Create", err, country)
		return nil, err
	}
	return country, nil
}

func (r *DBRepository) Find(query map[string]string) ([]Country, error) {
	var country []Country
	if err := r.db.Preload(clause.Associations).Find(&country, query).Error; err != nil {
		shared.LogError("error finding country", LogDBRepository, "Find", err)
		return nil, err
	}
	return country, nil
}

func (r *DBRepository) Get(countryID string) (*Country, error) {
	var country Country
	if err := r.db.Preload(clause.Associations).First(&country, countryID).Error; err != nil {
		shared.LogError("error getting country", LogDBRepository, "Find", err, countryID)
		return nil, err
	}
	return &country, nil
}

func (r *DBRepository) Update(country Country) (*Country, error) {
	var countryDB Country
	if err := r.db.First(&countryDB, country.ID).Error; err != nil {
		shared.LogError("error getting country", LogDBRepository, "Update", err, country.ID, country)
		return nil, err
	}

	if err := r.db.Model(&countryDB).Updates(country).Error; err != nil {
		shared.LogError("error updating country", LogDBRepository, "Update", err, country.ID, country, country)
		return nil, err
	}
	return &countryDB, nil
}

func (r *DBRepository) Delete(countryID string) (*Country, error) {
	var countryDB Country
	if err := r.db.First(&countryDB, countryID).Error; err != nil {
		shared.LogError("error getting country", LogDBRepository, "Delete", err, countryID)
		return nil, err
	}

	if err := r.db.Delete(&countryDB, countryID).Error; err != nil {
		shared.LogError("error deleting country", LogDBRepository, "Delete", err, countryDB)
		return nil, err
	}
	return &countryDB, nil
}
