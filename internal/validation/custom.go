package validation

import (
	"github.com/go-playground/validator/v10"
)

var registry []Validator

type Validator interface {
	Name() string
	Register(validate *validator.Validate) error
}

func RegisterCustomValidators(v *validator.Validate) error {
	for _, validator := range registry {
		if err := validator.Register(v); err != nil {
			return err
		}
	}
	return nil
}

func Register(v Validator) {
	registry = append(registry, v)
}
