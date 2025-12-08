package validation

import (
	"github.com/go-playground/validator/v10"
)

// Валидатор работает с указателями
type MinMaxValidator struct{}

const (
	MINMAX_VALIDATOR = "minmax"
)

func (v MinMaxValidator) Name() string {
	return MINMAX_VALIDATOR
}

func (v MinMaxValidator) Register(validate *validator.Validate) error {
	return validate.RegisterValidation(v.Name(), v.Validate)
}

func (MinMaxValidator) Validate(fl validator.FieldLevel) bool {
	minFieldName := fl.Param()
	if minFieldName == "" {
		return true
	}

	maxField := fl.Field()
	minField := fl.Parent().FieldByName(minFieldName)

	if !minField.IsValid() || !maxField.IsValid() {
		return true
	}

	if maxField.IsNil() || minField.IsNil() {
		return true
	}

	maxVal := maxField.Elem().Int()
	minVal := minField.Elem().Int()

	return minVal <= maxVal
}

// автоматическая регистрация валидатора при импорте
func init() {
	Register(MinMaxValidator{})
}
