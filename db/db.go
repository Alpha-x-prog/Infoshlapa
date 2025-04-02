package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

const dbFile = "news.db"

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

// Инициализация БД
func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	createTable := `CREATE TABLE IF NOT EXISTS news (
		article_id TEXT PRIMARY KEY,
		title TEXT,
		link TEXT,
		keywords TEXT,
		creator TEXT,
		video_url TEXT,
		description TEXT,
		content TEXT,
		pub_date TEXT,
		image_url TEXT,
		source_id TEXT,
		source_name TEXT,
		source_url TEXT,
		language TEXT,
		country TEXT,
		category TEXT,
		sentiment TEXT
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Сохранение новости в БД
func SaveToDB(db *sql.DB, article NewsArticle) error {
	_, err := db.Exec(
		`INSERT OR IGNORE INTO news (article_id, title, link, keywords, creator, video_url, description, content, pub_date, image_url, source_id, source_name, source_url, language, country, category, sentiment)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		article.ArticleID, article.Title, article.Link, fmt.Sprintf("%v", article.Keywords), fmt.Sprintf("%v", article.Creator), article.VideoURL,
		article.Description, article.Content, article.PubDate, article.ImageURL, article.SourceID, article.SourceName,
		article.SourceURL, article.Language, fmt.Sprintf("%v", article.Country), fmt.Sprintf("%v", article.Category), article.Sentiment,
	)
	return err
}
