package prices

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"fmt"
	"sem1-final-project-hard-level/internal/database/models"
	"sem1-final-project-hard-level/internal/validation"
	"strconv"
)

// создание архива с ценами
func createPricesArchive(prices []models.Prices, dataFileName string) ([]byte, error) {
	// csv в памяти сервиса, не лучший вариант, но работает
	csvData, err := createCSVData(prices)
	if err != nil {
		return nil, err
	}

	zipBuffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuffer)

	// файл в архиве
	csvFileInZip, err := zipWriter.Create(dataFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create file in ZIP: %v", err)
	}

	// записываем и возвращаем
	if _, err := csvFileInZip.Write(csvData); err != nil {
		return nil, fmt.Errorf("failed to write CSV to ZIP: %v", err)
	}
	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close ZIP: %v", err)
	}
	return zipBuffer.Bytes(), nil
}

// создания csv файла с ценами
func createCSVData(prices []models.Prices) ([]byte, error) {
	csvBuffer := new(bytes.Buffer)
	csvWriter := csv.NewWriter(csvBuffer)

	// создаем заголовок
	if err := csvWriter.Write([]string{"id", "name", "category", "price", "create_date"}); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %v", err)
	}

	// Записываем данные
	for _, price := range prices {
		record := []string{
			strconv.Itoa(price.ID),
			price.Name,
			price.Category,
			strconv.FormatFloat(price.Price, 'f', 2, 64),
			price.CreatedAt.Format(validation.TIMEFORMAT),
		}
		if err := csvWriter.Write(record); err != nil {
			return nil, fmt.Errorf("failed to write CSV record: %v", err)
		}
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return nil, fmt.Errorf("failed to flush CSV: %v", err)
	}

	return csvBuffer.Bytes(), nil
}
