package zones

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/shared"
	tablesPKG "github.com/BacoFoods/menu/pkg/tables"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	LogDBRepository = "pkg/zones/db_repository.go"
)

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db: db}
}

// Find method for find zones in database
func (r *DBRepository) Find(filters map[string]any) ([]Zone, error) {
	var zones []Zone
	if err := r.db.Preload(clause.Associations).Find(&zones, filters).Error; err != nil {
		shared.LogError("Error finding zones", LogDBRepository, "Find", err, filters)
		return nil, err
	}
	return zones, nil
}

// GetZone method for get a zone in database
func (r *DBRepository) GetZone(zoneID string) (*Zone, error) {
	var zone Zone
	if err := r.db.Preload(clause.Associations).First(&zone, zoneID).Error; err != nil {
		shared.LogError("Error finding zone", LogDBRepository, "GetZone", err, zoneID)
		return nil, err
	}
	return &zone, nil
}

// Create method for create a zone in database
func (r *DBRepository) Create(zone *Zone) (*Zone, error) {
	if err := r.db.Save(zone).Error; err != nil {
		shared.LogError("Error creating zone", LogDBRepository, "Create", nil, *zone)
		return nil, err
	}
	return zone, nil
}

// Update method for update a zone in database
func (r *DBRepository) Update(zoneID string, zone *Zone) (*Zone, error) {
	var zoneDB Zone
	if err := r.db.First(&zoneDB, zoneID).Error; err != nil {
		shared.LogError("Error finding zone", LogDBRepository, "Update", err, zoneID)
		return nil, err
	}

	if err := r.db.Model(&zoneDB).Updates(zone).Error; err != nil {
		shared.LogError("Error updating zone", LogDBRepository, "Update", err, *zone)
		return nil, err
	}

	return zone, nil
}

// Delete method for delete a zone in database
func (r *DBRepository) Delete(zoneID string) error {
	var zone Zone
	if err := r.db.First(&zone, zoneID).Error; err != nil {
		shared.LogError("Error finding zone", LogDBRepository, "Delete", err, zoneID)
		return err
	}

	if err := r.db.Delete(&zone).Error; err != nil {
		shared.LogError("Error deleting zone", LogDBRepository, "Delete", err, zone)
		return err
	}

	return nil
}

// AddTables method for add tables to zone in database
func (r *DBRepository) AddTables(zone *Zone, tables []uint) error {
	var tablesDB []tablesPKG.Table
	if err := r.db.Find(&tablesDB, tables).Error; err != nil {
		shared.LogError("Error finding tables", LogDBRepository, "AddTables", err, tables)
		return err
	}

	if len(tablesDB) == 0 {
		err := fmt.Errorf("some of tables:%v does not exist", tables)
		shared.LogError("Error finding tables", LogDBRepository, "RemoveTables", err, tables)
		return err
	}

	if err := r.db.Model(zone).Association("Tables").Append(tablesDB); err != nil {
		shared.LogError("Error adding tables to zone", LogDBRepository, "AddTables", err, *zone, tables)
		return err
	}

	return nil
}

// RemoveTables method for remove tables to zone in database
func (r *DBRepository) RemoveTables(zone *Zone, tables []uint) error {
	var tablesDB []tablesPKG.Table
	if err := r.db.Where("zone_id = ?", zone.ID).Find(&tablesDB, tables).Error; err != nil {
		shared.LogError("Error finding tables", LogDBRepository, "RemoveTables", err, tables)
		return err
	}

	if len(tablesDB) == 0 {
		err := fmt.Errorf("some of tables:%v does not exist", tables)
		shared.LogError("Error finding tables", LogDBRepository, "RemoveTables", err, tables)
		return err
	}

	if err := r.db.Model(zone).Association("Tables").Delete(tablesDB); err != nil {
		shared.LogError("Error deleting tables from zone", LogDBRepository, "RemoveTables", err, *zone, tables)
		return err
	}

	return nil
}

// Enable method for enable a zone in database
func (r *DBRepository) Enable(zoneID string) (*Zone, error) {
	var zone Zone
	if err := r.db.First(&zone, zoneID).Error; err != nil {
		shared.LogError("Error finding zone", LogDBRepository, "Enable", err, zoneID)
		return nil, err
	}

	if err := r.db.Model(&zone).Update("active", !zone.Active).Error; err != nil {
		shared.LogError("Error updating zone", LogDBRepository, "Enable", err, zone)
		return nil, err
	}

	return &zone, nil
}
