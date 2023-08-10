package zones

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/tables"
)

const (
	ErrorZoneFinding        = "Error finding zones"
	ErrorZoneCreating       = "Error creating zone"
	ErrorZoneUpdating       = "Error updating zone"
	ErrorZoneDeleting       = "Error deleting zone"
	ErrorZoneBadRequest     = "Error bad request"
	ErrorZoneNotFound       = "Error zone not found"
	ErrorZoneAddingTables   = "Error adding tables to zone"
	ErrorZoneRemovingTables = "Error removing tables from zone"
)

type Zone struct {
	ID      uint           `json:"id"`
	Name    string         `json:"name"`
	Tables  []tables.Table `json:"tables" gorm:"foreignKey:ZoneID"`
	StoreID *uint          `json:"store_id"`
}

type Repository interface {
	Find(filters map[string]any) ([]Zone, error)
	GetZone(zoneID string) (*Zone, error)
	Create(zone *Zone) (*Zone, error)
	Update(zoneID string, zone *Zone) (*Zone, error)
	Delete(zoneID string) error
	AddTables(zone *Zone, tables []uint) error
	RemoveTables(zone *Zone, tables []uint) error
}

func SetTables(zone *Zone, tableNumber, tableAmount int) {
	if tableNumber > 0 && tableAmount > 0 {
		for i := tableNumber; i <= tableNumber+tableAmount; i++ {
			table := tables.Table{
				Number:      i,
				DisplayName: fmt.Sprintf("Mesa %d", i),
				DisplayID:   fmt.Sprintf("%d", i),
				IsActive:    true,
			}
			zone.Tables = append(zone.Tables, table)
		}
	}
}
