package tables

import (
	"fmt"
	"strings"

	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
)

const LogRepository = "pkg/tables/table_repository.go"

type tableRepository struct {
	db *gorm.DB
}

func NewTableRepository(db *gorm.DB) *tableRepository {
	return &tableRepository{db}
}

func (r tableRepository) Get(id string) (*Table, error) {
	if strings.TrimSpace(id) == "" {
		err := fmt.Errorf(ErrorTableIDEmpty)
		shared.LogWarn("error getting table", LogRepository, "Get", err)
		return nil, err
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

func (r tableRepository) SwapTable(newTableID, oldTableID, orderID *uint) (*Table, error) {
	var newTable, oldlTable Table
	if err := r.db.First(&oldlTable, oldTableID).Error; err != nil {
		shared.LogError(ErrorTableUpdating, LogRepository, "SwapTable", err, *oldTableID, *orderID)
		return nil, err
	}

	if err := r.db.First(&newTable, newTableID).Error; err != nil {
		shared.LogError(ErrorTableUpdating, LogRepository, "SwapTable", err, *newTableID, *orderID)
		return nil, err
	}

	if oldTableID == newTableID {
		return &newTable, nil
	}

	// if new table has an order, return error
	if newTable.OrderID != nil {
		return nil, fmt.Errorf(ErrorTableHasOrder)
	}

	// if old table order is not the same as the orderID, return error
	if oldlTable.OrderID != nil && oldlTable.OrderID != orderID {
		return nil, fmt.Errorf("error swapping table, order is not the same")
	}

	// start tx
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// set new table order to old table order
		if err := r.db.Model(&newTable).Update("order_id", orderID).Error; err != nil {
			shared.LogError(ErrorTableUpdating, LogRepository, "SwapTable", err, *newTableID, *orderID)
			return err
		}

		// set old table order to nil
		if err := r.db.Model(&oldlTable).Update("order_id", gorm.Expr("NULL")).Error; err != nil {
			shared.LogError(ErrorTableUpdating, LogRepository, "SwapTable", err, *oldTableID, *orderID)
			return err
		}

		return nil
	})

	return &newTable, err
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
		err := fmt.Errorf(ErrorTableIDEmpty)
		shared.LogError("tableId is null, releasing table", LogRepository, "RemoveOrder", err, nil)
		return nil, err
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
	if err := r.db.Model(&table).Update("order_id", gorm.Expr("NULL")).Error; err != nil {
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
		Preload("Table.Zone").
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
