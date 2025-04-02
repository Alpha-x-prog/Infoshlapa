package parser

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"newsAPI/db" // Путь к пакету db
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
	PubDate     string   `json:"pubDate"`
	ImageURL    string   `json:"image_url"`
	SourceID    string   `json:"source_id"`
	SourceName  string   `json:"source_name"`
	SourceURL   string   `json:"source_url"`
	Language    string   `json:"language"`
	Country     []string `json:"country"`
	Category    []string `json:"category"`
	Sentiment   string   `json:"sentiment"`
}

type NewsResponse struct {
	Status  string        `json:"status"`
	Results []NewsArticle `json:"results"`
}

// Функция парсинга и сохранения новостей
func ParseAndSaveNews(apiURL string, apiKey string, database *sql.DB) error {
	url := fmt.Sprintf(apiURL, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("ошибка запроса: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	var news NewsResponse
	if err := json.Unmarshal(body, &news); err != nil {
		fmt.Println("Ответ API:", string(body))
		return fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	if news.Status == "error" {
		return fmt.Errorf("ошибка от API")
	}

	fmt.Println("Последние новости из России:")

	// Сохранение новостей в БД
	for _, article := range news.Results {
		// Преобразуем article из типа NewsArticle в тип db.NewsArticle
		dbArticle := db.NewsArticle{
			ArticleID:   article.ArticleID,
			Title:       article.Title,
			Link:        article.Link,
			Keywords:    article.Keywords,
			Creator:     article.Creator,
			VideoURL:    article.VideoURL,
			Description: article.Description,
			Content:     article.Content,
			PubDate:     article.PubDate,
			ImageURL:    article.ImageURL,
			SourceID:    article.SourceID,
			SourceName:  article.SourceName,
			SourceURL:   article.SourceURL,
			Language:    article.Language,
			Country:     article.Country,
			Category:    article.Category,
			Sentiment:   article.Sentiment,
		}

		if err := db.SaveToDB(database, dbArticle); err != nil {
			fmt.Println("Ошибка сохранения в БД:", err)
		}
	}

	return nil
}
