package main

import (
	"database/sql"
	"fmt"
	"log"
	"newsAPI/db"
	"newsAPI/parser"
	"os"
	"time"

	"github.com/joho/godotenv"
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
