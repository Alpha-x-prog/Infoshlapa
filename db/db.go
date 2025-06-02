package db

import (
	"database/sql"
	"fmt"
	"log"
	"newsAPI/collyan"
	"newsAPI/config"
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
	db, err := sql.Open("sqlite3", config.AppConfig.DBPath)
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

	// Создаем таблицу для закладок
	createBookmarksTable := `CREATE TABLE IF NOT EXISTS bookmarks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		article_id TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (article_id) REFERENCES news(article_id),
		UNIQUE(user_id, article_id)
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

	_, err = db.Exec(createBookmarksTable)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Сохранение новости в БД
func SaveToDB(db *sql.DB, article NewsArticle) error {
	// Проверяем наличие изображения
	if strings.TrimSpace(article.ImageURL) == "" {
		fmt.Printf("Пропущена новость без изображения: %s\n", article.Title)
		return nil
	}

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
	log.Printf("DB AddBookmark: Attempting to add bookmark for user %d, article %s", userID, articleID)

	// Check if bookmark already exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM bookmarks WHERE user_id = ? AND article_id = ?)", userID, articleID).Scan(&exists)
	if err != nil {
		log.Printf("DB AddBookmark: Error checking if bookmark exists: %v", err)
		return err
	}
	if exists {
		log.Printf("DB AddBookmark: Bookmark already exists for user %d, article %s", userID, articleID)
		return nil
	}

	result, err := db.Exec(
		`INSERT INTO bookmarks (user_id, article_id) VALUES (?, ?)`,
		userID, articleID,
	)
	if err != nil {
		log.Printf("DB AddBookmark: Error adding bookmark: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("DB AddBookmark: Error getting rows affected: %v", err)
		return err
	}

	log.Printf("DB AddBookmark: Successfully added bookmark for user %d, article %s (rows affected: %d)", userID, articleID, rowsAffected)
	return nil
}

// RemoveBookmark удаляет закладку
func RemoveBookmark(db *sql.DB, userID int, articleID string) error {
	log.Printf("DB RemoveBookmark: Attempting to remove bookmark for user %d, article %s", userID, articleID)

	_, err := db.Exec(
		`DELETE FROM bookmarks WHERE user_id = ? AND article_id = ?`,
		userID, articleID,
	)
	if err != nil {
		log.Printf("DB RemoveBookmark: Error removing bookmark: %v", err)
		return err
	}

	log.Printf("DB RemoveBookmark: Successfully removed bookmark for user %d, article %s", userID, articleID)
	return nil
}

// GetUserBookmarks получает все закладки пользователя
func GetUserBookmarks(db *sql.DB, userID int) ([]NewsArticle, error) {
	log.Printf("DB GetUserBookmarks: Fetching bookmarks for user %d", userID)

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
		log.Printf("DB GetUserBookmarks: Error querying bookmarks: %v", err)
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
			log.Printf("DB GetUserBookmarks: Error scanning bookmark row: %v", err)
			return nil, err
		}

		// Преобразуем строки в массивы
		article.Keywords = strings.Split(keywords, ", ")
		article.Creator = strings.Split(creator, ", ")
		article.Country = strings.Split(country, ", ")
		article.Category = strings.Split(category, ", ")

		bookmarks = append(bookmarks, article)
	}

	log.Printf("DB GetUserBookmarks: Successfully retrieved %d bookmarks for user %d", len(bookmarks), userID)
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

// GetUserChannelMessages получает все сообщения из каналов пользователя
func GetUserChannelMessages(db *sql.DB, userID int) ([]map[string]interface{}, error) {
	rows, err := db.Query(`
		SELECT tm.message_id, tm.message_text, tm.message_date, tm.media_url,
			   tc.channel_username, tc.channel_title, tm.summary
		FROM telegram_messages tm
		JOIN telegram_channels tc ON tm.channel_id = tc.channel_id
		JOIN user_channels uc ON tc.channel_username = uc.channel_username
		WHERE uc.user_id = ?
		ORDER BY tm.message_date DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []map[string]interface{}
	for rows.Next() {
		message := make(map[string]interface{})
		var messageID int
		var text, date, username, title, summary sql.NullString
		var mediaURL sql.NullString
		err := rows.Scan(&messageID, &text, &date, &mediaURL, &username, &title, &summary)
		if err != nil {
			return nil, err
		}
		message["message_id"] = messageID
		message["text"] = text.String
		message["date"] = date.String
		if mediaURL.Valid {
			message["media_url"] = mediaURL.String
		} else {
			message["media_url"] = ""
		}
		message["channel_username"] = username.String
		message["channel_title"] = title.String
		if summary.Valid {
			message["summary"] = summary.String
		} else {
			message["summary"] = ""
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// GetUserChannelMessagesByChannel получает сообщения из конкретного канала пользователя
func GetUserChannelMessagesByChannel(db *sql.DB, userID int, channelUsername string) ([]map[string]interface{}, error) {
	rows, err := db.Query(`
		SELECT tm.message_id, tm.message_text, tm.message_date, tm.media_url,
			   tc.channel_username, tc.channel_title
		FROM telegram_messages tm
		JOIN telegram_channels tc ON tm.channel_id = tc.channel_id
		JOIN user_channels uc ON tc.channel_username = uc.channel_username
		WHERE uc.user_id = ? AND tc.channel_username = ?
		ORDER BY tm.message_date DESC`,
		userID, channelUsername,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []map[string]interface{}
	for rows.Next() {
		message := make(map[string]interface{})
		var messageID int
		var text, date, username, title sql.NullString
		var mediaURL sql.NullString
		err := rows.Scan(&messageID, &text, &date, &mediaURL, &username, &title)
		if err != nil {
			return nil, err
		}
		message["message_id"] = messageID
		message["text"] = text.String
		message["date"] = date.String
		if mediaURL.Valid {
			message["media_url"] = mediaURL.String
		} else {
			message["media_url"] = ""
		}
		message["channel_username"] = username.String
		message["channel_title"] = title.String
		messages = append(messages, message)
	}

	return messages, nil
}

// GetPublicChannelMessages получает сообщения из публичных каналов
func GetPublicChannelMessages(db *sql.DB) ([]map[string]interface{}, error) {
	rows, err := db.Query(`
		SELECT tm.message_id, tm.message_text, tm.message_date, tm.media_url,
			   tc.channel_username, tc.channel_title, tm.summary
		FROM telegram_messages tm
		JOIN telegram_channels tc ON tm.channel_id = tc.channel_id
		WHERE tc.channel_username IN ('priem_mirea', 'mirea_esports', 'mireaprofkom')
		ORDER BY tm.message_date DESC
		LIMIT 15
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []map[string]interface{}
	for rows.Next() {
		message := make(map[string]interface{})
		var messageID int
		var text, date, username, title, summary sql.NullString
		var mediaURL sql.NullString
		err := rows.Scan(&messageID, &text, &date, &mediaURL, &username, &title, &summary)
		if err != nil {
			return nil, err
		}
		message["message_id"] = messageID
		message["text"] = text.String
		message["date"] = date.String
		if mediaURL.Valid {
			message["media_url"] = mediaURL.String
		} else {
			message["media_url"] = ""
		}
		message["channel_username"] = username.String
		message["channel_title"] = title.String
		if summary.Valid {
			message["summary"] = summary.String
		} else {
			message["summary"] = ""
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// SaveNewsArticle сохраняет статью в базу данных
func SaveNewsArticle(db *sql.DB, article *NewsArticle) error {
	log.Printf("DB SaveNewsArticle: Attempting to save article %s", article.ArticleID)

	// Prepare the arrays for storage
	keywords := strings.Join(article.Keywords, ", ")
	creator := strings.Join(article.Creator, ", ")
	country := strings.Join(article.Country, ", ")
	category := strings.Join(article.Category, ", ")

	// Check if article already exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM news WHERE article_id = ?)", article.ArticleID).Scan(&exists)
	if err != nil {
		log.Printf("DB SaveNewsArticle: Error checking if article exists: %v", err)
		return err
	}

	if exists {
		// Update existing article
		_, err = db.Exec(`
			UPDATE news SET 
				title = ?, link = ?, keywords = ?, creator = ?, video_url = ?,
				description = ?, content = ?, pub_date = ?, image_url = ?,
				source_id = ?, source_name = ?, source_url = ?, language = ?,
				country = ?, category = ?, sentiment = ?
			WHERE article_id = ?`,
			article.Title, article.Link, keywords, creator, article.VideoURL,
			article.Description, article.Content, article.PubDate, article.ImageURL,
			article.SourceID, article.SourceName, article.SourceURL, article.Language,
			country, category, article.Sentiment, article.ArticleID,
		)
	} else {
		// Insert new article
		_, err = db.Exec(`
			INSERT INTO news (
				article_id, title, link, keywords, creator, video_url,
				description, content, pub_date, image_url, source_id,
				source_name, source_url, language, country, category, sentiment
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			article.ArticleID, article.Title, article.Link, keywords, creator, article.VideoURL,
			article.Description, article.Content, article.PubDate, article.ImageURL,
			article.SourceID, article.SourceName, article.SourceURL, article.Language,
			country, category, article.Sentiment,
		)
	}

	if err != nil {
		log.Printf("DB SaveNewsArticle: Error saving article: %v", err)
		return err
	}

	log.Printf("DB SaveNewsArticle: Successfully saved article %s", article.ArticleID)
	return nil
}
