package tables

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

const (
	ErrorBadRequest        = "error bad request"
	ErrorTableUpdating     = "error updating table"
	ErrorTableDeleting     = "error deleting table"
	ErrorTableCreating     = "error creating table"
	ErrorTableGetting      = "error getting table"
	ErrorTableFinding      = "error finding table"
	ErrorTableHasOrder     = "error order was created but table already has an order"
	ErrorTableReleasing    = "error releasing table"
	ErrorTableNotFound     = "error table not found"
	ErrorTableScanningQR   = "error scanning qr"
	ErrorTableGeneratingQR = "error generating qr"
	ErrorTableIDEmpty      = "error table id empty"
)

const (
	ErrorZoneFinding        = "error finding zones"
	ErrorZoneCreating       = "error creating zone"
	ErrorZoneUpdating       = "error updating zone"
	ErrorZoneDeleting       = "error deleting zone"
	ErrorZoneBadRequest     = "error bad request"
	ErrorZoneNotFound       = "error zone not found"
	ErrorZoneAddingTables   = "error adding tables to zone"
	ErrorZoneRemovingTables = "error removing tables from zone"
	ErrorZoneEnabling       = "error enabling zone"
)

type Zone struct {
	ID      uint    `json:"id"`
	Name    string  `json:"name"`
	Tables  []Table `json:"tables" gorm:"foreignKey:ZoneID"`
	Active  bool    `json:"active" gorm:"default:true"`
	StoreID *uint   `json:"store_id"`
}

type ZoneRepository interface {
	Find(filters map[string]any) ([]Zone, error)
	GetZone(zoneID string) (*Zone, error)
	Create(zone *Zone) (*Zone, error)
	Update(zoneID string, zone *Zone) (*Zone, error)
	Delete(zoneID string) error
	AddTables(zone *Zone, tables []uint) error
	RemoveTables(zone *Zone, tables []uint) error
	Enable(zoneID string) (*Zone, error)
}

func SetTables(zone *Zone, tableNumber, tableAmount int) {
	if tableNumber > 0 && tableAmount > 0 {
		for i := tableNumber; i <= tableNumber+tableAmount; i++ {
			table := Table{
				Number:      i,
				DisplayName: fmt.Sprintf("Mesa %d", i),
				DisplayID:   fmt.Sprintf("%d", i),
				IsActive:    true,
			}
			zone.Tables = append(zone.Tables, table)
		}
	}
}

type Table struct {
	ID          uint            `json:"id,omitempty"`
	DisplayID   string          `json:"display_id" binding:"required"`
	DisplayName string          `json:"display_name" binding:"required"`
	Number      int             `json:"number" binding:"required"`
	XLocation   float64         `json:"xlocation,omitempty"`
	YLocation   float64         `json:"ylocation,omitempty"`
	IsActive    bool            `json:"is_active"`
	ZoneID      *uint           `json:"zone_id"`
	Zone        *Zone           `json:"zone"`
	OrderID     *uint           `json:"order_id"`
	QR          *QR             `json:"qr,omitempty" gorm:"foreignKey:TableID"`
	CreatedAt   *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type QR struct {
	ID        uint   `json:"id,omitempty"`
	TableID   *uint  `json:"table_id" binding:"required"`
	Table     *Table `json:"table,omitempty"`
	DisplayID string `json:"display_id" binding:"required"`
	IsActive  bool   `json:"is_active"`
	URL       string `json:"url"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *gorm.DeletedAt `swaggerignore:"true"`
}

type Repository interface {
	Get(id string) (*Table, error)
	Find(query map[string]any) ([]Table, error)
	Create(table *Table) (*Table, error)
	Update(id string, table *Table) (*Table, error)
	Delete(id string) error
	SetOrder(tableID, orderID *uint) (*Table, error)
	SwapTable(orderID, newTableID uint, oldTableID *uint) (*Table, error)
	RemoveOrder(tableID *uint) (*Table, error)
	ScanQR(qrID string) (*Table, error)
	CreateQR(qr QR) (*QR, error)
}
