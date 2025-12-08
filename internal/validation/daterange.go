package validation

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type DateRangeValidator struct{}

const (
	DATERANGE_VALIDATOR = "daterange"
)

func (v DateRangeValidator) Name() string {
	return DATERANGE_VALIDATOR
}

func (v DateRangeValidator) Register(validate *validator.Validate) error {
	return validate.RegisterValidation(v.Name(), v.Validate)
}

func (DateRangeValidator) Validate(fl validator.FieldLevel) bool {
	startFieldName := fl.Param()

	if startFieldName == "" {
		return true
	}

	startStr := fl.Parent().FieldByName(startFieldName).String()
	endStr := fl.Field().String()

	start, err1 := time.Parse(TIMEFORMAT, startStr)
	end, err2 := time.Parse(TIMEFORMAT, endStr)

	if err1 != nil || err2 != nil {
		return true
	}

	return !start.After(end)
}

// автоматическая регистрация валидатора при импорте
func init() {
	Register(DateRangeValidator{})
}
