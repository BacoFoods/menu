package account

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/BacoFoods/menu/internal"
	"github.com/BacoFoods/menu/pkg/brand"
	"github.com/BacoFoods/menu/pkg/channel"
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/BacoFoods/menu/pkg/store"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

const (
	ErrorAccountBadRequest      = "bad request"
	ErrorAccountPinBadRequest   = "bad request pin must have 4 digits"
	ErrorAccountCreation        = "error creating account"
	ErrorAccountPinFinding      = "error finding pin account"
	ErrorAccountPinCreation     = "error creating pin account"
	ErrorAccountPinAlreadyExist = "error this mail already has a pin"
	ErrorAccountDeleting        = "error deleting account"
	ErrorAccountLogin           = "error login account"
	ErrorAccountInvalidPassword = "error invalid password"
	ErrorAccountFinding         = "error finding account"
	ErrorAccountUpdating        = "error updating account"
	ErrorAccountGettingByID     = "error getting account by id"
	ErrorAccountIDEmpty         = "account id is empty"

	LogDomain = "pkg/account/domain"
)

type Repository interface {
	Create(*Account) (*Account, error)
	Get(username string) (*Account, error)
	Delete(accountID string) error
	Find(filter map[string]any) ([]Account, error)
	GetByUUID(uuid string) (*Account, error)
	GetByID(id string) (*Account, error)
	Update(*Account) (*Account, error)
}

type Account struct {
	Id          uint             `json:"id"`
	UUID        string           `json:"uuid"`
	DisplayName string           `json:"display_name"`
	Username    string           `json:"username"`
	Password    string           `json:"-" swaggerignore:"true"`
	Email       string           `json:"email"`
	ChannelID   *uint            `json:"channel_id"`
	Channel     *channel.Channel `json:"channel,omitempty"`
	StoreID     *uint            `json:"store_id"`
	Store       *store.Store     `json:"store,omitempty"`
	BrandID     *uint            `json:"brand_id"`
	Brand       *brand.Brand     `json:"brand,omitempty"`
	Role        string           `json:"role"`
	Disabled    bool             `json:"disabled"`
	CreatedAt   *time.Time       `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time       `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   gorm.DeletedAt   `json:"deleted_at,omitempty" swaggerignore:"true"`
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

func (a *Account) JWT() (string, error) {
	channelName := ""
	if a.Channel != nil {
		channelName = a.Channel.Name
	}

	storeName := ""
	if a.Store != nil {
		storeName = a.Store.Name
	}

	brandName := ""
	if a.Brand != nil {
		brandName = a.Brand.Name
	}

	exp := time.Now().Add(time.Hour * 12) // 12 hours expiration
	claims := jwt.MapClaims{
		"account_id":   a.Id,
		"uuid":         a.UUID,
		"name":         a.DisplayName,
		"email":        a.Email,
		"role":         a.Role,
		"channel":      a.ChannelID,
		"channel_name": channelName,
		"store":        a.StoreID,
		"store_name":   storeName,
		"brand":        a.BrandID,
		"brand_name":   brandName,
		"exp":          exp.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey, err := base64.StdEncoding.DecodeString(internal.Config.TokenSecret)
	if err != nil {
		shared.LogError("error decoding jwt key", LogDomain, "JWT", err, internal.Config.TokenSecret)
		return "", err
	}

	tokenString, err := token.SignedString(secretKey) // internal.Config.TokenSecret
	if err != nil {
		shared.LogError("error generating jwt", LogDomain, "JWT", err, *a)
		return "", err
	}

	return tokenString, nil
}
