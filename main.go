package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3" // Импорт драйвера SQLite
	"log"
	"net/http"
	"strconv"
)

var db *sql.DB

type News struct {
	ID          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Tags        string `json:"tags" db:"tags"`
	Description string `json:"description" db:"description"`
	URL         string `json:"url" db:"url"`
	URLToImage  string `json:"urlToImage" db:"urlToImage"`
	PublishedAt string `json:"publishedAt" db:"publishedAt"`
}

func main() {
	newsUpdate()
	var err error
	// Подключаемся к базе данных SQLite
	db, err = sql.Open("sqlite3", "./news.db")
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()

	// Проверка подключения к базе данных
	err = db.Ping()
	if err != nil {
		log.Fatal("Ошибка при пинге базы данных:", err)
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
	tag := c.DefaultQuery("tags", "") // Получаем параметр tags

	// Русские теги
	listTagsRU := []string{"общее", "бизнес", "развлечения", "здоровье", "наука", "спорт", "технологии"}

	// Проверяем, является ли переданный tag валидным
	validTag := false
	for _, t := range listTagsRU {
		if tag == t {
			validTag = true
			break
		}
	}

	// Логирование параметров запроса
	log.Printf("Получение новостей с лимитом %d, смещением %d, тегом: %s", limit, offset, tag)

	// SQL-запрос с фильтрацией по тегу (если он передан и валиден)
	var rows *sql.Rows
	var err error

	if validTag {
		rows, err = db.Query(
			"SELECT id, title, tags, description, url, urlToImage, publishedAt FROM news WHERE tags = ? ORDER BY publishedAt DESC LIMIT ? OFFSET ?",
			tag, limit, offset,
		)
	} else {
		rows, err = db.Query(
			"SELECT id, title, tags, description, url, urlToImage, publishedAt FROM news ORDER BY publishedAt DESC LIMIT ? OFFSET ?",
			limit, offset,
		)
	}

	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось выполнить запрос"})
		return
	}
	defer rows.Close()

	var news []News

	// Перебираем строки результата запроса
	for rows.Next() {
		var n News
		err := rows.Scan(&n.ID, &n.Title, &n.Tags, &n.Description, &n.URL, &n.URLToImage, &n.PublishedAt)
		if err != nil {
			log.Printf("Ошибка при сканировании строки: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обработке данных"})
			return
		}
		news = append(news, n)
	}

	// Проверка на ошибки после завершения перебора строк
	if err := rows.Err(); err != nil {
		log.Printf("Ошибка при обработке строк результата: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обработке данных"})
		return
	}

	// Отправляем ответ с данными
	c.JSON(http.StatusOK, news)
}
