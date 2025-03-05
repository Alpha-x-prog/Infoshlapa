package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("news.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	r := gin.Default()

	// Раздача статических файлов
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	r.GET("/news", getNews)

	r.Run(":8080")
}

func getNews(c *gin.Context) {
	limit := 15
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var news []News
	db.Order("publishedAt DESC").Limit(limit).Offset(offset).Find(&news)

	c.JSON(http.StatusOK, news)
}
