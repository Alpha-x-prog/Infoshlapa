package api

import (
	"database/sql"
	"log"
	"net/http"
	"newsAPI/gemini"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type NewsArticle struct {
	ArticleID   string   `json:"article_id"`
	Title       string   `json:"title"`
	Link        string   `json:"link"`
	Keywords    []string `json:"keywords"`
	Creator     []string `json:"creator"`
	VideoURL    string   `json:"video_url"`
	Description string   `json:"description"`
	Content     string   `json:"content"`
	PubDate     string   `json:"publishedAt"`
	ImageURL    string   `json:"urlToImage"`
	SourceID    string   `json:"source_id"`
	SourceName  string   `json:"source_name"`
	SourceURL   string   `json:"url"`
	Language    string   `json:"language"`
	Country     string   `json:"country"`
	Tags        string   `json:"tags"`
	Sentiment   string   `json:"sentiment"`
}

type Request struct {
	Prompt string `json:"prompt"`
}

type Response struct {
	Content string `json:"content"`
}

// GetNews обрабатывает запрос на получение новостей
func GetNews(c *gin.Context, database *sql.DB) {
	limit := 15
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	category := c.DefaultQuery("category", "")

	validCategories := []string{"top", "sports", "technology", "business", "science", "entertainment", "health", "world", "politics", "environment", "food"}
	validCategory := false

	for _, valid := range validCategories {
		if category == valid {
			validCategory = true
			break
		}
	}

	log.Printf("Получение новостей с лимитом %d, смещением %d, категорией: %s", limit, offset, category)

	var rows *sql.Rows
	var err error

	if validCategory {
		rows, err = database.Query(
			"SELECT article_id, title, link, keywords, creator, video_url, description, content, pub_date, image_url, source_id, source_name, source_url, language, country, category, sentiment "+
				"FROM news WHERE category LIKE ? ORDER BY pub_date DESC LIMIT ? OFFSET ?",
			"%"+category+"%", limit, offset,
		)
	} else {
		rows, err = database.Query(
			"SELECT article_id, title, link, keywords, creator, video_url, description, content, pub_date, image_url, source_id, source_name, source_url, language, country, category, sentiment "+
				"FROM news ORDER BY pub_date DESC LIMIT ? OFFSET ?",
			limit, offset,
		)
	}

	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось выполнить запрос"})
		return
	}
	defer rows.Close()

	var news []NewsArticle

	for rows.Next() {
		var n NewsArticle
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

		// Разбиваем строки с разделителями в срезы
		n.Keywords = strings.Split(keywordsStr, ",")
		n.Creator = strings.Split(creatorStr, ",")
		n.Country = strings.TrimSpace(strings.Split(countryStr, ",")[0])
		n.Tags = strings.TrimSpace(strings.Split(categoryStr, ",")[0])

		news = append(news, n)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Ошибка при обработке строк результата: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обработке данных"})
		return
	}

	c.JSON(http.StatusOK, news)
}

// GeminiAsk обрабатывает запросы к Gemini API
func GeminiAsk(c *gin.Context) {
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
