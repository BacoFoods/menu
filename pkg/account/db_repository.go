package account

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
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
	var account Account
	if err := r.db.Where("username = ?", username).First(&account).Error; err != nil {
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

	if err := r.db.Delete(&accountID, accountID).Error; err != nil {
		shared.LogError("error deleting account", LogDBRepository, "Delete", err, accountID)
		return err
	}

	return nil
}
