package payment

import (
	"errors"
	"fmt"

	"github.com/BacoFoods/menu/internal"
	"github.com/BacoFoods/menu/pkg/shared"
)

var (
	ErrCreatingPaylot = errors.New("error creating paylot")
)

type Service interface {
	Get(paymentID string) (*Payment, error)
	Find(filter map[string]any) ([]Payment, error)
	Create(payment *Payment) (*Payment, error)
	Update(payment *Payment) (*Payment, error)
	Delete(paymentID string) (*Payment, error)

	FindPaymentMethods(filter map[string]any) ([]PaymentMethod, error)
	GetPaymentMethod(string) (*PaymentMethod, error)
	CreatePaymentMethod(*PaymentMethod) (*PaymentMethod, error)
	UpdatePaymentMethod(*PaymentMethod) (*PaymentMethod, error)
	DeletePaymentMethod(string) (*PaymentMethod, error)

	CreatePaymentWithPaylot(invoiceID uint, total float64, tip float64, customerID *string) (*Payment, error)
}

type service struct {
	repository    Repository
	paylotAdapter PaylotsAPI
}

type PaylotsAPI interface {
	CreatePaylot(req PaylotReq) (*Paylot, error)
	PaylotStatus(paylotID string) (*PaylotStatus, error)
}

type Response[T any] struct {
	Message string `json:"msg"`
	Status  string `json:"status"`
	Data    T      `json:"data"`
}

type PaylotResponse Response[Paylot]

type Paylot struct {
	PaylotID    string `json:"paylot_id"`
	CheckoutURL string `json:"checkout_url"`
}

type PaylotReq struct {
	Reference        string `json:"reference"`
	Country          string `json:"country"`
	Amount           Amount `json:"amount"`
	IssuerID         string `json:"issuer_id"`
	CallbackURL      string `json:"callback_url"`
	WebhookURL       string `json:"webhook_url"`
	WebhookAuthToken string `json:"webhook_auth_token"`
}

type Amount struct {
	Currency string  `json:"currency" default:"COP" validate:"required"`
	Value    float64 `json:"value" validate:"required"`
}

type PaylotStatusResponse Response[PaylotStatus]

type PaylotStatus struct {
	Reference          string   `json:"reference"`
	Country            string   `json:"country"`
	Currency           string   `json:"currency"`
	TotalValue         float64  `json:"total_value"`
	IssuerID           string   `json:"issuer_id"`
	CallbackURL        string   `json:"callback_url,omitempty"`
	Status             string   `json:"status"`
	Observations       string   `json:"observations"`
	ProviderReferences []string `json:"provider_references"`
	Reason             string   `json:"reason"`
}

func NewService(repository Repository, api PaylotsAPI) service {
	return service{repository: repository, paylotAdapter: api}
}

func (s service) createPaylot(invoiceID uint, total float64, customerID *string) (*Paylot, error) {
	// TODO: change redirect url
	redirectUrl := fmt.Sprintf("%s/%d", internal.Config.OITHost, invoiceID)
	customer := ""
	if customerID != nil {
		customer = *customerID
	}
	req := PaylotReq{
		Reference: fmt.Sprint(invoiceID),
		Country:   "CO",
		Amount: Amount{
			Currency: "COP", // TODO!
			Value:    total,
		},
		IssuerID:    customer,
		CallbackURL: redirectUrl,
		WebhookURL:  s.webhookURL(invoiceID),
		// TODO: get auth token
		// WebhookAuthToken: config.Config.Payments.WebhookAuthToken,
	}

	paylot, err := s.paylotAdapter.CreatePaylot(req)
	if err != nil {
		shared.LogError("Error creating paylot for invoice", "payment/service.go", "Checkout", err, invoiceID)
		// return nil, ErrCreatingPaylot
		return nil, ErrCreatingPaylot
	}

	// TODO: payment

	return paylot, nil
}

func (s service) CreatePaymentWithPaylot(invoiceID uint, total, tip float64, customerID *string) (*Payment, error) {
	// TODO: asumimos un solo pago por invoice
	payments, err := s.Find(map[string]any{"invoice_id": invoiceID})
	if err != nil {
		return nil, err
	}

	var lastPayment *Payment
	if (len(payments)) > 0 {
		// check old payments
		for _, payment := range payments {
			// ignore canceled payments
			if payment.Status == "canceled" {
				continue
			}

			// if the invoide already has a paid payment, return it
			// this prevents a new paylot and payment to be created, and we asume the invoice has been paid in full
			// and invoices only have one payment.
			// TODO: This should change when split-the-bill is introduced
			if payment.Status == PaymentStatusPaid {
				return &payment, nil
			}

			// if a payment is pending and has the same value, return it and reuse the paylot
			sameValue := payment.Quantity == total && payment.Tip == tip
			if payment.Status == PaymentStatusPending && sameValue {
				lastPayment = &payment
				continue
			}

			// cancel any other pending payment
			if (payment.Status == PaymentStatusPending && !sameValue) || lastPayment != nil {
				payment.Status = PaymentStatusCanceled
				// TODO: cancel paylot ?
				_, err := s.Update(&payment)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	if lastPayment != nil {
		return lastPayment, nil
	}

	// create paylot
	paylotValue := total + tip
	paylot, err := s.createPaylot(invoiceID, paylotValue, customerID)
	if err != nil {
		return nil, err
	}

	// create payment
	return s.Create(&Payment{
		InvoiceID:   &invoiceID,
		Method:      "PagosBacu", // TODO: payment method category (?) - origin (?)
		Quantity:    total,
		Tip:         tip,
		TotalValue:  paylotValue,
		Code:        paylot.PaylotID,
		Status:      PaymentStatusPending,
		CheckoutURL: &paylot.CheckoutURL,
	})
}

// TODO: change and implement webhook
func (s *service) webhookURL(invoiceID uint) string {
	// return fmt.Sprintf("%s/api/menu/v1/public/checkout/webhook/%d", s.selfHost, orderID)
	return "TODO: unimplemented"
}

func (s service) Get(paymentID string) (*Payment, error) {
	return s.repository.Get(paymentID)
}

func (s service) Find(filter map[string]any) ([]Payment, error) {
	return s.repository.Find(filter)
}

func (s service) Create(payment *Payment) (*Payment, error) {
	return s.repository.Create(payment)
}

func (s service) Update(payment *Payment) (*Payment, error) {
	return s.repository.Update(payment)
}

func (s service) Delete(paymentID string) (*Payment, error) {
	return s.repository.Delete(paymentID)
}

func (s service) FindPaymentMethods(filter map[string]any) ([]PaymentMethod, error) {
	return s.repository.FindPaymentMethods(filter)
}

func (s service) GetPaymentMethod(paymentMethodID string) (*PaymentMethod, error) {
	return s.repository.GetPaymentMethod(paymentMethodID)
}

func (s service) CreatePaymentMethod(paymentMethod *PaymentMethod) (*PaymentMethod, error) {
	return s.repository.CreatePaymentMethod(paymentMethod)
}

func (s service) UpdatePaymentMethod(paymentMethod *PaymentMethod) (*PaymentMethod, error) {
	return s.repository.UpdatePaymentMethod(paymentMethod)
}

func (s service) DeletePaymentMethod(paymentMethodID string) (*PaymentMethod, error) {
	return s.repository.DeletePaymentMethod(paymentMethodID)
}
