package cashaudit

type DTOCashAudit struct {
	CashIncomesReported  float64 `json:"cash_incomes"`
	OtherIncomesReported float64 `json:"other_incomes"`
	CardIncomesReported  float64 `json:"card_incomes"`
}

func (dto DTOCashAudit) ToCashAudit() CashAudit {
	return CashAudit{
		CashIncomesReported:  dto.CashIncomesReported,
		CardIncomesReported:  dto.CardIncomesReported,
		OtherIncomesReported: dto.OtherIncomesReported,
	}
}

type DTOCashAuditConfirmationRequest struct {
	Observations string `json:"observations"`
	CashAuditID  string `json:"cash_audit_id"`
}

type DTOCashAuditCategories struct {
	TotalSell float64       `json:"total_sell"`
	BruteSell float64       `json:"brute_sell"`
	Orders    int           `json:"orders_length"`
	Seats     int           `json:"seats"`
	Variables []DTOVariable `json:"variables"`
	Incomes   DTOIncome     `json:"incomes"`
}

type DTOVariable struct {
	Discounts DTODiscounts `json:"discounts"`
	Tips      DTOTips      `json:"tips"`
}

type DTOIncome struct {
	Cash   float64    `json:"cash"`
	Cards  []DTOCard  `json:"cards"`
	Others []DTOOther `json:"others"`
}

type DTODiscounts struct {
	Name  string  `json:"name"`
	Total float64 `json:"unit"`
}

type DTOTips struct {
	Origin string  `json:"origin"`
	Total  float64 `json:"unit"`
}

type DTOCard struct {
	Origin string  `json:"origin"`
	Total  float64 `json:"unit"`
}

type DTOOther struct {
	Origin string  `json:"origin"`
	Total  float64 `json:"unit"`
}

func ToDTOCashAuditCategories(cashAudit CashAudit) DTOCashAuditCategories {
	cashIncomes := 0.0
	cardIncomes := make([]DTOCard, 0)
	otherIncomes := make([]DTOOther, 0)
	for _, income := range cashAudit.Incomes {
		if income.Type == IncomeTypeTip {
			continue
		}
		switch income.Type {
		case IncomeTypeCash:
			cashIncomes += income.Income
		case IncomeTypeCard:
			cardIncomes = append(cardIncomes, DTOCard{
				Origin: income.Origin,
				Total:  income.Income,
			})
		default:
			otherIncomes = append(otherIncomes, DTOOther{
				Origin: income.Origin,
				Total:  income.Income,
			})
		}
	}

	return DTOCashAuditCategories{
		TotalSell: cashAudit.TotalSell,
		BruteSell: cashAudit.BruteSell,
		Orders:    int(cashAudit.Orders),
		Seats:     int(cashAudit.Eaters),
		Variables: []DTOVariable{
			{
				Discounts: DTODiscounts{
					Name:  "Descuentos",
					Total: cashAudit.TotalDiscounts,
				},
				Tips: DTOTips{
					Origin: "Propinas",
					Total:  cashAudit.TotalTipsInvoices,
				},
			},
		},
		Incomes: DTOIncome{
			Cash:   cashIncomes,
			Cards:  cardIncomes,
			Others: otherIncomes,
		},
	}
}
