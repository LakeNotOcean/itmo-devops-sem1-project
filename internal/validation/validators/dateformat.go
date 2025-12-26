package validations

import (
	"time"

	"sem1-final-project-hard-level/internal/validation/registry"

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
	registry.Register(DateFormatValidator{})
}
