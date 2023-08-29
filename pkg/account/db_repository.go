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

func (r DBRepository) Login(username, password string) (*Account, error) {
	var account Account
	if err := r.db.Where("username = ? AND password = ?", username, password).First(&account).Error; err != nil {
		shared.LogError("error login account", LogDBRepository, "Login", err, username, password)
		return nil, err
	}

	return &account, nil
}
