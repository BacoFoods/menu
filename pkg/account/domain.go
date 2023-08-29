package account

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

const (
	ErrorBadRequest      = "bad request"
	ErrorAccountCreating = "error creating account"
	ErrorAccountDeleting = "error deleting account"
	ErrorAccountLogin    = "error login account"
)

type Account struct {
	ID        int64          `json:"ID"`
	Username  string         `json:"username"`
	Password  string         `json:"password"`
	Email     string         `json:"email"`
	ChannelID *int64         `json:"channel_id"`
	StoreID   *int64         `json:"store_id"`
	BrandID   *int64         `json:"brand_id"`
	Role      string         `json:"role"`
	Disabled  bool           `json:"disabled"`
	CreatedAt *time.Time     `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Role struct {
	ID          int64          `json:"ID"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	CreatedAt   *time.Time     `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time     `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Repository interface {
	Create(*Account) (*Account, error)
	Login(username, password string) (*Account, error)
	Delete(username string) error
}

func (a *Account) HashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		shared.LogError("error hashing password", "pkg/account/domain", "HashPassword", err, a)
		return err
	}

	a.Password = string(hashed)
	return nil
}

func (a *Account) CheckPassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password)); err != nil {
		shared.LogError("error checking password", "pkg/account/domain", "CheckPassword", err, a)
		return false
	}

	return true
}
