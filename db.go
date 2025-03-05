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

func selectTagBD10(tag string) []string {
	// Открываем подключение к базе данных
	db, err := sql.Open("sqlite3", "news.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Определяем нужный тег
	tag = "Культура" // Здесь можно подставить любое нужное значение

	// Запрос к базе данных
	rows, err := db.Query("SELECT title FROM news WHERE tags = ? ORDER BY id DESC LIMIT 10", tag)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Обход результатов
	var titles []string
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			log.Fatal(err)
		}
		titles = append(titles, title)
	}

	// Проверяем на ошибки при итерации
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Выводим заголовки
	//fmt.Println("Последние 10 новостей с тегом", tag, ":")
	//for _, title := range titles {
	//	fmt.Println("-", title)
	//}
	return titles
}
