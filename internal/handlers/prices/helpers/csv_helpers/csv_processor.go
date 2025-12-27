package csvhelpers

import (
	"encoding/csv"
	"fmt"
	"io"

	"sem1-final-project-hard-level/internal/database/models"
	"sem1-final-project-hard-level/internal/dto"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// обработка csv-файла с ценами
func ProcessCSV(db *gorm.DB, reader io.Reader, batchSize int) (*dto.UploadPricesResult, error) {
	// лучше выполнять в одной транзакции
	// если проблема - откат
	tx := db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", tx.Error)
	}
	defer func() {
		// в случае паники также откатываем транзакцию
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

	// обработка по пачке
	for {

		// получение пачки
		batch, err := readCSVBatch(csvReader, batchSize)
		// завершаем, если ошибка или пачка закончилась
		if err != nil && err != io.EOF {
			tx.Rollback()
			return nil, fmt.Errorf("failed to read CSV batch: %v", err)
		}
		if len(batch) == 0 {
			break
		}

		result.TotalCount += len(batch)

		// обработка пачки, валидация данных
		batchResult, err := processBatch(batch)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to process batch: %v", err)
		}

		if len(batchResult.validRecords) == 0 {
			continue
		}

		// непосредственная вставка валидных данных в БД
		insertResult := tx.Clauses(clause.OnConflict{Columns: []clause.Column{
			{Name: "create_date"},
			{Name: "name"},
			{Name: "category"},
			{Name: "price"},
		}, DoNothing: true}).Create(&batchResult.validRecords)
		if insertResult.Error != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to insert batch: %v", insertResult.Error)
		}

		// после вставки - учитываем результат
		result.TotalItems += int(insertResult.RowsAffected)
		// дубликаты не должны были быть добавлены
		result.DuplicatesCount += len(batchResult.validRecords) - int(insertResult.RowsAffected)

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
	// получение общего количества категорий и суммарной стоимости
	row := tx.Model(&models.Prices{}).Select("COUNT(DISTINCT category) as total_categories, COALESCE(SUM(price), 0) as total_price").Row()

	if err := row.Scan(&result.TotalCategories, &result.TotalPrice); err != nil {
		return fmt.Errorf("failed to calculate statistics: %v", err)
	}

	return nil
}
