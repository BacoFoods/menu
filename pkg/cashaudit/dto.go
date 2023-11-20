package cashaudit

type DTOCashAudit struct {
	TotalSellReported     float64 `json:"total_sell"`
	CashIncomesReported   float64 `json:"cash_incomes"`
	OnlineIncomesReported float64 `json:"online_incomes"`
	CardIncomesReported   float64 `json:"card_incomes"`
}

func (dto DTOCashAudit) ToCashAudit() CashAudit {
	return CashAudit{
		TotalSellReported:     dto.TotalSellReported,
		CashIncomesReported:   dto.CashIncomesReported,
		OnlineIncomesReported: dto.OnlineIncomesReported,
		CardIncomesReported:   dto.CardIncomesReported,
	}
}
