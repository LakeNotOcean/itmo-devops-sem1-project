package handlers

import (
	"archive/tar"
	"archive/zip"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sem1-final-project-hard-level/internal/config"
	custommiddleware "sem1-final-project-hard-level/internal/custom_middlewares"
	"sem1-final-project-hard-level/internal/database"
	"sem1-final-project-hard-level/internal/database/models"
	"sem1-final-project-hard-level/internal/dto"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type PriceHandler struct {
	db          *gorm.DB
	maxFileSize int64
	tempFileDir string
}

func NewPriceHandler(cfg *config.Config) *PriceHandler {
	return &PriceHandler{db: database.GetDb(), maxFileSize: cfg.MaxFileSize, tempFileDir: cfg.TempFileDir}
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
	// query-параметры для типа файла
	params, err := custommiddleware.GetQueryParamsFromContext[dto.UploadPricesQueryParams](r.Context())
	if err != nil {
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
	var result *dto.UploadPricesResult
	if params.Format.String() == "tar" {
		result, err = h.processTarArchive(tempFile.Name())
	} else if params.Format.String() == "zip" {
		result, err = h.processZipArchive(tempFile.Name())
	} else {
		http.Error(w, "Unsupported archive format", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// результат
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *PriceHandler) processTarArchive(filePath string) (*dto.UploadPricesResult, error) {
	// Открываем tar архив
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	tarReader := tar.NewReader(file)

	// Ищем файл data.csv в архиве
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if header.Name == "data.csv" || filepath.Base(header.Name) == "data.csv" {
			return h.processCSV(tarReader)
		}
	}

	return nil, fmt.Errorf("data.csv not found in archive")
}

func (h *PriceHandler) processZipArchive(filePath string) (*dto.UploadPricesResult, error) {
	// Открываем zip архив
	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, err
	}
	defer zipReader.Close()

	// Ищем файл data.csv в архиве
	for _, f := range zipReader.File {
		if f.Name == "data.csv" || filepath.Base(f.Name) == "data.csv" {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			return h.processCSV(rc)
		}
	}

	return nil, fmt.Errorf("data.csv not found in archive")
}

func (h *PriceHandler) processCSV(reader io.Reader) (*dto.UploadPricesResult, error) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = ','          // Указываем разделитель
	csvReader.FieldsPerRecord = -1 // Разрешаем разное количество полей

	// Пропускаем заголовок
	_, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %v", err)
	}

	var records []models.Prices
	seenIDs := make(map[int]bool)
	existingIDs := make(map[int]bool)
	totalCount := 0
	duplicatesInFile := 0

	// Читаем данные построчно
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			// Пропускаем некорректные строки
			continue
		}

		totalCount++

		// Проверяем количество полей
		if len(record) < 5 {
			duplicatesInFile++ // считаем как дубликат/некорректную запись
			continue
		}

		// Парсим ID
		id, err := strconv.Atoi(strings.TrimSpace(record[0]))
		if err != nil {
			duplicatesInFile++
			continue
		}

		// Проверяем дубликаты в файле
		if seenIDs[id] {
			duplicatesInFile++
			continue
		}
		seenIDs[id] = true

		// Проверяем существование в БД
		var exists bool
		h.db.Model(&models.Prices{}).Select("count(*) > 0").Where("id = ?", id).Find(&exists)
		if exists {
			existingIDs[id] = true
			duplicatesInFile++
			continue
		}

		// Парсим дату
		createdAt, err := time.Parse("2006-01-02", strings.TrimSpace(record[1]))
		if err != nil {
			duplicatesInFile++
			continue
		}

		// Проверяем имя
		name := strings.TrimSpace(record[2])
		if len(name) > 255 || name == "" {
			duplicatesInFile++
			continue
		}

		// Проверяем категорию
		category := strings.TrimSpace(record[3])
		if len(category) > 255 || category == "" {
			duplicatesInFile++
			continue
		}

		// Парсим цену
		price, err := strconv.ParseFloat(strings.TrimSpace(record[4]), 64)
		if err != nil || price < 0 {
			duplicatesInFile++
			continue
		}

		records = append(records, models.Prices{
			ID:        id,
			CreatedAt: createdAt,
			Name:      name,
			Category:  category,
			Price:     price,
		})
	}

	// Вставляем данные пачками
	if len(records) > 0 {
		batchSize := 100 // Размер пачки
		for i := 0; i < len(records); i += batchSize {
			end := i + batchSize
			if end > len(records) {
				end = len(records)
			}

			batch := records[i:end]
			if err := h.db.Create(&batch).Error; err != nil {
				// Пробуем вставить по одной записи при ошибке
				for _, record := range batch {
					if err := h.db.Create(&record).Error; err != nil {
						log.Printf("Failed to insert record with ID %d: %v", record.ID, err)
					}
				}
			}
		}
	}

	// Считаем итоговую статистику
	var result dto.UploadPricesResult
	result.TotalCount = totalCount
	result.DuplicatesCount = duplicatesInFile
	result.TotalItems = len(records)

	// Получаем общее количество категорий
	h.db.Model(&models.Prices{}).Distinct("category").Count(&result.TotalCategories)

	// Получаем суммарную стоимость
	h.db.Model(&models.Prices{}).Select("COALESCE(SUM(price), 0)").Scan(&result.TotalPrice)

	return &result, nil
}
