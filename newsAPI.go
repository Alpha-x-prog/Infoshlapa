package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

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
		aiDescribe := geminiResponse("Напиши краткую выжимку на русском языке по тексту в 3 предложения без вступления от лица я" + textik + "Если текста нет, выведи слово error")
		result := strings.Index(aiDescribe, "error")
		if result != -1 {
			aiDescribe = geminiResponse("Переведи данный текст на русский язык без вступления от лица я" + article.Description)
			fmt.Println("Сработала тревога")
		}
		//fmt.Println(aiDescribe)
		article := News{
			Title:       article.Title,
			Tags:        tag,
			Description: aiDescribe,
			URL:         article.URL,
			URLToImage:  article.URLToImage,
			PublishedAt: article.PublishedAt,
		}

		exists, errSim := checkSimilarTitle(db, article.Title, tag)
		if errSim != nil {
			log.Fatal(err)
		}
		if exists {
			fmt.Println("Новость с таким заголовком уже есть в последних 20 записях данной категории.")
		} else {
			err = InsertNews(db, article)
			if err != nil {
				log.Println("Failed to insert article:", err)
			} else {
				fmt.Println("Article inserted successfully")
			}
		}
	}
}

func newsUpdate() {
	listTagsUS := []string{"general", "business", "entertainment", "health", "science", "sports", "technology"}
	listTagsRU := []string{"общее", "бизнес", "развлечения", "здоровье", "наука", "спорт", "технологии"}
	apiKeyNews := "8edd255699aa4fdfa562f51ac68de15e"
	for i, tag := range listTagsUS {
		fmt.Println("Пошло поехало", tag)
		//url := "https://newsapi.org/v2/" + tag + "&pageSize=1&apiKey=" + apiKeyNews
		url := "https://newsapi.org/v2/top-headlines?country=us&category=" + tag + "&pageSize=1&apiKey=" + apiKeyNews
		newsTagUpdate(url, listTagsRU[i])
	}
}
