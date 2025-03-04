package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Обработчик запроса: получение 9 последних новостей
func getNews(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "ничего ты не получишь")
	//db, err := setupDatabase()
	//if err != nil {
	//	log.Fatal("Ошибка подключения к БД:", err)
	//}
	//
	//var news []News
	//result := db.Order("id DESC").Limit(9).Find(&news) // Берём последние 9 записей
	//if result.Error != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
	//	return
	//}
	//
	//c.JSON(http.StatusOK, news)
}

func main() {
	// Создаём новый сервер
	r := gin.Default()

	// Разрешаем CORS, чтобы можно было делать запросы из браузера
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Next()
	})

	// Роут для получения новостей
	r.GET("/api/news", getNews)

	// Запуск сервера на порту 8080
	fmt.Println("Сервер запущен на http://127.0.0.1:8080")
	newsUpdate()
	r.Run(":8080")
}
