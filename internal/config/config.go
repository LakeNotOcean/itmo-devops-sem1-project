package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        int
	IdleTimeout int
	MaxFileSize int64

	DbHost     string
	DbUser     string
	DbPassword string
	DbPort     int
	DbName     string
}

func Load() *Config {
	err := godotenv.Load("./config/env")
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Port:        getEnv("PORT", 3000),
		IdleTimeout: getEnv("IDLE_TIMEOUT", 60),
		MaxFileSize: getEnv[int64]("MAX_FILE_SIZE", 100),
		DbHost:      getEnv("DB_HOST", "localhost"),
		DbUser:      getEnv("DB_USER", "user"),
		DbPassword:  getEnv("DB_PASSWORD", "password"),
		DbPort:      getEnv("DB_PORT", 5432),
		DbName:      getEnv("DB_NAME", "postgres"),
	}

}

func getEnv[T string | bool | int | int64](key string, defaultValue T) T {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	var result T
	switch any(result).(type) {
	case string:
		return any(value).(T)
	case bool:
		if boolValue, err := strconv.ParseBool(value); err != nil {
			return any(boolValue).(T)
		}
	case int:
		if intValue, err := strconv.Atoi(value); err != nil {
			return any(intValue).(T)
		}
	default:
	}
	return defaultValue
}
