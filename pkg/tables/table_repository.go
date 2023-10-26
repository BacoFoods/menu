package tables

import (
	"fmt"

	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
)

const LogRepository = "pkg/tables/repository"

type tableRepository struct {
	db *gorm.DB
}

func NewTableRepository(db *gorm.DB) *tableRepository {
	return &tableRepository{db}
}

func (r tableRepository) Get(id string) (*Table, error) {
	if id == "" {
		shared.LogWarn("error getting table", LogRepository, "Get", shared.ErrorIDEmpty)
		return nil, shared.ErrorIDEmpty
	}

	var table Table

	if err := r.db.First(&table, id).Error; err != nil {
		shared.LogError(ErrorTableGetting, LogRepository, "Get", err, id)
		return nil, err
	}

	return &table, nil
}

func (r tableRepository) Find(query map[string]any) ([]Table, error) {
	var tables []Table

	if err := r.db.Where(query).Preload("QR").Find(&tables).Error; err != nil {
		shared.LogError(ErrorTableFinding, LogRepository, "Find", err, query)
		return nil, err
	}

	return tables, nil
}

func (r tableRepository) Create(table *Table) (*Table, error) {
	if err := r.db.Create(table).Error; err != nil {
		shared.LogError(ErrorTableCreating, LogRepository, "Create", err, *table)
		return nil, err
	}

	return table, nil
}

func (r tableRepository) Update(id string, table *Table) (*Table, error) {
	var tableDB Table
	if err := r.db.First(&tableDB, id).Error; err != nil {
		shared.LogError(ErrorTableUpdating, LogRepository, "Update", err, id, *table)
		return nil, err
	}

	if err := r.db.Model(&tableDB).Updates(table).Error; err != nil {
		shared.LogError(ErrorTableUpdating, LogRepository, "Update", err, id, *table)
		return nil, err
	}

	return &tableDB, nil
}

func (r tableRepository) Delete(id string) error {
	var table Table
	if err := r.db.First(&table, id).Error; err != nil {
		shared.LogError(ErrorTableDeleting, LogRepository, "Delete", err, id)
		return err
	}

	if err := r.db.Delete(&table).Error; err != nil {
		shared.LogError(ErrorTableDeleting, LogRepository, "Delete", err, id, table)
		return err
	}

	return nil
}

func (r tableRepository) SetOrder(tableID, orderID *uint) (*Table, error) {
	var table Table
	if err := r.db.First(&table, tableID).Error; err != nil {
		shared.LogError(ErrorTableUpdating, LogRepository, "SetOrder", err, *tableID, *orderID)
		return nil, err
	}

	if table.OrderID != nil && table.OrderID == orderID {
		return &table, nil
	}

	if table.OrderID != nil && table.OrderID != orderID {
		shared.LogError(ErrorTableHasOrder, LogRepository, "SetOrder", nil, *tableID, *orderID)
		return nil, fmt.Errorf(ErrorTableHasOrder)
	}

	if err := r.db.Model(&table).Update("order_id", orderID).Error; err != nil {
		shared.LogError(ErrorTableUpdating, LogRepository, "SetOrder", err, *tableID, *orderID)
		return nil, err
	}

	return &table, nil
}

func (r tableRepository) RemoveOrder(tableID *uint) (*Table, error) {
	if tableID == nil {
		shared.LogWarn("tableId is null, releasing table", LogRepository, "RemoveOrder", nil, nil)
		return nil, nil
	}

	var table Table
	if err := r.db.First(&table, tableID).Error; err != nil {
		shared.LogError(ErrorTableUpdating, LogRepository, "RemoveOrder", err, *tableID)
		return nil, err
	}

	if table.OrderID == nil {
		return &table, nil
	}

	table.OrderID = nil
	if err := r.db.Model(&table).Update("order_id", nil).Error; err != nil {
		shared.LogError(ErrorTableUpdating, LogRepository, "RemoveOrder", err, *tableID)
		return nil, err
	}

	return &table, nil
}

func (r tableRepository) ScanQR(qrID string) (*Table, error) {
	var qr QR
	q := r.db.Where("display_id = ?", qrID).
		Where("deleted_at is null").
		Where("is_active").
		Preload("Table.").
		First(&qr)

	if err := q.Error; err != nil {
		shared.LogError(ErrorTableScanningQR, LogRepository, "ScanQR", err, qrID)
		return nil, err
	}

	return qr.Table, nil
}

func (r tableRepository) CreateQR(qr QR) (*QR, error) {
	if err := r.db.Create(&qr).Error; err != nil {
		shared.LogError(ErrorTableGeneratingQR, LogRepository, "CreateQR", err, qr)
		return nil, err
	}

	return &qr, nil
}
