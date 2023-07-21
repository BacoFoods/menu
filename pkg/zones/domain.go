package zones

import (
	"github.com/BacoFoods/menu/pkg/tables"
)

const (
	ErrorZoneFinding      = "Error finding zones"
	ErrorZoneCreating     = "Error creating zone"
	ErrorZoneUpdating     = "Error updating zone"
	ErrorZoneDeleting     = "Error deleting zone"
	ErrorZoneBadRequest   = "Error bad request"
	ErrorZoneNotFound     = "Error zone not found"
	ErrorZoneAddingTables = "Error adding tables to zone"
)

type Zone struct {
	ID      *uint          `json:"id"`
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
	AddTablesToZone(zoneID string, tables []uint) error
}
