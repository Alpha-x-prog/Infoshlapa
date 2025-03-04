package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type News struct {
	ID          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Tags        string `json:"tags" db:"tags"`
	Description string `json:"description" db:"description"`
	URL         string `json:"url" db:"url"`
	URLToImage  string `json:"urlToImage" db:"urlToImage"`
	PublishedAt string `json:"publishedAt" db:"publishedAt"`
}

type NewsResponse struct {
	Articles []News `json:"articles"`
}

func main() {
	apiKeyNews := "8edd255699aa4fdfa562f51ac68de15e"
	url := "https://newsapi.org/v2/top-headlines?country=us&apiKey=" + apiKeyNews

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка запроса:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var news NewsResponse
	json.Unmarshal(body, &news)

	for _, article := range news.Articles {
		fmt.Println("Заголовок:", article.Title)
		fmt.Println("Ссылка:", article.URL)
		fmt.Println("Картинка:", article.URLToImage)
		fmt.Println("----")
		geminiResponse(article.URL)
		break
	}
	db, err := InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	article := News{
		Title:       "Sample News Title",
		Tags:        "Tech, AI",
		Description: "This is a sample news article.",
		URL:         "https://example.com/news",
		URLToImage:  "https://example.com/image.jpg",
		PublishedAt: "2025-03-04",
	}

	err = InsertNews(db, article)
	if err != nil {
		log.Println("Failed to insert article:", err)
	} else {
		fmt.Println("Article inserted successfully")
	}
}
