package databasehelpers

import (
	"sem1-final-project-hard-level/internal/database/models"
	"sem1-final-project-hard-level/internal/dto"
	"sem1-final-project-hard-level/internal/validation"
	"time"

	"gorm.io/gorm"
)

// список цен согласно фильтрам
func FetchPricesFromDB(db *gorm.DB, params *dto.GetPricesQueryParamsDto) ([]models.Prices, error) {
	query := buildQuery(db, params)

	var prices []models.Prices
	if err := query.Find(&prices).Error; err != nil {
		return nil, err
	}

	return prices, nil
}

func buildQuery(db *gorm.DB, params *dto.GetPricesQueryParamsDto) *gorm.DB {
	query := db.Model(&models.Prices{})

	// фильтр по дате начала
	if params.Start != "" {
		startTime, err := time.Parse(validation.TIMEFORMAT, params.Start)
		if err == nil {
			query = query.Where("create_date >= ?", startTime)
		}
	}

	// фильтр по дате окончания
	if params.End != "" {
		endTime, err := time.Parse(validation.TIMEFORMAT, params.End)
		if err == nil {
			query = query.Where("create_date <= ?", endTime)
		}
	}

	// фильтр по минимальной цене
	if params.Min != nil {
		query = query.Where("price >= ?", *params.Min)
	}

	// фильтр по максимальной цене
	if params.Max != nil {
		query = query.Where("price <= ?", *params.Max)
	}

	return query
}
