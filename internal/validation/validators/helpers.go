package validations

import (
	"reflect"
)

// разыменовывание указателей и возвращание значения типа T или zero value
func Deref[T any](field reflect.Value) (T, bool) {
	var zero T

	if !field.IsValid() {
		return zero, false
	}

	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return zero, false
		}
		field = field.Elem()
	}

	if !field.CanInterface() {
		return zero, false
	}

	value, ok := field.Interface().(T)
	if !ok {
		return zero, false
	}

	return value, true
}
