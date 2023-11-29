package scheduler

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"time"
)

type RequestSchedule struct {
	ID      uint   `json:"id,omitempty"`
	StoreID *uint  `json:"store_id" validate:"required"`
	BrandID *uint  `json:"brand_id" validate:"required"`
	Day     string `json:"day" validate:"required,oneof=monday tuesday wednesday thursday friday saturday sunday holiday" enums:"monday,tuesday,wednesday,thursday,friday,saturday,sunday,holiday"`
	Opening string `json:"opening" validate:"required,availablehour"`
	Closing string `json:"closing" validate:"required,availablehour"`
	Enable  bool   `json:"enable"`
}

func (r *RequestSchedule) ToSchedule() *Schedule {
	return &Schedule{
		ID:      r.ID,
		StoreID: r.StoreID,
		BrandID: r.BrandID,
		Day:     r.Day,
		Opening: r.Opening,
		Closing: r.Closing,
		Enable:  r.Enable,
	}
}

// RequestValidate DTO function for request validation
func (r *RequestSchedule) RequestValidate() error {
	_en := en.New()
	uni := ut.New(_en, _en)
	trans, _ := uni.GetTranslator("en")

	errMsg := ""

	v := validator.New()
	v.RegisterTranslation("oneof", trans, shared.OneOfValidationTranslator, shared.OneOfValidation)
	v.RegisterTranslation("required", trans, shared.RequiredValidationTranslator, shared.RequiredValidation)
	v.RegisterValidation("availablehour", availableHour)
	v.RegisterTranslation("availablehour", trans, shared.HourValidationTranslator, shared.HourValidation)

	if err := v.Struct(r); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errMsg += err.Translate(trans) + "; "
		}
		return fmt.Errorf(errMsg)
	}

	return nil
}

// availableHour DTO function for hour validation
var availableHour validator.Func = func(fl validator.FieldLevel) bool {
	hour, ok := fl.Field().Interface().(string)
	if ok {
		_, err := time.Parse("15:04", hour)
		if err != nil {
			return false
		}
	}
	return true
}

type RequestEnableStore struct {
	Enable bool `json:"enable" validate:"required"`
}

type ResponseBrand struct {
	Stores []ResponseStore `json:"stores"`
}

type ResponseStore struct {
	ID        uint                `json:"id"`
	Name      string              `json:"name"`
	Open      bool                `json:"open"`
	Schedules []ResponseSchedules `json:"schedules"`
}

type ResponseSchedules struct {
	ID      uint   `json:"id"`
	Day     string `json:"day"`
	StoreID *uint  `json:"store_id"`
	BrandID *uint  `json:"brand_id"`
	Opening string `json:"opening"`
	Closing string `json:"closing"`
	Enable  bool   `json:"enable"`
}

type RequestHoliday struct {
	ID        uint   `json:"id,omitempty"`
	BrandID   *uint  `json:"brand_id" binding:"required"`
	Day       string `json:"day" binding:"required" example:"2021-01-01" format:"2006-01-02"`
	Enable    bool   `json:"enable"`
	CountryID *uint  `json:"country_id" binding:"required"`
}

func (r *RequestHoliday) ToHoliday() (*Holiday, error) {
	day, err := time.Parse("2006-01-02", r.Day)
	if err != nil {
		return nil, err
	}
	return &Holiday{
		ID:        r.ID,
		BrandID:   r.BrandID,
		Day:       day,
		Enable:    r.Enable,
		CountryID: r.CountryID,
	}, err
}
