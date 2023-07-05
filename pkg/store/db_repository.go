package store

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
)

const LogDBRepository string = "pkg/store/db_repository"

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db: db}
}

// Create method for create a new store in database
func (r *DBRepository) Create(store *Store) (*Store, error) {
	if err := r.db.Save(store).Error; err != nil {
		shared.LogError("error creating store", LogDBRepository, "Create", err, store)
		return nil, err
	}
	return store, nil
}

// Find method for find stores in database
func (r *DBRepository) Find(filters map[string]string) ([]Store, error) {
	var stores []Store
	if err := r.db.Find(&stores, filters).Error; err != nil {
		shared.LogError("error getting stores", LogDBRepository, "Find", err, filters)
		return nil, err
	}
	return stores, nil
}

// Get method for get a store in database
func (r *DBRepository) Get(storeID string) (*Store, error) {
	var store Store
	if err := r.db.First(&store, storeID).Error; err != nil {
		shared.LogError("error getting store", LogDBRepository, "Get", err, storeID)
		return nil, err
	}
	return &store, nil
}

// Update method for update a store in database
func (r *DBRepository) Update(store *Store) (*Store, error) {
	var storeDB Store
	if err := r.db.First(&storeDB, store.ID).Error; err != nil {
		shared.LogError("error getting store", LogDBRepository, "Update", err, store)
		return nil, err
	}
	if err := r.db.Model(&storeDB).Updates(store).Error; err != nil {
		shared.LogError("error updating store", LogDBRepository, "Update", err, store)
		return nil, err
	}
	return &storeDB, nil
}

// Delete method for delete a store in database
func (r *DBRepository) Delete(storeID string) (*Store, error) {
	var store Store
	if err := r.db.First(&store, storeID).Error; err != nil {
		shared.LogError("error getting store", LogDBRepository, "Delete", err, storeID)
		return nil, err
	}

	if err := r.db.Delete(&store).Error; err != nil {
		shared.LogError("error deleting store", LogDBRepository, "Delete", err, store)
		return nil, err
	}

	return &store, nil
}
