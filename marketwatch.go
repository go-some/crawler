package crawler

import (
	"github.com/gocolly/colly"
	"strings"
)

type MarketWatch struct {
}

func (rc *MarketWatch) Run(wtr DocsWriter) {
	c := colly.NewCollector()

	docs := make([]News, 0, 100)

	c.OnHTML("div.element--article", func(e *colly.HTMLElement) {
		// site-specific patterns
		title := strings.Trim(e.ChildText("h3.article__headline"), " ")
		url := e.ChildAttr(".article__headline a", "href")
		body := ""
		time := e.ChildText(".article__timestamp")
		origin := "https://www.marketwatch.com/"
		// [TODO] add validation
		doc := News{title, body, time, url, origin}
		docs = append(docs, doc)
	})

	c.OnScraped(func(r *colly.Response) {
		wtr.WriteDocs(docs)
	})

	c.Visit("https://www.marketwatch.com/investing/stocks?mod=exchange-traded-funds")
}
