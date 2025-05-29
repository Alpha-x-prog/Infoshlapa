package main

import (
	"fmt"
	"log"
	"newsAPI/api"
	"newsAPI/db"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Инициализация подключения к БД
	database, err := db.InitDB()
	if err != nil {
		fmt.Println("Ошибка инициализации БД:", err)
		return
	}
	defer database.Close()

	r := gin.Default()

	// API маршруты
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/news", func(c *gin.Context) {
			api.GetNews(c, database)
		})

		apiGroup.POST("/ask", api.GeminiAsk)

		// Public routes
		apiGroup.POST("/profile", func(c *gin.Context) {
			api.ProfileAuthHandler(c, database)
		})

		// Protected routes (require JWT)
		protected := apiGroup.Group("/protected")
		protected.Use(api.JWTAuthMiddleware())
		{
			protected.GET("/profile", func(c *gin.Context) {
				userID, exists := c.Get("user_id")
				if !exists {
					c.JSON(500, gin.H{"success": false, "message": "User ID not found in context"})
					return
				}
				c.JSON(200, gin.H{"success": true, "message": "Welcome to your profile!", "user_id": userID})
			})

			// Bookmark routes
			protected.POST("/bookmarks", func(c *gin.Context) {
				api.AddBookmark(c, database)
			})

			protected.DELETE("/bookmarks", func(c *gin.Context) {
				api.RemoveBookmark(c, database)
			})

			protected.GET("/bookmarks", func(c *gin.Context) {
				api.GetBookmarks(c, database)
			})
		}
	}

	// Раздача статических файлов из папки shlapa/dist
	r.Static("/js", "./shlapa/dist/js")
	r.Static("/css", "./shlapa/dist/css")
	r.StaticFile("/", "./shlapa/dist/index.html")
	r.StaticFile("/favicon.ico", "./shlapa/dist/favicon.ico")

	// Обработка всех остальных маршрутов для SPA
	r.NoRoute(func(c *gin.Context) {
		c.File("./shlapa/dist/index.html")
	})

	// Запуск сервера
	log.Println("Server starting on http://localhost:8080")
	r.Run(":8080")
}
