package crawler

import (
	"fmt"
	"github.com/gocolly/colly"
	"regexp"
	"strings"
)

type MarketWatch struct {
}

func (rc *MarketWatch) Run(wtr DocsWriter) {
	rootCollector := colly.NewCollector(
		colly.MaxDepth(3),
		colly.URLFilters(
			regexp.MustCompile("https://www\\.marketwatch\\.com/"),
			regexp.MustCompile("https://www\\.marketwatch\\.com/story/.+"),
		),
		//colly.DisallowedURLFilters(
		//	regexp.MustCompile("https://www\\.usatoday\\.com/opinion/"),
		//),
	)
	rootCollector.AllowURLRevisit = false

	articleCollector := colly.NewCollector()

	rootCollector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if strings.Index(link, "marketwatch.com/story/") != -1 {
			articleCollector.Visit(link)
		} else {
			e.Request.Visit(link)
		}
	})

	articleCollector.OnHTML(".region--primary", func(e *colly.HTMLElement) {
		doc := News{
			Title:  e.ChildText("h1.article__headline"),
			Body:   e.ChildText("div.article__body "),
			Time:   e.ChildText(".timestamp "),
			Url:    e.Request.URL.String(),
			Origin: "MarketWatch",
		}
		cnt, err := wtr.WriteDocs([]News{doc})
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(cnt, "docs saved")
		}
	})

	rootCollector.Visit("https://www.marketwatch.com/")
}
