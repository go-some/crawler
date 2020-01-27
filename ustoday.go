package crawler

import (
	"github.com/gocolly/colly"
	"strings"
)

type USToday struct {
}

func (rc *USToday) Run(wtr DocsWriter) {
	c := colly.NewCollector()

	docs := make([]News, 0, 100)

	c.OnHTML("a.gnt_m_flm_a", func(e *colly.HTMLElement) {
		// site-specific patterns
		title := strings.Trim(e.Text, " ")
		url := e.Attr("href")
		body := strings.Trim(e.Attr("data-c-br"), " ")
		time := e.ChildAttr("div", "data-c-dt")
		origin := "https://www.usatoday.com/"
		// [TODO] add validation
		doc := News{title, body, time, url, origin}
		docs = append(docs, doc)
	})

	c.OnScraped(func(r *colly.Response) {
		wtr.WriteDocs(docs)
	})

	c.Visit("https://www.usatoday.com/money/")
}
