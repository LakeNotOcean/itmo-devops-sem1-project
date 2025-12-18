package models

import (
	"time"
)

type Prices struct {
	ID         int       `gorm:"type:integer;primaryKey;not null;autoIncrement:false" json:"id"`
	CreateDate time.Time `gorm:"type:date" json:"create_date"`
	Name       string    `gorm:"type:varchar(255);not null" json:"name"`
	Category   string    `gorm:"type:varchar(255);not null;index" json:"category"`
	Price      float64   `gorm:"type:decimal(10,2);not null;check:price >= 0" json:"price"`
}

func (Prices) TableName() string {
	return "prices"
}
