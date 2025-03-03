package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type NewsResponse struct {
	Articles []struct {
		Title      string `json:"title"`
		Url        string `json:"url"`
		UrlToImage string `json:"urlToImage"`
	} `json:"articles"`
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
		fmt.Println("Ссылка:", article.Url)
		fmt.Println("Картинка:", article.UrlToImage)
		fmt.Println("----")
		geminiResponse(article.Url)
		break
	}
}
