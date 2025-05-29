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

	// Создаем таблицу для Telegram каналов пользователей
	createUserChannelsTable := `CREATE TABLE IF NOT EXISTS user_channels (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		channel_url TEXT NOT NULL,
		channel_username TEXT,
		channel_name TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		UNIQUE(user_id, channel_url)
	);`

	// Создаем таблицу для информации о каналах Telegram
	createTelegramChannelsTable := `CREATE TABLE IF NOT EXISTS telegram_channels (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		channel_id INTEGER UNIQUE,
		channel_username TEXT UNIQUE NOT NULL,
		channel_title TEXT,
		last_message_id INTEGER DEFAULT 0,
		added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		is_active BOOLEAN DEFAULT TRUE
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

	_, err = db.Exec(createUserChannelsTable)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(createTelegramChannelsTable)
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

// AddBookmark добавляет закладку
func AddBookmark(db *sql.DB, userID int, articleID string) error {
	_, err := db.Exec(
		`INSERT OR IGNORE INTO bookmarks (user_id, article_id) VALUES (?, ?)`,
		userID, articleID,
	)
	return err
}

// RemoveBookmark удаляет закладку
func RemoveBookmark(db *sql.DB, userID int, articleID string) error {
	_, err := db.Exec(
		`DELETE FROM bookmarks WHERE user_id = ? AND article_id = ?`,
		userID, articleID,
	)
	return err
}

// GetUserBookmarks получает все закладки пользователя
func GetUserBookmarks(db *sql.DB, userID int) ([]NewsArticle, error) {
	rows, err := db.Query(`
		SELECT n.article_id, n.title, n.link, n.keywords, n.creator, n.video_url, 
		       n.description, n.content, n.pub_date, n.image_url, n.source_id, 
		       n.source_name, n.source_url, n.language, n.country, n.category, n.sentiment
		FROM bookmarks b
		JOIN news n ON b.article_id = n.article_id
		WHERE b.user_id = ?
		ORDER BY b.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookmarks []NewsArticle
	for rows.Next() {
		var article NewsArticle
		var keywords, creator, country, category string
		err := rows.Scan(
			&article.ArticleID, &article.Title, &article.Link,
			&keywords, &creator, &article.VideoURL,
			&article.Description, &article.Content, &article.PubDate,
			&article.ImageURL, &article.SourceID, &article.SourceName,
			&article.SourceURL, &article.Language, &country,
			&category, &article.Sentiment,
		)
		if err != nil {
			return nil, err
		}

		// Преобразуем строки в массивы
		article.Keywords = strings.Split(keywords, ", ")
		article.Creator = strings.Split(creator, ", ")
		article.Country = strings.Split(country, ", ")
		article.Category = strings.Split(category, ", ")

		bookmarks = append(bookmarks, article)
	}

	return bookmarks, nil
}

// AddUserChannel добавляет Telegram канал для пользователя
func AddUserChannel(db *sql.DB, userID int, channelURL, channelUsername, channelName string) error {
	_, err := db.Exec(
		`INSERT OR IGNORE INTO user_channels (user_id, channel_url, channel_username, channel_name) VALUES (?, ?, ?, ?)`,
		userID, channelURL, channelUsername, channelName,
	)
	return err
}

// RemoveUserChannel удаляет Telegram канал пользователя
func RemoveUserChannel(db *sql.DB, userID int, channelURL string) error {
	_, err := db.Exec(
		`DELETE FROM user_channels WHERE user_id = ? AND channel_url = ?`,
		userID, channelURL,
	)
	return err
}

// GetUserChannels получает все Telegram каналы пользователя
func GetUserChannels(db *sql.DB, userID int) ([]map[string]string, error) {
	rows, err := db.Query(`
		SELECT channel_url, channel_name, created_at
		FROM user_channels
		WHERE user_id = ?
		ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []map[string]string
	for rows.Next() {
		channel := make(map[string]string)
		var url, name, createdAt string
		err := rows.Scan(&url, &name, &createdAt)
		if err != nil {
			return nil, err
		}
		channel["url"] = url
		channel["name"] = name
		channel["created_at"] = createdAt
		channels = append(channels, channel)
	}

	return channels, nil
}

// DeleteAllUsers удаляет всех пользователей из базы данных
func DeleteAllUsers(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM users")
	return err
}
