package middleware

import (
	"context"
	"net/http"
)

func ValidateParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if errors:=validation
	})
}
