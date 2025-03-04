package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const apiURLGemini = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent"
const apiKeyGemini = "AIzaSyAOdy--rJXYEsXy7Pvdh0YYhNEsBpSiQFg" // Замени на свой API-ключ

// Структура запроса
type RequestBody struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

// Ответ запроса
type Response struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func geminiResponse(question string) {
	// Прокси-сервер с логином и паролем
	proxyURL, _ := url.Parse("http://user204274:wdumt6@193.37.197.158:5167")
	transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}

	// Создаём клиента с прокси
	client := &http.Client{Transport: transport}

	// Создаём запрос
	reqBody := RequestBody{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{
				Parts: []struct {
					Text string `json:"text"`
				}{
					{Text: question},
				}, //https://gizmodo.com/musk-and-trumps-fort-knox-trip-is-about-bitcoin-2000569420
			},
		},
	}

	// Кодируем в JSON
	data, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", apiURLGemini+"?key="+apiKeyGemini, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка запроса:", err)
		return
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, _ := ioutil.ReadAll(resp.Body)
	var result Response
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Ошибка парсинга JSON:", err)
		return
	}

	// Достаем текст
	//fmt.Println(string(body))
	if len(result.Candidates) > 0 && len(result.Candidates[0].Content.Parts) > 0 {
		fmt.Println(result.Candidates[0].Content.Parts[0].Text)
	} else {
		fmt.Println("Нет данных в ответе")
	}
}
