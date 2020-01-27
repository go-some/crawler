package crawler

import (
	"github.com/gocolly/colly"
	"strings"
)

type Reuters struct {
}

func (rc *Reuters) Run(wtr DocsWriter) {
	// Instantiate default NewCollector
	c := colly.NewCollector()
	docs := make([]News, 0, 100)

	// Create another collector to scrape each news article
	articleCollector := c.Clone()

	c.OnHTML(".story", func(e *colly.HTMLElement) {
		//find article url and visit
		article_url := e.Request.AbsoluteURL(e.ChildAttr("a", "href"))
		if strings.Index(article_url, "reuters.com/article") != -1 {
			articleCollector.Visit(article_url)
		}
	})

	articleCollector.OnHTML("div.StandardArticle_inner-container", func(e *colly.HTMLElement) {
		//read article and save
		doc := News{
			Title:  e.ChildText(".ArticleHeader_headline"),
			Body:   e.ChildText("div.StandardArticleBody_body"),
			Time:   e.ChildText(".ArticleHeader_date"),
			Url:    e.Request.URL.String(),
			Origin: "Reuters",
		}
		docs = append(docs, doc)
	})

	c.OnScraped(func(r *colly.Response) {
		wtr.WriteDocs(docs)
	})

	c.Visit("https://www.reuters.com/finance")
}
