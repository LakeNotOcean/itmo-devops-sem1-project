package handlers

import (
	"fmt"
	"net/http"
	custommiddleware "sem1-final-project-hard-level/internal/custom_middlewares"
	"sem1-final-project-hard-level/internal/dto"
)

type PriceHandler struct {
	// ct
}

func NewPriceHandler() *PriceHandler {
	return &PriceHandler{}
}

func (h *PriceHandler) GetPrices(w http.ResponseWriter, r *http.Request) {
	params, err := custommiddleware.GetQueryParamsFromContext[dto.GetPricesQueryParamsDto](r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Используем params
	fmt.Printf("Max: %v, Min: %v, Start: %v, End: %v \n", params.Max, params.Min, params.Start, params.End)
	w.WriteHeader(http.StatusOK)
}

// func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
//     params, err := middleware.GetQueryParamsFromContext[ListUsersQuery](r.Context())
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     // Используем params
//     fmt.Printf("Page: %d, Limit: %d\n", params.Page, params.Limit)
// }
