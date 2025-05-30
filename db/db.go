package db

import (
	"database/sql"
	"newsAPI/collyan"
	"newsAPI/gemini"
	_ "newsAPI/gemini"
	"strings"
	"time"

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
	PubDate     string   `json:"publishedAt"`
	ImageURL    string   `json:"urlToImage"`
	SourceID    string   `json:"source_id"`
	SourceName  string   `json:"source_name"`
	SourceURL   string   `json:"url"`
	Language    string   `json:"language"`
	Country     []string `json:"country"`
	Category    []string `json:"tags"`
	Sentiment   string   `json:"sentiment"`
}

// Инициализация БД
func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	// Создаем таблицу новостей
	createNewsTable := `CREATE TABLE IF NOT EXISTS news (
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

	// Создаем таблицу пользователей
	createUsersTable := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// Создаем таблицу для истории запросов к AI
	createConversationsTable := `CREATE TABLE IF NOT EXISTS conversations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		question TEXT NOT NULL,
		answer TEXT NOT NULL,
		timestamp INTEGER NOT NULL
	);`

	// Выполняем создание таблиц
	_, err = db.Exec(createNewsTable)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(createUsersTable)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(createConversationsTable)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Сохранение новости в БД
func SaveToDB(db *sql.DB, article NewsArticle) error {
	if article.Description == "" {
		article.Description = collyan.ScrapperCollyan(article.Link)
		article.Description = gemini.GeminiResponse("Сделай краткое описание в 2-3 предолжения: " + article.Description)
	}
	_, err := db.Exec(
		`INSERT OR IGNORE INTO news (article_id, title, link, keywords, creator, video_url, description, content, pub_date, image_url, source_id, source_name, source_url, language, country, category, sentiment)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		article.ArticleID, article.Title, article.Link,
		strings.Join(article.Keywords, ", "),
		strings.Join(article.Creator, ", "),
		article.VideoURL, article.Description, article.Content,
		article.PubDate, article.ImageURL, article.SourceID,
		article.SourceName, article.SourceURL, article.Language,
		strings.Join(article.Country, ", "),
		strings.Join(article.Category, ", "),
		article.Sentiment,
	)
	return err
}

func saveToDBAI(db *sql.DB, question, answer string) error {
	_, err := db.Exec("INSERT INTO conversations (question, answer, timestamp) VALUES (?, ?, ?)", question, answer, time.Now().Unix())
	return err
}
