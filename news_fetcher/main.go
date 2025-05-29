package main

import (
	"database/sql"
	"fmt"

	"github.com/tidwall/gjson"
)

func processNewsItem(item gjson.Result, db *sql.DB) error {
	// Проверяем наличие изображения
	if !item.Get("image_url").Exists() || item.Get("image_url").String() == "" {
		return nil // Пропускаем новости без изображений
	}

	// Проверяем наличие обязательных полей
	if !item.Get("title").Exists() || !item.Get("link").Exists() {
		return nil
	}

	// Получаем значения полей
	title := item.Get("title").String()
	link := item.Get("link").String()
	description := item.Get("description").String()
	pubDate := item.Get("pubDate").String()
	imageURL := item.Get("image_url").String()
	category := item.Get("category").String()
	source := item.Get("source").String()

	// Проверяем, существует ли уже такая новость
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM news WHERE link = $1)", link).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking existing news: %v", err)
	}

	if exists {
		return nil // Новость уже существует, пропускаем
	}

	// Добавляем новость в базу данных
	_, err = db.Exec(`
        INSERT INTO news (title, link, description, pub_date, image_url, category, source)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, title, link, description, pubDate, imageURL, category, source)

	if err != nil {
		return fmt.Errorf("error inserting news: %v", err)
	}

	return nil
}
