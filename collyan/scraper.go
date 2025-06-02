package collyan

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func ScrapperCollyan(url string) string {
	// Создаем коллектор без ограничений по доменам
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.MaxDepth(1),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"),
	)

	// Устанавливаем заголовки для имитации браузера
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
		r.Headers.Set("Cache-Control", "no-cache")
		r.Headers.Set("Pragma", "no-cache")
		r.Headers.Set("Sec-Ch-Ua", `"Chromium";v="122", "Not(A:Brand";v="24", "Google Chrome";v="122"`)
		r.Headers.Set("Sec-Ch-Ua-Mobile", "?0")
		r.Headers.Set("Sec-Ch-Ua-Platform", `"Windows"`)
		r.Headers.Set("Sec-Fetch-Dest", "document")
		r.Headers.Set("Sec-Fetch-Mode", "navigate")
		r.Headers.Set("Sec-Fetch-Site", "none")
		r.Headers.Set("Sec-Fetch-User", "?1")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
	})

	var articleText strings.Builder
	var title string

	// Получаем заголовок
	c.OnHTML("h1, .articleHeader h1, .article-title, .post-title, .entry-title, .headline", func(e *colly.HTMLElement) {
		title = strings.TrimSpace(e.Text)
	})

	// Получаем основной текст статьи
	c.OnHTML(`
		div.article__text p, 
		div.article__block p, 
		div.b-material p,
		.article-content p,
		#article-body p,
		.WYSIWYG p,
		.article p,
		.post-content p,
		.entry-content p,
		.content p,
		.text p,
		[itemprop="articleBody"] p
	`, func(e *colly.HTMLElement) {
		text := strings.TrimSpace(e.Text)
		if text != "" {
			articleText.WriteString(text)
			articleText.WriteString("\n")
		}
	})

	// Обработка ошибок
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Ошибка при парсинге %s: %v\n", url, err)
		fmt.Printf("Статус код: %d\n", r.StatusCode)
		fmt.Printf("Заголовки ответа: %v\n", r.Headers)
	})

	// Добавляем задержку перед запросом
	time.Sleep(1 * time.Second)

	// Посещаем страницу
	err := c.Visit(url)
	if err != nil {
		fmt.Printf("Ошибка при посещении %s: %v\n", url, err)
		return ""
	}

	// Формируем полный текст статьи
	fullText := title + "\n\n" + articleText.String()

	// Очищаем текст от лишних пробелов и переносов строк
	fullText = strings.TrimSpace(fullText)
	fullText = strings.ReplaceAll(fullText, "\n\n\n", "\n\n")

	// Проверяем, что текст не пустой
	if fullText == "" {
		fmt.Printf("Не удалось получить текст статьи с %s\n", url)
		return ""
	}

	return fullText
}
