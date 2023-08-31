package account

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/BacoFoods/menu/pkg/shared"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

const (
	ErrorBadRequest             = "bad request"
	ErrorAccountPinBadRequest   = "bad request pin must have 4 digits"
	ErrorAccountCreation        = "error creating account"
	ErrorAccountPinCreation     = "error creating pin account"
	ErrorAccountDeleting        = "error deleting account"
	ErrorAccountLogin           = "error login account"
	ErrorAccountPinLogin        = "error login with pin account"
	ErrorAccountInvalidPassword = "error invalid password"
	ErrorAccountFinding         = "error finding account"
)

type Repository interface {
	Create(*Account) (*Account, error)
	Get(username string) (*Account, error)
	Delete(accountID string) error
	Find(filter map[string]any) ([]Account, error)
}

type Account struct {
	Id          uint           `json:"id"`
	DisplayName string         `json:"display_name"`
	Username    string         `json:"username"`
	Password    string         `json:"-" swaggerignore:"true"`
	Email       string         `json:"email"`
	ChannelID   *int64         `json:"channel_id"`
	StoreID     *int64         `json:"store_id"`
	BrandID     *int64         `json:"brand_id"`
	Role        string         `json:"role"`
	Disabled    bool           `json:"disabled"`
	CreatedAt   *time.Time     `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time     `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Role struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	CreatedAt   *time.Time     `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time     `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

func (a *Account) HashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		shared.LogError("error hashing password", "pkg/account/domain", "HashPassword", err, *a)
		return err
	}

	a.Password = string(hashed)
	return nil
}

func (a *Account) CheckPassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password)); err != nil {
		shared.LogWarn("error checking password", "pkg/account/domain", "CheckPassword", err, *a)
		return false
	}

	return true
}

func (a *Account) HashPin() {
	hasher := sha256.New()
	hasher.Write([]byte(a.Password))
	hashBytes := hasher.Sum(nil)
	a.Username = hex.EncodeToString(hashBytes)
	a.Password = hex.EncodeToString(hashBytes)
}
