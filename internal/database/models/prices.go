package models

import (
	"time"
)

type Prices struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	CreateDate time.Time `gorm:"type:date;uniqueIndex:idx_prices_unique" json:"create_date"`
	Name       string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_prices_unique" json:"name"`
	Category   string    `gorm:"type:varchar(255);not null;index;uniqueIndex:idx_prices_unique" json:"category"`
	Price      float64   `gorm:"type:decimal(10,2);not null;uniqueIndex:idx_prices_unique;check:price >= 0" json:"price"`
}

func (Prices) TableName() string {
	return "prices"
}
