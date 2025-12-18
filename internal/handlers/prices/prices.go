package prices

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sem1-final-project-hard-level/internal/config"
	custommiddleware "sem1-final-project-hard-level/internal/custom_middleware"
	"sem1-final-project-hard-level/internal/database"
	"sem1-final-project-hard-level/internal/dto"
	archivehelpers "sem1-final-project-hard-level/internal/handlers/prices/helpers/archive_helpers"
	databasehelpers "sem1-final-project-hard-level/internal/handlers/prices/helpers/database_helpers"

	"gorm.io/gorm"
)

type PriceHandler struct {
	db           *gorm.DB
	maxFileSize  int64
	tempFileDir  string
	dataFileName string
	batchSize    int
}

func NewPriceHandler(cfg *config.Config) *PriceHandler {
	return &PriceHandler{db: database.GetDb(),
		maxFileSize:  cfg.MaxFileSize,
		tempFileDir:  cfg.TempFileDir,
		dataFileName: cfg.DataFileName,
		batchSize:    cfg.BatchSize}
}

func (h *PriceHandler) GetPrices(w http.ResponseWriter, r *http.Request) {
	params, err := custommiddleware.GetQueryParamsFromContext[dto.GetPricesQueryParamsDto](r.Context())

	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("GetPrices params - Start: %s, End: %s, Min: %v, Max: %v\n",
		params.Start, params.End, params.Min, params.Max)

	// берем все записи согласно входным параметрам
	// получение данных упрощено по сравнению с загрузкой - отсуствуют ограничения, возможно их стоит поставить на период
	prices, err := databasehelpers.FetchPricesFromDB(h.db, params)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, fmt.Sprintf("Failed to fetch prices: %v", err), http.StatusInternalServerError)
		return
	}
	archiveBytes, err := archivehelpers.CreatePricesArchive(prices, h.dataFileName)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create archive: %v", err), http.StatusInternalServerError)
		return
	}
	archivehelpers.SendArchiveToClient(w, archiveBytes, h.dataFileName)
}

func (h *PriceHandler) UploadPrices(w http.ResponseWriter, r *http.Request) {
	// query-параметры для типа файла
	params, err := custommiddleware.GetQueryParamsFromContext[dto.UploadPricesQueryParams](r.Context())
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = r.ParseMultipartForm(h.maxFileSize << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// загружаем файл и сохраняем во временную директорию
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	tempFile, err := os.CreateTemp(h.tempFileDir, fmt.Sprintf("upload-*.%s", params.Format.String()))
	if err != nil {
		http.Error(w, "Unable to create temp file", http.StatusInternalServerError)
		return
	}
	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	tempFile.Close()

	// после сохранения - обрабатываем загруженный файл в зависимости от его типа
	result, err := handlePricesFile(h.db, tempFile.Name(), h.dataFileName, h.batchSize, params.Format)

	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// результат
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
