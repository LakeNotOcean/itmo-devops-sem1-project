package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/gorm/logger"
)

type Config struct {
	Port         int
	IdleTimeout  int
	MaxFileSize  int64
	TempFileDir  string
	DataFileName string
	BatchSize    int
	logLevel     logger.LogLevel

	DbName     string
	DbUser     string
	DbPassword string
	DbPort     int
	DbHost     string
}

func Load() *Config {
	err := godotenv.Load("./configs/.env")
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}
	exe, _ := os.Executable()
	exeDir := filepath.Dir(exe)

	return &Config{
		Port:         getEnv("PORT", 8080),
		IdleTimeout:  getEnv("IDLE_TIMEOUT", 60),
		MaxFileSize:  getEnv[int64]("MAX_FILE_SIZE", 100),
		TempFileDir:  filepath.Join(exeDir, getEnv("TEMP_FILE_DIR", "temp")),
		DataFileName: getEnv("DATA_FILE_NAME", "data"),
		BatchSize:    getEnv("BATCH_SIZE", 500),
		logLevel:     getEnv("LOG_LEVEL", logger.Info),
		DbName:       getEnv("POSTGRES_DB", "postgres"),
		DbUser:       getEnv("POSTGRES_USER", "postgres"),
		DbPassword:   getEnv("POSTGRES_PASSWORD", "password"),
		DbPort:       getEnv("DB_PORT", 5432),
		DbHost:       getEnv("DB_HOST", "localhost"),
	}

}

func getEnv[T string | bool | int | int64 | logger.LogLevel](key string, defaultValue T) T {
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
