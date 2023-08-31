package account

import "fmt"

type RequestLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RequestLoginPin struct {
	Pin int `json:"pin"`
}

type RequestAccount struct {
	DisplayName string `json:"display_name" binding:"required"`
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Role        string `json:"role"`
	BrandID     int64  `json:"brand_id" binding:"required"`
	StoreID     int64  `json:"store_id" binding:"required"`
	ChannelID   int64  `json:"channel_id"`
}

func (r RequestAccount) ToAccount() *Account {
	displayName := r.DisplayName
	if r.DisplayName == "" {
		displayName = r.Username
	}

	return &Account{
		DisplayName: displayName,
		Username:    r.Username,
		Password:    r.Password,
		Email:       r.Email,
		Role:        r.Role,
		BrandID:     &r.BrandID,
		StoreID:     &r.StoreID,
		ChannelID:   &r.ChannelID,
	}
}

type RequestPinUser struct {
	DisplayName string `json:"display_name"`
	Pin         int64  `json:"pin" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Role        string `json:"role"`
	BrandID     int64  `json:"brand_id" binding:"required"`
	StoreID     int64  `json:"store_id" binding:"required"`
	ChannelID   int64  `json:"channel_id"`
}

func (r RequestPinUser) ToAccount() *Account {
	var displayName = r.DisplayName
	if r.DisplayName == "" {
		displayName = r.Email
	}

	return &Account{
		DisplayName: displayName,
		Username:    fmt.Sprintf("%d", r.Pin),
		Password:    fmt.Sprintf("%d", r.Pin),
		Email:       r.Email,
		Role:        r.Role,
		BrandID:     &r.BrandID,
		StoreID:     &r.StoreID,
		ChannelID:   &r.ChannelID,
	}
}
