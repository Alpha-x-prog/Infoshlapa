package main

import (
	"fmt"
	"log"
	"newsAPI/api"
	"newsAPI/db"
	"newsAPI/handlers"
	"newsAPI/middleware"
	"newsAPI/telegram"

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

	// Инициализация обработчиков
	handlers.InitDB(database)

	// Запускаем фоновый процесс получения сообщений
	telegram.StartMessageFetcher()

	// Создаем роутер
	router := gin.Default()

	// Настраиваем CORS
	router.Use(middleware.CORSMiddleware())

	// Группа API
	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/news", func(c *gin.Context) {
			api.GetNews(c, database)
		})

		apiGroup.POST("/ask", api.GeminiAsk)

		// Auth routes
		apiGroup.POST("/register", func(c *gin.Context) {
			api.RegisterHandler(c, database)
		})

		apiGroup.POST("/login", func(c *gin.Context) {
			api.LoginHandler(c, database)
		})

		// Public routes
		public := apiGroup.Group("/public")
		{
			public.GET("/channels", handlers.GetPublicChannels)
			public.GET("/channels/messages", handlers.GetPublicChannelMessages)
		}

		// Protected routes (require JWT)
		protected := apiGroup.Group("/protected")
		protected.Use(middleware.AuthMiddleware())
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

			// Channel routes
			protected.POST("/channels", func(c *gin.Context) {
				api.AddChannel(c, database)
			})

			protected.DELETE("/channels", func(c *gin.Context) {
				api.RemoveChannel(c, database)
			})

			protected.GET("/channels", func(c *gin.Context) {
				api.GetUserChannels(c, database)
			})

			// Channel messages routes
			protected.GET("/channels/messages", func(c *gin.Context) {
				api.GetUserChannelMessages(c, database)
			})

			protected.GET("/channels/:username/messages", func(c *gin.Context) {
				api.GetUserChannelMessagesByChannel(c, database)
			})

			// Admin routes
			protected.DELETE("/users/all", func(c *gin.Context) {
				api.DeleteAllUsers(c, database)
			})

			// Gemini AI settings routes
			protected.PUT("/ai/settings", func(c *gin.Context) {
				api.UpdateGeminiSettings(c, database)
			})
		}
	}

	// Раздача статических файлов из папки shlapa/dist
	router.Static("/js", "./shlapa/dist/js")
	router.Static("/css", "./shlapa/dist/css")
	router.StaticFile("/", "./shlapa/dist/index.html")
	router.StaticFile("/favicon.ico", "./shlapa/dist/favicon.ico")

	// Обработка всех остальных маршрутов для SPA
	router.NoRoute(func(c *gin.Context) {
		c.File("./shlapa/dist/index.html")
	})

	// Запуск сервера
	log.Println("Server starting on http://localhost:8080")
	router.Run(":8080")
}
