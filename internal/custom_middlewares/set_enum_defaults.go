// custom_middlewares/query_parser_middleware.go
package custommiddleware

import (
	"reflect"
	"sem1-final-project-hard-level/internal/validation"
	"strings"
)

// Устанавливает значения по умолчанию для полей с валидатором enum
func SetEnumDefaults[T any](params *T) {
	v := reflect.ValueOf(params).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		validateTag := fieldType.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		if strings.Contains(validateTag, validation.ENUM_VALIDATOR) {
			parts := strings.Split(validateTag, validation.ENUM_VALIDATOR+"=")
			if len(parts) < 2 {
				panic(validation.GetInvalidValidatorFormatMessage(fieldType.Name))
			}

			enumValues := strings.Split(parts[1], ";")

			if len(enumValues) < 2 {
				panic(validation.GetInvalidValidatorFormatMessage(fieldType.Name))
			}

			defaultValue := enumValues[0]

			if field.Kind() == reflect.String && field.String() == "" {
				field.SetString(defaultValue)
			}
		}
	}
}
