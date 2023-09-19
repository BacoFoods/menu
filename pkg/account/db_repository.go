package account

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	LogDBRepository string = "pkg/account/db_repository"
)

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db}
}

func (r DBRepository) Create(account *Account) (*Account, error) {
	if err := r.db.Save(account).Error; err != nil {
		shared.LogError("error creating account", LogDBRepository, "Create", err, account)
		return nil, err
	}

	return account, nil
}

func (r DBRepository) Get(username string) (*Account, error) {
	if username == "" {
		shared.LogWarn("error getting account", LogDBRepository, "Get", shared.ErrorIDEmpty)
		return nil, shared.ErrorIDEmpty
	}

	var account Account
	if err := r.db.
		Preload(clause.Associations).
		Where("username = ?", username).First(&account).Error; err != nil {
		shared.LogError("error getting account", LogDBRepository, "Get", err, username)
		return nil, err
	}

	return &account, nil
}

func (r DBRepository) Delete(accountID string) error {
	var account Account
	if err := r.db.First(&account, accountID).Error; err != nil {
		shared.LogError("error getting account", LogDBRepository, "Delete", err, accountID)
		return err
	}

	if err := r.db.Delete(&account, account.Id).Error; err != nil {
		shared.LogError("error deleting account", LogDBRepository, "Delete", err, accountID)
		return err
	}

	return nil
}

func (r DBRepository) Find(filter map[string]any) ([]Account, error) {
	var accounts []Account
	if err := r.db.Find(&accounts, filter).Error; err != nil {
		shared.LogError("error getting accounts", LogDBRepository, "Find", err, filter)
		return nil, err
	}

	return accounts, nil
}
