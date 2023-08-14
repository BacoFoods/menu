package tables

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
)

const LogRepository = "pkg/tables/repository"

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) DBRepository {
	return DBRepository{db}
}

func (r DBRepository) Get(id string) (*Table, error) {
	var table Table

	if err := r.db.First(&table, id).Error; err != nil {
		shared.LogError(ErrorTableGetting, LogRepository, "Get", err, id)
		return nil, err
	}

	return &table, nil
}

func (r DBRepository) Find(query map[string]any) ([]Table, error) {
	var tables []Table

	if err := r.db.Where(query).Find(&tables).Error; err != nil {
		shared.LogError(ErrorTableFinding, LogRepository, "Find", err, query)
		return nil, err
	}

	return tables, nil
}

func (r DBRepository) Create(table *Table) (*Table, error) {
	if err := r.db.Create(table).Error; err != nil {
		shared.LogError(ErrorTableCreating, LogRepository, "Create", err, *table)
		return nil, err
	}

	return table, nil
}

func (r DBRepository) Update(id string, table *Table) (*Table, error) {
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

func (r DBRepository) Delete(id string) error {
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

func (r DBRepository) SetOrder(tableID, orderID *uint) (*Table, error) {
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

func (r DBRepository) RemoveOrder(tableID *uint) (*Table, error) {
	var table Table
	if err := r.db.First(&table, tableID).Error; err != nil {
		shared.LogError(ErrorTableUpdating, LogRepository, "RemoveOrder", err, *tableID)
		return nil, err
	}

	if table.OrderID == nil {
		return &table, nil
	}

	if err := r.db.Model(&table).Update("order_id", nil).Error; err != nil {
		shared.LogError(ErrorTableUpdating, LogRepository, "RemoveOrder", err, *tableID)
		return nil, err
	}

	return &table, nil
}
