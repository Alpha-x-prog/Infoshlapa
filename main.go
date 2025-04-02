package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"newsAPI/db" // Замените на актуальный путь к пакету db
	"newsAPI/gemini"
	_ "newsAPI/gemini"
	"newsAPI/parser"
	_ "newsAPI/parser"
	"strconv"
	"strings"
)

func main() {
	// Проверка на наличие API-ключа
	//apiKey := os.Getenv("NEWSDATA_API_KEY")
	apiKey := "pub_77741bf0488e9be1b2e35d8e319363ee57db7"
	if apiKey == "" {
		fmt.Println("API-ключ не найден. Установите переменную окружения NEWSDATA_API_KEY.")
		return
	}

	// Инициализация подключения к БД
	database, err := db.InitDB()
	if err != nil {
		fmt.Println("Ошибка инициализации БД:", err)
		return
	}
	defer database.Close()

	apiURL := "https://newsdata.io/api/1/latest?apikey=%s&category=top&language=ru&country=ru"
	// Запуск парсинга и сохранения новостей в БД
	err = parser.ParseAndSaveNews(apiURL, apiKey, database)
	if err != nil {
		fmt.Println("Ошибка при парсинге и сохранении новостей:", err)
	}
	r := gin.Default()

	// Раздача статических файлов
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	r.GET("/news", func(c *gin.Context) {
		getNews(c, database)
	})

	//помощник
	r.POST("/ask", func(c *gin.Context) {
		geminiASK(c)
	})

	r.Run(":8080")
}

func getNews(c *gin.Context, database *sql.DB) {
	// Получаем параметры запроса
	limit := 15
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	category := c.DefaultQuery("category", "") // Категория из запроса

	// Список валидных категорий
	validCategories := []string{"top", "sports", "technology", "business", "science", "entertainment", "health", "world", "politics", "environment", "food"}

	// Проверка валидности категории
	var validCategory bool
	for _, valid := range validCategories {
		if category == valid {
			validCategory = true
			break
		}
	}

	// Логирование параметров запроса
	log.Printf("Получение новостей с лимитом %d, смещением %d, категорией: %s", limit, offset, category)

	var rows *sql.Rows
	var err error

	// Формируем SQL-запрос в зависимости от переданной категории
	if validCategory {
		// Если категория валидна, фильтруем по ней
		rows, err = database.Query(
			"SELECT article_id, title, link, keywords, creator, video_url, description, content, pub_date, image_url, source_id, source_name, source_url, language, country, category, sentiment "+
				"FROM news WHERE category LIKE ? ORDER BY pub_date DESC LIMIT ? OFFSET ?",
			"%"+category+"%", limit, offset,
		)
	} else {
		// Если категория не передана или не валидна, выводим все новости
		rows, err = database.Query(
			"SELECT article_id, title, link, keywords, creator, video_url, description, content, pub_date, image_url, source_id, source_name, source_url, language, country, category, sentiment "+
				"FROM news ORDER BY pub_date DESC LIMIT ? OFFSET ?",
			limit, offset,
		)
	}

	// Обработка ошибок
	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось выполнить запрос"})
		return
	}
	defer rows.Close()

	var news []db.NewsArticle

	// Перебираем строки результата запроса
	for rows.Next() {
		var n db.NewsArticle
		var keywordsStr, creatorStr, countryStr, categoryStr string

		err := rows.Scan(
			&n.ArticleID, &n.Title, &n.Link, &keywordsStr, &creatorStr, &n.VideoURL,
			&n.Description, &n.Content, &n.PubDate, &n.ImageURL, &n.SourceID, &n.SourceName, &n.SourceURL,
			&n.Language, &countryStr, &categoryStr, &n.Sentiment,
		)

		if err != nil {
			log.Printf("Ошибка при сканировании строки: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обработке данных"})
			return
		}

		// Преобразуем строки с разделителями в срезы
		n.Keywords = strings.Split(keywordsStr, ",")
		n.Creator = strings.Split(creatorStr, ",")
		n.Country = strings.Split(countryStr, ",")
		n.Category = strings.Split(categoryStr, ",")

		// Добавляем новость в список
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

type Request struct {
	Prompt string `json:"prompt"`
}

type Response struct {
	Content string `json:"content"`
}

func geminiASK(c *gin.Context) {
	var req Request

	// Декодируем JSON из тела запроса в структуру Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка декодирования запроса"})
		return
	}
	userQuery := req.Prompt
	// Получаем ответ от функции geminiResponse

	responseContent := gemini.GeminiResponse("Напиши кратко ответ на вопрос: " + userQuery)
	// Отправляем JSON-ответ с полученным ответом
	c.JSON(http.StatusOK, Response{Content: responseContent})
}
