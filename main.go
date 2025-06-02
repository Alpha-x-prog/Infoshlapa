package main

import (
	"fmt"
	"log"
	"newsAPI/api"
	"newsAPI/config"
	"newsAPI/db"
	"newsAPI/handlers"
	"newsAPI/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация конфигурации
	if err := config.Init(); err != nil {
		log.Fatal("Error initializing config:", err)
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

	// Раздача статических файлов
	router.Static("/js", config.AppConfig.DistPath+"/js")
	router.Static("/css", config.AppConfig.DistPath+"/css")
	router.StaticFile("/", config.AppConfig.DistPath+"/index.html")
	router.StaticFile("/favicon.ico", config.AppConfig.DistPath+"/favicon.ico")

	// Обработка всех остальных маршрутов для SPA
	router.NoRoute(func(c *gin.Context) {
		c.File(config.AppConfig.DistPath + "/index.html")
	})

	// Запуск сервера
	log.Printf("Server starting on http://localhost:%s", config.AppConfig.Port)
	router.Run(":" + config.AppConfig.Port)
}
