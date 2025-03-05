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

func newsTagUpdate(url string, tag string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка запроса:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var news NewsResponse
	json.Unmarshal(body, &news)

	db, err := InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	for _, article := range news.Articles {
		textik := scrapperCollyan(article.URL)
		fmt.Println("Заголовок:", article.Title)
		fmt.Println("Теги:", article.Tags)
		fmt.Println("Ссылка:", article.URL)
		fmt.Println("Картинка:", article.URLToImage)
		fmt.Println("----")
		aiDescribe := geminiResponse("Напиши краткую выжимку по тексту в 3 предложения" + textik + "Если текста нет, выведи слово error")
		//fmt.Println(aiDescribe)
		article := News{
			Title:       article.Title,
			Tags:        tag,
			Description: aiDescribe,
			URL:         article.URL,
			URLToImage:  article.URLToImage,
			PublishedAt: article.PublishedAt,
		}

		selectTagBD10("культура")
		err = InsertNews(db, article)
		if err != nil {
			log.Println("Failed to insert article:", err)
		} else {
			fmt.Println("Article inserted successfully")
		}
	}
}

func newsUpdate() {
	listTagsUS := []string{"general"}
	listTagsRU := []string{"Общее", "Политика", "Культура", "Спорт", "Экономика", "Криптовалюта", "Технологии"}
	apiKeyNews := "8edd255699aa4fdfa562f51ac68de15e"
	for i, tag := range listTagsUS {
		fmt.Println("ПОшло поехало", tag)
		//url := "https://newsapi.org/v2/" + tag + "&pageSize=1&apiKey=" + apiKeyNews
		url := "https://newsapi.org/v2/top-headlines?country=us&category=" + tag + "&pageSize=1&apiKey=" + apiKeyNews
		newsTagUpdate(url, listTagsRU[i])
	}
}
