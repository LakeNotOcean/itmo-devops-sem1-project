package prices

import (
	"fmt"
	"sem1-final-project-hard-level/internal/database/models"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// результат валидации пачки записей
type batchResult struct {
	validRecords    []models.Prices
	validCount      int
	duplicatesCount int
}

// валидация пачки записей
func processBatch(tx *gorm.DB, batch [][]string, seenIDs map[int]bool) (*batchResult, error) {
	result := &batchResult{}
	// нужны для проверке в БД на дубликаты
	var recordIDs []int
	recordMap := make(map[int]*models.Prices)

	// парсим записи и собираем ID для проверки дубликатов по ID в БД
	for _, record := range batch {
		parsed, err := parseAndValidateRecord(record)
		// если данные неполные, то просто игнорируем
		if err != nil {
			continue
		}

		// дубликаты в файле
		if seenIDs[parsed.ID] {
			result.duplicatesCount++
			continue
		}

		seenIDs[parsed.ID] = true
		recordIDs = append(recordIDs, parsed.ID)
		recordMap[parsed.ID] = parsed
	}

	// получаем дубликаты в БД
	existingIDs, err := getExistingIDs(tx, recordIDs)
	if err != nil {
		return nil, err
	}

	// оставшиеся валидные записи собираем для вставки
	for _, parsed := range recordMap {
		if existingIDs[parsed.ID] {
			result.duplicatesCount++
			continue
		}

		result.validRecords = append(result.validRecords, *parsed)
		result.validCount++
	}

	return result, nil
}

// парсинг строки CSV в структуру
func parseAndValidateRecord(record []string) (*models.Prices, error) {
	// количество полей проверяем в первую очередь
	if len(record) < 5 {
		return nil, fmt.Errorf("record has less than 5 fields")
	}

	// парсим ID
	id, err := strconv.Atoi(strings.TrimSpace(record[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %v", err)
	}

	// имя
	name := strings.TrimSpace(record[1])
	if len(name) > 255 || name == "" {
		return nil, fmt.Errorf("invalid name")
	}

	// категория
	category := strings.TrimSpace(record[2])
	if len(category) > 255 || category == "" {
		return nil, fmt.Errorf("invalid category")
	}

	// цена - 2 знака после запятой
	priceStr := strings.TrimSpace(record[3])
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return nil, err
	}

	// парсим дату
	createdAt, err := time.Parse("2006-01-02", strings.TrimSpace(record[4]))
	if err != nil {
		return nil, fmt.Errorf("invalid date: %v", err)
	}

	return &models.Prices{
		ID:        id,
		CreatedAt: createdAt,
		Name:      name,
		Category:  category,
		Price:     price,
	}, nil
}

// получение ID в БД
func getExistingIDs(tx *gorm.DB, ids []int) (map[int]bool, error) {
	existingIDs := make(map[int]bool)

	if len(ids) == 0 {
		return existingIDs, nil
	}

	// ID из БД
	var existing []int
	if err := tx.Model(&models.Prices{}).
		Where("id IN ?", ids).
		Pluck("id", &existing).Error; err != nil {
		return nil, err
	}

	for _, id := range existing {
		existingIDs[id] = true
	}

	return existingIDs, nil
}
