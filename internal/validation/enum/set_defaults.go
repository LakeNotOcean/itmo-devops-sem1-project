package enum

import (
	"reflect"
	"strings"
)

// установка значений по умолчанию для полей с типом enum
func SetDefaults[T any](params *T) {
	v := reflect.ValueOf(params).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		validateTag := fieldType.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		if strings.Contains(validateTag, ENUM_VALIDATOR) {
			parts := strings.Split(validateTag, ENUM_VALIDATOR+"=")
			if len(parts) < 2 {
				panic(getInvalidValidatorFormatMessage(fieldType.Name))
			}

			enumValues := strings.Split(parts[1], ";")

			if len(enumValues) < 2 {
				panic(getInvalidValidatorFormatMessage(fieldType.Name))
			}

			defaultValue := enumValues[0]

			if field.Kind() == reflect.String && field.String() == "" {
				field.SetString(defaultValue)
			}
		}
	}
}
