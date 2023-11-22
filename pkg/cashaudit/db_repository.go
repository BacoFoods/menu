package cashaudit

import (
	"errors"
	"fmt"
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	LogDBRepository = "pkg/cashaudit/repository"
)

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db}
}

func (r *DBRepository) Get(storeID string) (*CashAudit, error) {
	var cashAudit CashAudit
	if err := r.db.Where("store_id = ?", storeID).First(&cashAudit).Error; err != nil {
		shared.LogError("error getting cash audit", LogDBRepository, "Get", err, storeID)
		return nil, fmt.Errorf(ErrorCashAuditGetting)
	}

	return &cashAudit, nil
}

func (r *DBRepository) Create(cashAudit *CashAudit) (*CashAudit, error) {
	if err := r.db.Save(cashAudit).Error; err != nil {
		shared.LogError("error creating cash audit", LogDBRepository, "Create", err, cashAudit)
		return nil, fmt.Errorf(ErrorCashAuditCreating)
	}

	return cashAudit, nil
}

func (r *DBRepository) GetTodayCashAudit(storeID string) (*CashAudit, error) {
	var cashAudit CashAudit

	if err := r.db.Preload(clause.Associations).
		Where("store_id = ? AND created_at >= NOW() - INTERVAL '1' DAY", storeID).
		First(&cashAudit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		shared.LogError("error getting orders", LogDBRepository, "GetLastDayOrders", err, storeID)
		return nil, err
	}

	return &cashAudit, nil
}
