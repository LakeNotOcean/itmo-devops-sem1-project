package validation

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type DateFormatValidator struct{}

const (
	DATEFORMAT_VALIDATOR = "dateformat"
	TIMEFORMAT           = "2006-01-02"
)

func (v DateFormatValidator) Name() string {
	return DATEFORMAT_VALIDATOR
}

func (v DateFormatValidator) Register(validate *validator.Validate) error {
	return validate.RegisterValidation(v.Name(), v.Validate)
}

func (DateFormatValidator) Validate(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	if str == "" {
		return true
	}

	_, err := time.Parse(TIMEFORMAT, str)
	return err == nil
}

// автоматическая регистрация валидатора при импорте
func init() {
	Register(DateFormatValidator{})
}

// type PriceQueryParamsDto struct {
//     Start string `query:"start" validate:"omitempty,dateformat"`
//     End   string `query:"end"   validate:"omitempty,dateformat,daterange=Start"`
//     Min   *int   `query:"min"   validate:"omitempty,gt=0"`
//     Max   *int   `query:"max"   validate:"omitempty,gt=0,minmax=Min"`
// }
