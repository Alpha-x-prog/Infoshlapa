package collyan

import (
	"fmt"
	"github.com/gocolly/colly"
)

func ScrapperCollyan(url string) string {
	c := colly.NewCollector()

	var articleText string
	// Find and visit all links
	c.OnHTML("p", func(e *colly.HTMLElement) {
		articleText += e.Text
	})

	//c.OnRequest(func(r *colly.Request) {
	//	fmt.Println("Visiting", r.URL)
	//})

	//c.Visit("https://www.cnn.com/2025/03/03/entertainment/vanity-fair-oscar-party-2025/index.html")
	fmt.Println("Передаётся ссылка", url)
	c.Visit(url)

	return articleText

}
