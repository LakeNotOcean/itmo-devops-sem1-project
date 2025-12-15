package handlers

import (
	"log"
	"net/http"
	"sem1-final-project-hard-level/internal/config"
	custommiddleware "sem1-final-project-hard-level/internal/custom_middlewares"
	"sem1-final-project-hard-level/internal/database"
	"sem1-final-project-hard-level/internal/dto"
	"strings"

	"gorm.io/gorm"
)

type PriceHandler struct {
	db          *gorm.DB
	maxFileSize int64
}

func NewPriceHandler(cfg *config.Config) *PriceHandler {
	return &PriceHandler{db: database.GetDb(), maxFileSize: cfg.MaxFileSize}
}

func (h *PriceHandler) GetPrices(w http.ResponseWriter, r *http.Request) {
	params, err := custommiddleware.GetQueryParamsFromContext[dto.GetPricesQueryParamsDto](r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Используем params
	log.Printf("Max: %v, Min: %v, Start: %v, End: %v \n", params.Max, params.Min, params.Start, params.End)
	w.WriteHeader(http.StatusOK)
}

func (h *PriceHandler) UploadPrices(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, h.maxFileSize<<20)

	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "multipart/form-data") {
		http.Error(w, "Expected multipart/form-data", http.StatusBadRequest)
		return
	}

}
