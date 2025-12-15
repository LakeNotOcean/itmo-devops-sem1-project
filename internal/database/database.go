package database

import (
	"fmt"
	"log"
	"sem1-final-project-hard-level/internal/config"
	"sem1-final-project-hard-level/internal/database/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDb(cfg *config.Config) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", cfg.DbHost, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.Port)

	var err error
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	err = Db.AutoMigrate(&models.Product{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// на всякий случай
	Db.Exec(`
        ALTER TABLE products 
        ADD CONSTRAINT IF NOT EXISTS price_check CHECK (price >= 0)
    `)

	log.Println("Database connection established and migrated successfully")
	return nil
}

func GetDb() *gorm.DB {
	return Db
}
