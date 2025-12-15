package models

import (
	"time"
)

type Product struct {
	ID        int       `gorm:"type:serial;primaryKey;not null" json:"id"`
	CreatedAt time.Time `gorm:"type:date;default:CURRENT_DATE" json:"created_at"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Category  string    `gorm:"type:varchar(255);not null;index" json:"category"`
	Price     float64   `gorm:"type:decimal(10,2);not null;check:price >= 0" json:"price"`
}
