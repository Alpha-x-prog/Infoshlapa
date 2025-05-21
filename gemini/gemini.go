package gemini

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const apiURLGemini = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent"

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
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func GeminiResponse(question string) string {
	// Получаем API ключ из переменных окружения
	apiKeyGemini := os.Getenv("GEMINI_API_KEY")
	if apiKeyGemini == "" {
		fmt.Println("GEMINI_API_KEY не найден в переменных окружения")
		return "error"
	}

	// Получаем настройки прокси из переменных окружения
	proxyURL := os.Getenv("PROXY_URL")
	proxyUser := os.Getenv("PROXY_USER")
	proxyPass := os.Getenv("PROXY_PASSWORD")

	// Формируем URL прокси с учетными данными
	proxyURLWithAuth := fmt.Sprintf("http://%s:%s@%s", proxyUser, proxyPass, proxyURL)
	proxy, _ := url.Parse(proxyURLWithAuth)
	transport := &http.Transport{Proxy: http.ProxyURL(proxy)}

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
				},
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
		return "error"
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Ответ от API:", string(body))

	var result Response
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Ошибка парсинга JSON:", err)
		return "error"
	}

	// Проверяем наличие ошибки в ответе
	if result.Error.Code != 0 {
		fmt.Printf("Ошибка API: %s (код: %d)\n", result.Error.Message, result.Error.Code)
		return "error"
	}

	// Достаем текст
	if len(result.Candidates) > 0 && len(result.Candidates[0].Content.Parts) > 0 {
		return result.Candidates[0].Content.Parts[0].Text
	} else {
		fmt.Println("Нет данных в ответе")
		return "error"
	}
}
