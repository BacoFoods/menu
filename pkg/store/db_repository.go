package store

import (
	"github.com/BacoFoods/menu/pkg/channel"
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/BacoFoods/menu/pkg/zones"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	if err := r.db.Preload(clause.Associations).Save(store).Error; err != nil {
		shared.LogError("error creating store", LogDBRepository, "Create", err, store)
		return nil, err
	}
	return store, nil
}

// Find method for find stores in database
func (r *DBRepository) Find(filters map[string]string) ([]Store, error) {
	var stores []Store
	if err := r.db.Preload(clause.Associations).Find(&stores, filters).Error; err != nil {
		shared.LogError("error getting stores", LogDBRepository, "Find", err, filters)
		return nil, err
	}
	return stores, nil
}

// Get method for get a store in database
func (r *DBRepository) Get(storeID string) (*Store, error) {
	var store Store
	if err := r.db.Preload(clause.Associations).First(&store, storeID).Error; err != nil {
		shared.LogError("error getting store", LogDBRepository, "Get", err, storeID)
		return nil, err
	}
	return &store, nil
}

// Update method for update a store in database
func (r *DBRepository) Update(store *Store) (*Store, error) {
	var storeDB Store
	if err := r.db.Preload(clause.Associations).First(&storeDB, store.ID).Error; err != nil {
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

// FindByIDs method for find stores in database by storeIDs
func (r *DBRepository) FindByIDs(storeIDs []string) ([]Store, error) {
	var stores []Store
	if err := r.db.Preload(clause.Associations).Find(&stores, storeIDs).Error; err != nil {
		shared.LogError("error getting stores", LogDBRepository, "FindByStores", err, storeIDs)
		return nil, err
	}
	return stores, nil
}

// AddChannel method for add a channel in store
func (r *DBRepository) AddChannel(storeID string, channel *channel.Channel) (*Store, error) {
	var store Store
	if err := r.db.Preload(clause.Associations).First(&store, storeID).Error; err != nil {
		shared.LogError("error getting store", LogDBRepository, "AddChannel", err, storeID)
		return nil, err
	}

	if err := r.db.Model(&store).Association("Channels").Append(channel); err != nil {
		shared.LogError("error adding channel to store", LogDBRepository, "AddChannel", err, storeID)
		return nil, err
	}

	return &store, nil
}

// FindZonesByStore method for get zones by store
func (r *DBRepository) FindZonesByStore(storeID string) ([]zones.Zone, error) {
	var zoneList []zones.Zone

	if err := r.db.Preload(clause.Associations).Where("store_id = ?", storeID).Find(&zoneList).Error; err != nil {
		shared.LogError("error getting zones", LogDBRepository, "GetZonesByStore", err, storeID)
		return nil, err
	}

	return zoneList, nil
}

// GetZoneByStore method for get zone by store
func (r *DBRepository) GetZoneByStore(storeID, zoneID string) (*zones.Zone, error) {
	var zone zones.Zone

	if err := r.db.Preload(clause.Associations).Where("store_id = ? AND id = ?", storeID, zoneID).First(&zone).Error; err != nil {
		shared.LogError("error getting zone", LogDBRepository, "GetZoneByStore", err, storeID, zoneID)
		return nil, err
	}

	return &zone, nil
}
