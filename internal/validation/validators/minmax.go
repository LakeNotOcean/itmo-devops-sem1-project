package validations

import (
	"sem1-final-project-hard-level/internal/validation/registry"

	"github.com/go-playground/validator/v10"
)

// валидатор работает с указателями *int
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

	maxVal, maxOk := Deref[int](maxField)
	minVal, minOk := Deref[int](minField)

	if !maxOk || !minOk {
		return true
	}

	return minVal <= maxVal
}

// автоматическая регистрация валидатора при импорте
func init() {
	registry.Register(MinMaxValidator{})
}
