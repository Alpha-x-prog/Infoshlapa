package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	DBPath           string
	Port             string
	TelegramBotToken string
	StaticPath       string
	DistPath         string
}

var AppConfig Config

func Init() error {
	// Загружаем переменные окружения из .env файла
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	// Получаем текущую директорию
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %v", err)
	}

	AppConfig = Config{
		DBPath:           filepath.Join(currentDir, getEnv("DB_PATH", "news.db")),
		Port:             getEnv("PORT", "8080"),
		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		StaticPath:       filepath.Join(currentDir, "static"),
		DistPath:         filepath.Join(currentDir, "shlapa/dist"),
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
