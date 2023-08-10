package order

type OrderTDP struct {
	OrderType string           `json:"order_type"`
	BrandID   *uint            `json:"brand_id" binding:"required"`
	StoreID   *uint            `json:"store_id" binding:"required"`
	ChannelID *uint            `json:"channel_id" binding:"required"`
	TableID   *uint            `json:"table_id"`
	Comments  string           `json:"comments"`
	Detail    []OrderDetailTDP `json:"items"`
	Seats     int              `json:"seats"`
}

func (o OrderTDP) ToOrder() Order {
	details := make([]OrderItem, len(o.Detail))
	for i, d := range o.Detail {
		details[i] = d.ToOrderDetail()
	}

	return Order{
		OrderType: o.OrderType,
		BrandID:   o.BrandID,
		StoreID:   o.StoreID,
		ChannelID: o.ChannelID,
		TableID:   o.TableID,
		Comments:  o.Comments,
		Seats:     o.Seats,
		Items:     details,
	}
}

type OrderDetailTDP struct {
	ProductID *uint  `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"`
	Comments  string `json:"comments"`
	Course    string `json:"course"`
}

func (o OrderDetailTDP) ToOrderDetail() OrderItem {
	return OrderItem{
		ProductID: o.ProductID,
		Quantity:  o.Quantity,
		Comments:  o.Comments,
		Course:    o.Course,
	}
}
