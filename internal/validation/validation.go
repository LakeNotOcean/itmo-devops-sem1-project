package validation

import (
	enum "sem1-final-project-hard-level/internal/validation/enum"
	registry "sem1-final-project-hard-level/internal/validation/registry"
	validations "sem1-final-project-hard-level/internal/validation/validators"
)

const (
	TIMEFORMAT = validations.TIMEFORMAT
)

// функция-обертка для того, чтобы решить проблему с импортом
func SetDefaults[T any](params *T) {
	enum.SetDefaults(params)
}

var (
	RegisterCustomValidators = registry.RegisterCustomValidators
)
