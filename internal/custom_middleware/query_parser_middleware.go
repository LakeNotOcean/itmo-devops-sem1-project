package custommiddleware

import (
	"context"
	"fmt"
	"net/http"
	"sem1-final-project-hard-level/internal/validation"

	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
)

type key int

const (
	queryParams key = iota
)

type QueryParserConfig struct {
	Validator *validator.Validate
	Decoder   *form.Decoder
}

// парсинг и валидация query-параметров запроса
func QueryParserMiddleware[T any](config *QueryParserConfig) func(http.Handler) http.Handler {
	if config == nil {
		config = &QueryParserConfig{
			Validator: validator.New(),
			Decoder:   form.NewDecoder(),
		}
	}

	// кастомные валидаторы из пакета validation
	validation.RegisterCustomValidators(config.Validator)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var params T

			// сначала парсим
			if err := config.Decoder.Decode(&params, r.URL.Query()); err != nil {
				http.Error(w, "Invalid query parameters", http.StatusBadRequest)
				return
			}

			// затем ставим значения по умолчанию
			validation.SetDefaults(&params)

			// в конце валидируем итог
			if err := config.Validator.Struct(params); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			ctx := context.WithValue(r.Context(), queryParams, params)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetQueryParamsFromContext[T any](ctx context.Context) (*T, error) {
	value := ctx.Value(queryParams)
	if value == nil {
		return nil, fmt.Errorf("query params not found in context")
	}

	params, ok := value.(T)
	if !ok {
		return nil, fmt.Errorf("invalid query params type")
	}

	return &params, nil
}
