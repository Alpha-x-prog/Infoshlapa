package main

import (
	"database/sql"
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
