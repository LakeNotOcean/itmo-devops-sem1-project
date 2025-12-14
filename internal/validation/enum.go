// validation/enum_with_default.go
package validation

import (
	"fmt"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Используется в формате enum_value=default_value|value1|value2
// default_value устанавливается при отсутствии значения
// Работает для указателей
type EnumValidator struct{}

const (
	ENUM_VALIDATOR = "enum_value"
)

func GetInvalidValidatorFormatMessage(fieldName string) string {
	return fmt.Sprintf("Must be in the format enum_value=default_value|value1|value2 for field \"%s\"", fieldName)
}

func (v EnumValidator) Name() string {
	return ENUM_VALIDATOR
}

func (v EnumValidator) Register(validate *validator.Validate) error {
	return validate.RegisterValidation(v.Name(), v.Validate)
}

func (EnumValidator) Validate(fl validator.FieldLevel) bool {
	params := fl.Param()
	if params == "" {
		return true
	}

	parts := strings.Split(params, ";")
	if len(parts) < 2 {
		panic(GetInvalidValidatorFormatMessage(fl.FieldName()))
	}

	allowedValues := parts[1:]

	str := fl.Field().String()

	return str == "" || slices.Contains(allowedValues, str)
}

// автоматическая регистрация валидатора при импорте
func init() {
	Register(EnumValidator{})
}
