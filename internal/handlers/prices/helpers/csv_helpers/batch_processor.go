package csvhelpers

import (
	"fmt"
	"log"
	"strings"
	"time"

	"sem1-final-project-hard-level/internal/database/models"
	"sem1-final-project-hard-level/internal/handlers/prices/helpers"
	"sem1-final-project-hard-level/internal/validation"
)

// результат валидации пачки записей
type batchResult struct {
	validRecords []models.Prices
}

// валидация пачки записей
func processBatch(batch [][]string) (*batchResult, error) {
	result := &batchResult{}

	// парсим записи
	for _, record := range batch {
		parsed, err := parseAndValidateRecord(record)
		log.Println("Parsed: ", parsed)
		// если данные неполные, то просто игнорируем
		if err != nil {
			continue
		}

		result.validRecords = append(result.validRecords, *parsed)
	}

	return result, nil
}

// парсинг строки CSV в структуру
func parseAndValidateRecord(record []string) (*models.Prices, error) {
	// количество полей проверяем в первую очередь
	if len(record) < 5 {
		return nil, fmt.Errorf("record has less than 5 fields")
	}

	// игнорируем ID, так как используем инкрементацию
	// имя
	name := strings.TrimSpace(record[1])
	if len(name) > 255 || name == "" {
		return nil, fmt.Errorf("invalid name: %v", name)
	}

	// категория
	category := strings.TrimSpace(record[2])
	if len(category) > 255 || category == "" {
		return nil, fmt.Errorf("invalid category: %v", category)
	}

	// цена
	priceStr := strings.TrimSpace(record[3])
	price, err := helpers.ParsePriceWithRegex(priceStr)
	if err != nil {
		return nil, err
	}

	// парсим дату
	createdAt, err := time.Parse(validation.TIMEFORMAT, strings.TrimSpace(record[4]))
	if err != nil {
		return nil, fmt.Errorf("invalid date: %v", err)
	}

	return &models.Prices{
		ID:         0,
		CreateDate: createdAt,
		Name:       name,
		Category:   category,
		Price:      price,
	}, nil
}
