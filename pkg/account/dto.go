package account

type RequestLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RequestAccount struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Role      string `json:"role"`
	BrandID   int64  `json:"brand_id" binding:"required"`
	StoreID   int64  `json:"store_id" binding:"required"`
	ChannelID int64  `json:"channel_id"`
}

func (r RequestAccount) ToAccount() *Account {
	return &Account{
		Username:  r.Username,
		Password:  r.Password,
		Email:     r.Email,
		Role:      r.Role,
		BrandID:   &r.BrandID,
		StoreID:   &r.StoreID,
		ChannelID: &r.ChannelID,
	}
}
