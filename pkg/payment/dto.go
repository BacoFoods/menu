package payment

type DTOPayment struct {
	InvoiceID uint    `json:"invoice_id" binding:"required"`
	Method    string  `json:"method" binding:"required"`
	Quantity  float32 `json:"quantity" binding:"required"`
}

func (dto DTOPayment) ToPayment() *Payment {
	return &Payment{
		InvoiceID: &dto.InvoiceID,
		Method:    dto.Method,
		Quantity:  dto.Quantity,
	}
}
