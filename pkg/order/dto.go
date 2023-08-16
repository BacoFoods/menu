package order

type OrderDTO struct {
	OrderType string         `json:"order_type"`
	BrandID   *uint          `json:"brand_id" binding:"required"`
	StoreID   *uint          `json:"store_id" binding:"required"`
	ChannelID *uint          `json:"channel_id" binding:"required"`
	TableID   *uint          `json:"table_id"`
	Comments  string         `json:"comments"`
	Items     []OrderItemDTO `json:"items"`
	Seats     int            `json:"seats"`
}

func (o OrderDTO) ToOrder() Order {
	items := make([]OrderItem, len(o.Items))
	for i, d := range o.Items {
		items[i] = d.ToOrderItem()
	}

	return Order{
		OrderType: o.OrderType,
		BrandID:   o.BrandID,
		StoreID:   o.StoreID,
		ChannelID: o.ChannelID,
		TableID:   o.TableID,
		Comments:  o.Comments,
		Seats:     o.Seats,
		Items:     items,
	}
}

type OrderItemDTO struct {
	ProductID *uint              `json:"product_id" binding:"required"`
	Comments  string             `json:"comments"`
	Course    string             `json:"course"`
	Modifiers []OrderModifierDTO `json:"modifiers"`
}

func (o OrderItemDTO) ToOrderItem() OrderItem {
	modifiers := make([]OrderModifier, len(o.Modifiers))
	for i, d := range o.Modifiers {
		modifiers[i] = d.ToOrderModifier()
	}

	return OrderItem{
		ProductID: o.ProductID,
		Comments:  o.Comments,
		Course:    o.Course,
		Modifiers: modifiers,
	}
}

type OrderModifierDTO struct {
	Comments  string `json:"comments"`
	ProductID uint   `json:"product_id" binding:"required"`
}

func (o OrderModifierDTO) ToOrderModifier() OrderModifier {
	productID := o.ProductID
	return OrderModifier{
		Comments:  o.Comments,
		ProductID: &productID,
	}
}
