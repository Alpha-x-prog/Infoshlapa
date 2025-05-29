package main

import (
	"database/sql"
	"fmt"
	"log"
	"newsAPI/db"
	"newsAPI/parser"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/tidwall/gjson"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get API key from environment variables
	apiKey := os.Getenv("NEWSDATA_API_KEY")
	if apiKey == "" {
		log.Fatal("NEWSDATA_API_KEY not found in environment variables")
	}

	// Initialize database connection
	database, err := db.InitDB()
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer database.Close()

	categories := []string{"top", "health", "politics", "sports", "business", "science", "food"}

	// Start news fetcher for each category
	for _, category := range categories {
		go runNewsFetcher(apiKey, category, database)
	}

	// Keep the main goroutine running
	select {}
}

func runNewsFetcher(apiKey, category string, database *sql.DB) {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	apiURL := fmt.Sprintf("https://newsdata.io/api/1/latest?apikey=%%s&category=%s&language=ru&country=ru", category)

	for {
		log.Printf("Starting news fetch for category: %s", category)
		err := parser.ParseAndSaveNews(apiURL, apiKey, database)
		if err != nil {
			log.Printf("Error fetching news (%s): %v", category, err)
		}
		<-ticker.C
	}
}

func processNewsItem(item gjson.Result, db *sql.DB) error {
	// Проверяем наличие изображения
	imageURL := item.Get("image_url").String()
	if !item.Get("image_url").Exists() || strings.TrimSpace(imageURL) == "" {
		fmt.Printf("Пропущена новость без изображения: %s\n", item.Get("title").String())
		return nil // Пропускаем новости без изображений
	}

	// Проверяем наличие обязательных полей
	if !item.Get("title").Exists() || !item.Get("link").Exists() {
		return nil
	}

	// Получаем значения полей
	title := item.Get("title").String()
	link := item.Get("link").String()
	description := item.Get("description").String()
	pubDate := item.Get("pubDate").String()
	category := item.Get("category").String()
	source := item.Get("source").String()

	// Проверяем, существует ли уже такая новость
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM news WHERE link = $1)", link).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking existing news: %v", err)
	}

	if exists {
		return nil // Новость уже существует, пропускаем
	}

	// Добавляем новость в базу данных
	_, err = db.Exec(`
		INSERT INTO news (title, link, description, pub_date, image_url, category, source)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, title, link, description, pubDate, imageURL, category, source)

	if err != nil {
		return fmt.Errorf("error inserting news: %v", err)
	}

	fmt.Printf("Добавлена новость с изображением: %s\n", title)
	return nil
}
