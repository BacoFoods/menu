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
	items := make([]OrderItem, 0)
	for _, d := range o.Items {
		for j := 0; j < d.Quantity; j++ {
			items = append(items, d.ToOrderItem())
		}
		continue
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
	Quantity  int                `json:"quantity" binding:"required"`
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

type RequestUpdateOrderSeats struct {
	Seats int `json:"seats" binding:"required"`
}

type RequestUpdateOrderStatus struct {
	Status string `json:"status" binding:"required" enum:"ordering,cooking,delivered,invoicing,canceled,completed" example:"ordering,cooking,delivered,invoicing,canceled,completed"`
}

type RequestUpdateOrderProduct struct {
	Price    float64 `json:"price"`
	Unit     string  `json:"unit"`
	Quantity int     `json:"quantity"`
	Comments string  `json:"comments"`
	Course   string  `json:"course"`
}

type RequestUpdateOrderItem struct {
	Price    float64 `json:"price"`
	Comments string  `json:"comments"`
	Course   string  `json:"course"`
}

func (r RequestUpdateOrderItem) ToOrderItem() OrderItem {
	return OrderItem{
		Price:    r.Price,
		Comments: r.Comments,
		Course:   r.Course,
	}
}

type RequestAddProducts struct {
	Items []OrderItemDTO `json:"items" binding:"required"`
}

type RequestModifiers struct {
	Modifiers []OrderModifierDTO `json:"modifiers" binding:"required"`
}

type RequestUpdateOrderComments struct {
	Comments string `json:"comments" binding:"required"`
}
