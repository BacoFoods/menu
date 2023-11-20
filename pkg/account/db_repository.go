package account

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
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
	if strings.TrimSpace(username) == "" {
		err := fmt.Errorf(ErrorAccountIDEmpty)
		shared.LogWarn("error getting account", LogDBRepository, "Get", err)
		return nil, err
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

func (r DBRepository) GetByUUID(uuid string) (*Account, error) {
	var account Account
	if err := r.db.Where("uuid = ?", uuid).First(&account).Error; err != nil {
		shared.LogError("error getting account", LogDBRepository, "GetByUUID", err, uuid)
		return nil, err
	}

	return &account, nil
}

func (r DBRepository) Update(account *Account) (*Account, error) {
	var accountDB Account
	if err := r.db.First(&accountDB, account.Id).Error; err != nil {
		shared.LogError("error getting account", LogDBRepository, "Update", err, *account)
		return nil, err
	}

	if err := r.db.Model(&accountDB).Updates(account).Error; err != nil {
		shared.LogError("error updating account", LogDBRepository, "Update", err, *account)
		return nil, err
	}

	return account, nil
}

func (r DBRepository) GetByID(id string) (*Account, error) {
	var account Account
	if err := r.db.Where("id = ?", id).First(&account).Error; err != nil {
		shared.LogError("error getting account", LogDBRepository, "GetByID", err, id)
		return nil, err
	}

	return &account, nil
}
