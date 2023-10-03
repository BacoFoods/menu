package payment

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
)

const (
	LogRepository string = "pkg/payment/repository"
)

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db: db}
}

func (r *DBRepository) Get(paymentID string) (*Payment, error) {
	if paymentID == "" {
		err := shared.ErrorIDEmpty
		shared.LogWarn("payment id is empty", LogRepository, "Get", err, paymentID)
		return nil, err
	}

	var payment Payment
	if err := r.db.Where("id = ?", paymentID).First(&payment).Error; err != nil {
		shared.LogError(ErrorPaymentGetting, LogRepository, "Get", err, paymentID)
		return nil, err
	}

	return &payment, nil
}

func (r *DBRepository) Find(filter map[string]any) ([]Payment, error) {
	var payments []Payment
	if err := r.db.Where(filter).Find(&payments).Error; err != nil {
		shared.LogError(ErrorPaymentFinding, LogRepository, "Find", err, filter)
		return nil, err
	}

	return payments, nil
}

func (r *DBRepository) Create(payment *Payment) (*Payment, error) {
	if err := r.db.Create(payment).Error; err != nil {
		shared.LogError(ErrorPaymentCreating, LogRepository, "Create", err, payment)
		return nil, err
	}

	return payment, nil
}

func (r *DBRepository) Update(payment *Payment) (*Payment, error) {
	if err := r.db.Save(payment).Error; err != nil {
		shared.LogError(ErrorPaymentUpdating, LogRepository, "Update", err, payment)
		return nil, err
	}

	return payment, nil
}

func (r *DBRepository) Delete(paymentID string) (*Payment, error) {
	if paymentID == "" {
		err := shared.ErrorIDEmpty
		shared.LogWarn("payment id is empty", LogRepository, "Delete", err, paymentID)
		return nil, err
	}

	payment, err := r.Get(paymentID)
	if err != nil {
		shared.LogError(ErrorPaymentDeleting, LogRepository, "Delete", err, paymentID)
		return nil, err
	}

	if err := r.db.Delete(payment).Error; err != nil {
		shared.LogError(ErrorPaymentDeleting, LogRepository, "Delete", err, paymentID)
		return nil, err
	}

	return payment, nil
}

func (r *DBRepository) FindPaymentMethods(filter map[string]any) ([]PaymentMethod, error) {
	var paymentMethods []PaymentMethod
	if err := r.db.Where(filter).Find(&paymentMethods).Error; err != nil {
		shared.LogError(ErrorPaymentMethodFinding, LogRepository, "FindPaymentMethods", err)
		return nil, err
	}

	return paymentMethods, nil
}

func (r *DBRepository) GetPaymentMethod(paymentMethodID string) (*PaymentMethod, error) {
	if paymentMethodID == "" {
		err := shared.ErrorIDEmpty
		shared.LogWarn("payment method id is empty", LogRepository, "GetPaymentMethod", err, paymentMethodID)
		return nil, err
	}

	var paymentMethod PaymentMethod
	if err := r.db.Where("id = ?", paymentMethodID).First(&paymentMethod).Error; err != nil {
		shared.LogError(ErrorPaymentMethodFinding, LogRepository, "GetPaymentMethod", err, paymentMethodID)
		return nil, err
	}

	return &paymentMethod, nil
}

func (r *DBRepository) CreatePaymentMethod(paymentMethod *PaymentMethod) (*PaymentMethod, error) {
	if err := r.db.Create(paymentMethod).Error; err != nil {
		shared.LogError(ErrorPaymentMethodFinding, LogRepository, "CreatePaymentMethod", err, paymentMethod)
		return nil, err
	}

	return paymentMethod, nil
}

func (r *DBRepository) UpdatePaymentMethod(paymentMethod *PaymentMethod) (*PaymentMethod, error) {
	if err := r.db.Save(paymentMethod).Error; err != nil {
		shared.LogError(ErrorPaymentMethodFinding, LogRepository, "UpdatePaymentMethod", err, paymentMethod)
		return nil, err
	}

	return paymentMethod, nil
}

func (r *DBRepository) DeletePaymentMethod(paymentMethodID string) (*PaymentMethod, error) {
	if paymentMethodID == "" {
		err := shared.ErrorIDEmpty
		shared.LogWarn("payment method id is empty", LogRepository, "DeletePaymentMethod", err, paymentMethodID)
		return nil, err
	}

	paymentMethod, err := r.GetPaymentMethod(paymentMethodID)
	if err != nil {
		shared.LogError(ErrorPaymentMethodFinding, LogRepository, "DeletePaymentMethod", err, paymentMethodID)
		return nil, err
	}

	if err := r.db.Delete(paymentMethod).Error; err != nil {
		shared.LogError(ErrorPaymentMethodFinding, LogRepository, "DeletePaymentMethod", err, paymentMethodID)
		return nil, err
	}

	return paymentMethod, nil
}
