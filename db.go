package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "news.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}
func InsertNews(db *sql.DB, article News) error {
	query := `INSERT INTO news (title, tags, description, url, urlToImage, publishedAt) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, article.Title, article.Tags, article.Description, article.URL, article.URLToImage, article.PublishedAt)
	if err != nil {
		log.Println("Error inserting news:", err)
		return err
	}
	return nil
}

func checkSimilarTitle(db *sql.DB, title string, tag string) (bool, error) {
	query := `
		SELECT COUNT(*) FROM (
			SELECT title FROM news
			WHERE tags = ?
			ORDER BY id DESC
			LIMIT 20
		) WHERE title = ?
	`
	var count int
	err := db.QueryRow(query, tag, title).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("ошибка запроса: %v", err)
	}
	//query := "SELECT COUNT(*) FROM news WHERE category = ? AND title = ?"
	//err := db.QueryRow(query, category, title).Scan(&title)
	//if err != nil {
	//	if err == sql.ErrNoRows {
	//		fmt.Println("В таблице news нет данных.")
	//	} else {
	//		log.Fatalf("Ошибка запроса: %v", err)
	//	}
	//} else {
	//	fmt.Printf("Успешно! Заголовок первой новости: %s\n", title)
	//}
	return count > 0, nil
}
