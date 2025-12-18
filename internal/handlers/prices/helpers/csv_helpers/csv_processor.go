package csvhelpers

import (
	"encoding/csv"
	"fmt"
	"io"
	"sem1-final-project-hard-level/internal/database/models"
	"sem1-final-project-hard-level/internal/dto"

	"gorm.io/gorm"
)

func ProcessCSV(db *gorm.DB, reader io.Reader, batchSize int) (*dto.UploadPricesResult, error) {
	// лучше выполнять в одной транзакции
	// если что - откат
	tx := db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", tx.Error)
	}
	defer func() {
		// В случае паники также откатываем транзакцию
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	csvReader := csv.NewReader(reader)
	csvReader.Comma = ','
	csvReader.FieldsPerRecord = -1

	result := &dto.UploadPricesResult{}

	// заголовок не интересен
	_, err := csvReader.Read()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to read CSV header: %v", err)
	}

	seenIDs := make(map[int]bool)
	var batch [][]string

	// обработка по пачке
	for {
		batch, err = readCSVBatch(csvReader, batchSize)
		// завершаем, если ошибка или пачка закончилась
		if err != nil && err != io.EOF {
			tx.Rollback()
			return nil, fmt.Errorf("failed to read CSV batch: %v", err)
		}
		if len(batch) == 0 {
			break
		}

		result.TotalCount += len(batch)

		// обработка пачки
		batchResult, err := processBatch(tx, batch, seenIDs)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to process batch: %v", err)
		}

		result.DuplicatesCount += batchResult.duplicatesCount
		result.TotalItems += batchResult.validCount

		if len(batchResult.validRecords) == 0 {
			continue
		}

		// непосредственная вставка валидных данных в БД
		if err := tx.Create(&batchResult.validRecords).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to insert batch: %v", err)
		}
		if err == io.EOF {
			break
		}
	}

	// получаем статистику
	if err := calculateStatistics(tx, result); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to calculate statistics: %v", err)
	}

	// фиксируем транзакцию в конце
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return result, nil
}

// получение пачки строк
func readCSVBatch(reader *csv.Reader, batchSize int) ([][]string, error) {
	var batch [][]string
	for range batchSize {
		record, err := reader.Read()
		if err != nil {
			return batch, err
		}
		batch = append(batch, record)
	}
	return batch, nil
}

// получение статистики
func calculateStatistics(tx *gorm.DB, result *dto.UploadPricesResult) error {
	// общее количество категорий
	if err := tx.Model(&models.Prices{}).Distinct("category").Count(&result.TotalCategories).Error; err != nil {
		return fmt.Errorf("failed to count categories: %v", err)
	}

	// суммарная стоимость
	if err := tx.Model(&models.Prices{}).Select("COALESCE(SUM(price), 0)").Scan(&result.TotalPrice).Error; err != nil {
		return fmt.Errorf("failed to calculate total price: %v", err)
	}

	return nil
}
