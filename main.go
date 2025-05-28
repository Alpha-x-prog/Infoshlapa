package main

import (
	"database/sql"
	"fmt"
	"log"
	"newsAPI/api"
	"newsAPI/db"
	"newsAPI/parser"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Получаем API ключ из переменных окружения
	apiKey := os.Getenv("NEWSDATA_API_KEY")
	if apiKey == "" {
		log.Fatal("NEWSDATA_API_KEY не найден в переменных окружения")
	}

	// Инициализация подключения к БД
	database, err := db.InitDB()
	if err != nil {
		fmt.Println("Ошибка инициализации БД:", err)
		return
	}
	defer database.Close()

	categories := []string{"top", "health", "politics", "sports", "business", "science", "food"}

	// Запускаем горутину для каждого типа категории
	for _, category := range categories {
		go startNewsFetcher(apiKey, category, database)
	}

	r := gin.Default()

	// Раздача статических файлов
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	r.GET("/news", func(c *gin.Context) {
		api.GetNews(c, database)
	})

	// Помощник
	r.POST("/ask", api.GeminiAsk)

	// Public routes
	r.POST("/register", func(c *gin.Context) {
		api.RegisterHandler(c, database)
	})

	r.POST("/login", func(c *gin.Context) {
		api.LoginHandler(c, database)
	})

	// Protected routes (require JWT)
	protected := r.Group("/protected")
	protected.Use(api.JWTAuthMiddleware())
	{
		// Example protected route
		protected.GET("/profile", func(c *gin.Context) {
			userID, exists := c.Get("user_id")
			if !exists {
				c.JSON(500, gin.H{"success": false, "message": "User ID not found in context"})
				return
			}
			c.JSON(200, gin.H{"success": true, "message": "Welcome to your profile!", "user_id": userID})
		})
	}

	r.Run(":8080")
}

func startNewsFetcher(apiKey, category string, database *sql.DB) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	apiURL := fmt.Sprintf("https://newsdata.io/api/1/latest?apikey=%%s&category=%s&language=ru&country=ru", category)

	for {
		log.Printf("Запуск парсинга для категории: %s", category)
		err := parser.ParseAndSaveNews(apiURL, apiKey, database)
		if err != nil {
			log.Printf("Ошибка при парсинге новостей (%s): %v", category, err)
		}
		<-ticker.C
	}
}
