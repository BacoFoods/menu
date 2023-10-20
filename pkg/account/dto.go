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
	BrandID     int64  `json:"brand_id"`
	StoreID     int64  `json:"store_id"`
	ChannelID   int64  `json:"channel_id"`
}

func (r RequestAccount) ToAccount() *Account {
	displayName := r.DisplayName
	if r.DisplayName == "" {
		displayName = r.Username
	}

	var brandID *int64
	if r.BrandID != 0 {
		brandID = &r.BrandID
	}

	var storeID *int64
	if r.StoreID != 0 {
		storeID = &r.StoreID
	}

	var channelID *int64
	if r.ChannelID != 0 {
		channelID = &r.ChannelID
	}

	return &Account{
		DisplayName: displayName,
		Username:    r.Username,
		Password:    r.Password,
		Email:       r.Email,
		Role:        r.Role,
		BrandID:     brandID,
		StoreID:     storeID,
		ChannelID:   channelID,
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

	var channelID *int64
	if r.ChannelID != 0 {
		channelID = &r.ChannelID
	}

	return &Account{
		DisplayName: displayName,
		Username:    fmt.Sprintf("%d", r.Pin),
		Password:    fmt.Sprintf("%d", r.Pin),
		Email:       r.Email,
		Role:        r.Role,
		BrandID:     &r.BrandID,
		StoreID:     &r.StoreID,
		ChannelID:   channelID,
	}
}

type RequestAccountUpdate struct {
	ID          uint   `json:"id" binding:"required"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	BrandID     int64  `json:"brand_id"`
	StoreID     int64  `json:"store_id"`
	ChannelID   int64  `json:"channel_id"`
	UUID        string `json:"uuid"`
}

func (r RequestAccountUpdate) ToAccount() *Account {
	var displayName = r.DisplayName
	if r.DisplayName == "" {
		displayName = r.Email
	}

	var brandID *int64
	if r.BrandID != 0 {
		brandID = &r.BrandID
	}

	var storeID *int64
	if r.StoreID != 0 {
		storeID = &r.StoreID
	}

	var channelID *int64
	if r.ChannelID != 0 {
		channelID = &r.ChannelID
	}

	return &Account{
		Id:          r.ID,
		DisplayName: displayName,
		Email:       r.Email,
		Role:        r.Role,
		BrandID:     brandID,
		StoreID:     storeID,
		ChannelID:   channelID,
		UUID:        r.UUID,
	}
}
